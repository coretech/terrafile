package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/nritholtz/stdemuxerhook"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var opts struct {
	ModulePath string `short:"p" long:"module_path" default:"./.terrafile/vendor" description:"File path to install generated terraform modules"`

	TerrafilePath string `short:"f" long:"terrafile_file" default:"./.terrafile/Terrafile" description:"File path to the Terrafile file"`
}

// To be set by goreleaser on build
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var previousRepo = "/tmp" // Used to speed up git clone using --reference-if-able 27s vs. 8.6s for 7 clones

func init() {
	// Needed to redirect logrus to proper stream STDOUT vs STDERR
	log.AddHook(stdemuxerhook.New(log.StandardLogger()))
}

func gitClone(repositoryPath string, version string, referenceRepo string) string {
	pathParts := strings.Split(repositoryPath, ":")
	repositoryName := pathParts[1]
	targetPath := fmt.Sprintf("%s/refs/%s", repositoryName, version)

	log.Printf("[*] Checking out %s of %s -> %s/%s\n", version, repositoryPath, opts.ModulePath, targetPath)

	args := []string{"clone", "-b", version, repositoryPath, targetPath}
	if referenceRepo != "" {
		args = append(args, "--reference-if-able", referenceRepo)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = opts.ModulePath
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	return targetPath
}

func main() {
	//fmt.Printf("Terrafile: version %v, commit %v, built at %v \n", version, commit, date)
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
	var config map[string][]string
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatalln(err)
	}

	// Clone modules
	os.RemoveAll(opts.ModulePath)
	os.MkdirAll(opts.ModulePath, os.ModePerm)
	for source, refs := range config {
		var referenceRepo string
		for _, ref := range refs {
			referenceRepo = gitClone(source, ref, referenceRepo)
		}
	}
}
