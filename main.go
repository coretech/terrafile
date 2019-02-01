package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/nritholtz/stdemuxerhook"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type module struct {
	Source  string `yaml:"source"`
	Version string `yaml:"version"`
	Commit  string `yaml:"commit"`
	Path    string `yaml:"path"`
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

func gitClonePath(repository string, path string, version string, moduleName string) {
	moduleDir := opts.ModulePath + "/" + moduleName
	tmpModuleDir := opts.ModulePath + "/.tmppath_" + moduleName

	os.MkdirAll(tmpModuleDir+"/.git/info", os.ModePerm)
	fh, err := os.Create(tmpModuleDir + "/.git/info/sparse-checkout")
	if err != nil {
		log.Fatalln(err)
	}
	defer fh.Close()
	_, err = io.Copy(fh, strings.NewReader(path))
	if err != nil {
		log.Fatalln(err)
	}

	for _, command := range []string{
		"git init . ",
		fmt.Sprintf("git remote add origin %s", repository),
		"git config core.sparsecheckout true",
		fmt.Sprintf("git pull origin %s", version),
	} {
		c := strings.Fields(command)[0]
		a := strings.Fields(command)[1:]

		cmd := exec.Command(c, a...)
		cmd.Dir = tmpModuleDir
		err := cmd.Run()
		if err != nil {
			log.Fatalln(err)
		}
	}
	err = os.Rename(tmpModuleDir+"/"+path, moduleDir)
	err = os.RemoveAll(tmpModuleDir)
	if err != nil {
		log.Fatalln(err)
	}
}

func gitCheckoutCommit(commit string, moduleName string) {
	cmd := exec.Command("git", "checkout", commit)
	cmd.Dir = opts.ModulePath + "/" + moduleName
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func rmGit(moduleName string) {
	gitDir := opts.ModulePath + "/" + moduleName + "/.git"
	err := os.RemoveAll(gitDir)
	if err != nil {
		log.Fatalln(err)
	}
}

func logFetch(repository string, version string, commit string, path string, moduleName string) {
	var moduleVersion string
	var pathFrom string
	moduleVersion = version
	pathFrom = ""
	if commit != "" {
		moduleVersion = commit
	}
	if path != "" {
		pathFrom = fmt.Sprintf("%s from ", path)
	}
	log.Printf("[%s] Checking out %s%s of %s", moduleName, pathFrom, moduleVersion, repository)
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
		logFetch(module.Source, module.Version, module.Commit, module.Path, key)

		// Checkout path or sparse checkout specified path
		if module.Path == "" {
			gitClone(module.Source, module.Version, key)
		} else {
			gitClonePath(module.Source, module.Path, module.Version, key)
		}
		// Checkout a commit if specified
		if module.Commit != "" {
			gitCheckoutCommit(module.Commit, key)
		}
		rmGit(key)
	}
}
