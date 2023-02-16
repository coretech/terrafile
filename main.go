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
	Source       string   `yaml:"source"`
	Version      string   `yaml:"version"`
	Destinations []string `yaml:"destinations"`
}

var opts struct {
	ModulePath string `short:"p" long:"module_path" default:"./vendor/modules" description:"File path to install generated terraform modules, if not overridden by 'destinations:' field"`

	TerrafilePath string `short:"f" long:"terrafile_file" default:"./Terrafile" description:"File path to the Terrafile file"`

	Clean bool `short:"c" long:"clean" description:"Remove everything from destinations and module path upon fetching module(s)\n !!! WARNING !!! Removes all files and folders in the destinations including non-modules."`
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
	cleanupPath := filepath.Join(destinationDir, moduleName)
	log.Printf("[*] Removing previously cloned artifacts at %s", cleanupPath)
	_ = os.RemoveAll(cleanupPath)
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

	workDirAbsolutePath, err := os.Getwd()
	if err != nil {
		log.Errorf("failed to get working directory absolute path due to: %s", err)
	}

	// Read File
	yamlFile, err := os.ReadFile(opts.TerrafilePath)
	if err != nil {
		log.Fatalf("failed to read configuration in file %s due to error: %s", opts.TerrafilePath, err)
	}

	// Parse File
	var config map[string]module
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatalf("failed to parse yaml file due to error: %s", err)
	}

	if opts.Clean {
		cleanDestinations(config)
	}

	// Clone modules
	var wg sync.WaitGroup
	_ = os.RemoveAll(opts.ModulePath)
	_ = os.MkdirAll(opts.ModulePath, os.ModePerm)

	for key, mod := range config {
		wg.Add(1)
		go func(m module, key string) {
			defer wg.Done()

			// path to clone module
			cloneDestination := opts.ModulePath
			// list of paths to link module to. empty, unless Destinations are more than 1 location
			var linkDestinations []string

			if m.Destinations != nil && len(m.Destinations) > 0 {
				// set first in Destinations as location to clone to
				cloneDestination = filepath.Join(m.Destinations[0], opts.ModulePath)
				// the rest of Destinations are locations to link module to
				linkDestinations = m.Destinations[1:]

			}

			// create folder to clone into
			if err := os.MkdirAll(cloneDestination, os.ModePerm); err != nil {
				log.Errorf("failed to create folder %s due to error: %s", cloneDestination, err)

				// no reason to continue as failed to create folder
				return
			}

			// clone repository
			gitClone(m.Source, m.Version, key, cloneDestination)

			for _, d := range linkDestinations {
				// the source location as folder where module was cloned and module folder name
				moduleSrc := filepath.Join(workDirAbsolutePath, cloneDestination, key)
				// append destination path with module path
				dst := filepath.Join(d, opts.ModulePath)

				log.Infof("[*] Creating folder %s", dst)
				if err := os.MkdirAll(dst, os.ModePerm); err != nil {
					log.Errorf("failed to create folder %s due to error: %s", dst, err)
					return
				}

				dst = filepath.Join(dst, key)

				log.Infof("[*] Remove existing artifacts at %s", dst)
				if err := os.RemoveAll(dst); err != nil {
					log.Errorf("failed to remove location %s due to error: %s", dst, err)
					return
				}

				log.Infof("[*] Link %s to %s", moduleSrc, dst)
				if err := os.Symlink(moduleSrc, dst); err != nil {
					log.Errorf("failed to link module from %s to %s due to error: %s", moduleSrc, dst, err)
				}
			}
		}(mod, key)
	}

	wg.Wait()
}

func cleanDestinations(config map[string]module) {

	// Map filters duplicate destinations with key being each destination's file path
	uniqueDestinations := make(map[string]bool)

	// Range over config and gather all unique destinations
	for _, m := range config {
		if len(m.Destinations) == 0 {
			uniqueDestinations[opts.ModulePath] = true
			continue
		}

		// range over Destinations and put them into map
		for _, dst := range m.Destinations {
			// Destination supposed to be conjunction of destination defined in file with module path
			d := filepath.Join(dst, opts.ModulePath)
			uniqueDestinations[d] = true
		}
	}

	for dst := range uniqueDestinations {

		log.Infof("[*] Removing artifacts from %s", dst)
		if err := os.RemoveAll(dst); err != nil {
			log.Errorf("Failed to remove artifacts from %s due to error: %s", dst, err)
		}
	}
}
