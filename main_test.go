package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/assert"
)

var terrafileBinaryPath string
var workingDirectory string

func init() {
	var err error
	workingDirectory, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	terrafileBinaryPath = workingDirectory + "/terrafile"
}
func TestTerraformWithTerrafilePath(t *testing.T) {
	folder, back := setup(t)
	defer back()

	testcli.Run(terrafileBinaryPath, "-f", fmt.Sprint(folder, "/Terrafile"))

	if !testcli.Success() {
		t.Fatalf("Expected to succeed, but failed: %q with message: %q", testcli.Error(), testcli.Stderr())
	}
	// Assert output
	for _, output := range []string{
		"Checking out v1.46.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc -> ./.terrafile/vendor/terraform-aws-modules/terraform-aws-vpc/refs/v1.46.0",
		"Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc -> ./.terrafile/vendor/terraform-aws-modules/terraform-aws-vpc/refs/master",
	} {
		assert.Contains(t, testcli.Stdout(), output)
	}
	// Assert files exist
	for _, moduleName := range []string{
		"terraform-aws-modules/terraform-aws-vpc/refs/master",
		"terraform-aws-modules/terraform-aws-vpc/refs/v1.46.0",
	} {
		assert.DirExists(t, path.Join(workingDirectory, "./.terrafile/vendor", moduleName))
	}
}

func setup(t *testing.T) (current string, back func()) {
	folder, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	createTerrafile(t, folder)
	return folder, func() {
		assert.NoError(t, os.RemoveAll(folder))
	}
}

func createFile(t *testing.T, filename string, contents string) {
	assert.NoError(t, ioutil.WriteFile(filename, []byte(contents), 0644))
}

func createTerrafile(t *testing.T, folder string) {
	var yaml = `git@github.com:terraform-aws-modules/terraform-aws-vpc:
  - v1.46.0
  - master
`
	createFile(t, path.Join(folder, "Terrafile"), yaml)
}
