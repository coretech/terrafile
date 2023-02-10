/*
Copyright 2022 IDT Corp.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/nritholtz/stdemuxerhook"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type module struct {
	Source      string   `yaml:"source"`
	Version     string   `yaml:"version"`
	Destination []string `yaml:"destination"`
}

var opts struct {
	ModulePath string `short:"p" long:"module_path" default:"./vendor/modules" description:"File path to install generated terraform modules, if not overridden by 'destination:' field"`

	TerrafilePath string `short:"f" long:"terrafile_file" default:"./Terrafile" description:"File path to the Terrafile file"`
}

// To be set by goreleaser on build
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	// Needed to redirect logrus to proper stream STDOUT vs STDERR
	log.AddHook(stdemuxerhook.New(log.StandardLogger()))
}

func gitClone(repository string, version string, moduleName string, destinationDir string) {
	log.Printf("[*] Deleting existing clonned artifacts")
	os.RemoveAll(filepath.Join(destinationDir, moduleName))
	log.Printf("[*] Checking out %s of %s \n", version, repository)
	cmd := exec.Command("git", "clone", "--single-branch", "--depth=1", "-b", version, repository, moduleName)
	cmd.Dir = destinationDir
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to clone repository %s due to error: %s", cmd.String(), err)
	}
}

func main() {
	fmt.Printf("Terrafile: version %v, commit %v, built at %v \n", version, commit, date)
	_, err := flags.Parse(&opts)

	// Invalid choice
	if err != nil {
		log.Errorf("failed to parse flags due to: %s", err)
		os.Exit(1)
	}

	// Read File
	yamlFile, err := os.ReadFile(opts.TerrafilePath)
	if err != nil {
		log.Fatalf("failed to read configuration in fie %s due to error: %s", opts.TerrafilePath, err)
	}

	// Parse File
	var config map[string]module
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatalf("failed to parce yaml file due to error: %s", err)
	}

	// Clone modules
	var wg sync.WaitGroup
	_ = os.RemoveAll(opts.ModulePath)
	_ = os.MkdirAll(opts.ModulePath, os.ModePerm)

	for key, mod := range config {
		wg.Add(1)
		go func(m module, key string) {
			defer wg.Done()

			firstDestination := opts.ModulePath
			skipCopy := true
			if m.Destination != nil && len(m.Destination) > 0 {
				firstDestination = filepath.Join(m.Destination[0], opts.ModulePath)
				skipCopy = false
			}

			if err := os.MkdirAll(firstDestination, os.ModePerm); err != nil {
				log.Errorf("failed to create folder %s due to error: %s", firstDestination, err)
				return
			}

			gitClone(m.Source, m.Version, key, firstDestination)
			// Delete .git folder
			_ = os.RemoveAll(filepath.Join(firstDestination, key, ".git"))

			if skipCopy {
				return
			}

			for _, d := range m.Destination[1:] {
				dst := filepath.Join(d, opts.ModulePath)
				wg.Add(1)
				go func(dst string, m module, key string) {
					defer wg.Done()
					if err := os.MkdirAll(dst, os.ModePerm); err != nil {
						log.Errorf("failed to create folder %s due to error: %s", dst, err)
						return
					}
					os.RemoveAll(filepath.Join(dst, key))
					moduleSrc := filepath.Join(firstDestination, key)
					cmd := exec.Command("cp", "-Rf", moduleSrc, dst)
					if err := cmd.Run(); err != nil {
						log.Errorf("failed to copy module from %s to %s due to error: %s", moduleSrc, dst, err)
					}
				}(dst, m, key)
			}
		}(mod, key)
	}

	wg.Wait()
}
