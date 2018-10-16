package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/jessevdk/go-flags"
	"github.com/nritholtz/stdemuxerhook"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type module struct {
	Source  string `yaml:"source"`
	Version string `yaml:"version"`
}

var opts struct {
	ModulePath string `short:"p" long:"module_path" default:"./vendor/modules" description:"File path to install generated terraform modules"`

	TerrafilePath string `short:"f" long:"terrafile_file" default:"." description:"File path in which the Terrafile file is located"`
}

func init() {
	// Needed to redirect logrus to proper stream STDOUT vs STDERR
	log.AddHook(stdemuxerhook.New(log.StandardLogger()))
}

func gitClone(repository string, version string, moduleName string) {
	log.Printf("[*] Checking out %s of %s \n", version, repository)
	cmd := exec.Command("git", "clone", "-b", version, repository, moduleName)
	cmd.Dir = opts.ModulePath
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	_, err := flags.Parse(&opts)

	// Invalid choice
	if err != nil {
		os.Exit(1)
	}

	// Read File
	yamlFile, err := ioutil.ReadFile(fmt.Sprint(opts.TerrafilePath, "/Terrafile"))
	if err != nil {
		log.Fatalln(err)
	}

	// Parse File
	var config map[string]module
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatalln(err)
	}

	// Clone modules
	os.RemoveAll(opts.ModulePath)
	os.MkdirAll(opts.ModulePath, os.ModePerm)
	for key, module := range config {
		gitClone(module.Source, module.Version, key)
	}
}
