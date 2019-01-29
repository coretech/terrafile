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
		"[tf-aws-vpc] Checking out v1.46.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"[tf-aws-vpc-experimental] Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"[tf-aws-vpc-commit] Checking out 01601169c00c68f37d5df8a80cc17c88f02c04d0 of git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"[tf-aws-vpc-default] Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"[tf-aws-vpc-path] Checking out examples/simple-vpc from master of git@github.com:terraform-aws-modules/terraform-aws-vpc",
	} {
		assert.Contains(t, testcli.Stdout(), output)
	}
	// Assert files exist
	for _, moduleName := range []string{
		"tf-aws-vpc",
		"tf-aws-vpc-experimental",
		"tf-aws-vpc-commit",
		"tf-aws-vpc-default",
		"tf-aws-vpc-path",
	} {
		assert.DirExists(t, path.Join(workingDirectory, "vendor/modules", moduleName))
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
	var yaml = `tf-aws-vpc:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: "v1.46.0"
tf-aws-vpc-experimental:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: "master"
tf-aws-vpc-commit:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  commit: "01601169c00c68f37d5df8a80cc17c88f02c04d0"
tf-aws-vpc-default:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
tf-aws-vpc-path:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  path: "examples/simple-vpc"
  version: "master"
`
	createFile(t, path.Join(folder, "Terrafile"), yaml)
}
