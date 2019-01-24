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
	Commit  string `yaml:"commit"`
}

var opts struct {
	ModulePath string `short:"p" long:"module_path" default:"./vendor/modules" description:"File path to install generated terraform modules"`

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

func gitClone(repository string, version string, moduleName string) {
	cmd := exec.Command("git", "clone", "-b", version, repository, moduleName)
	cmd.Dir = opts.ModulePath
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func gitCheckout(commit string, moduleName string) {
	cmd := exec.Command("git", "checkout", commit)
	cmd.Dir = opts.ModulePath + "/" + moduleName
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func logFetch(repository string, version string, commit string, moduleName string) {
	var moduleVersion string
	moduleVersion = version
	if commit != "" {
		moduleVersion = commit
	}
	log.Printf("[%s] Checking out %s of %s", moduleName, moduleVersion, repository)
}

func main() {
	fmt.Printf("Terrafile: version %v, commit %v, built at %v \n", version, commit, date)
	_, err := flags.Parse(&opts)

	// Invalid choice
	if err != nil {
		os.Exit(1)
	}

	// Read File
	yamlFile, err := ioutil.ReadFile(opts.TerrafilePath)
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
		if len(module.Version) == 0 {
			module.Version = "master"
		}
		logFetch(module.Source, module.Version, module.Commit, key)
		gitClone(module.Source, module.Version, key)
		// Checkout a commit if specified
		if module.Commit != "" {
			gitCheckout(module.Commit, key)
		}
	}
}
