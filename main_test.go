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

	defer func() {
		assert.NoError(t, os.RemoveAll("testdata/"))
	}()

	if !testcli.Success() {
		t.Fatalf("Expected to succeed, but failed: %q with message: %q", testcli.Error(), testcli.Stderr())
	}
	// Assert output
	for _, output := range []string{
		"Checking out v1.46.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"Checking out v1.36.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"Checking out v1.35.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc",
	} {
		assert.Contains(t, testcli.Stdout(), output)
	}
	// Assert folder exist with default destination
	for _, moduleName := range []string{
		"tf-aws-vpc",
		"tf-aws-vpc-experimental",
	} {
		assert.DirExists(t, path.Join(workingDirectory, "vendor/modules", moduleName))
	}

	// Assert folder exist with non-default destination
	for _, moduleName := range []string{
		"testdata/stackA/vendor/modules/tf-aws-vpc-legacy",
		"testdata/stackA/vendor/modules/tf-aws-vpc-legacy2",
		"testdata/stackB/vendor/modules/tf-aws-vpc-legacy2",
		"testdata/stackC/vendor/modules/tf-aws-vpc-legacy2",
	} {
		assert.DirExists(t, path.Join(workingDirectory, moduleName))
	}

	// Assert files exist with default destination
	for _, moduleName := range []string{
		"tf-aws-vpc/main.tf",
		"tf-aws-vpc-experimental/main.tf",
	} {
		assert.FileExists(t, path.Join(workingDirectory, "vendor/modules", moduleName))
	}

	// Assert files exist with non-default destination
	for _, moduleName := range []string{
		"testdata/stackA/vendor/modules/tf-aws-vpc-legacy/main.tf",
		"testdata/stackA/vendor/modules/tf-aws-vpc-legacy2/main.tf",
		"testdata/stackB/vendor/modules/tf-aws-vpc-legacy2/main.tf",
		"testdata/stackC/vendor/modules/tf-aws-vpc-legacy2/main.tf",
	} {
		assert.FileExists(t, path.Join(workingDirectory, moduleName))
	}

	// Assert checked out correct version
	for moduleName, cloneOptions := range map[string]map[string]string{
		"tf-aws-vpc": map[string]string{
			"repository": "git@github.com:terraform-aws-modules/terraform-aws-vpc",
			"version":    "v1.46.0",
		},
		"tf-aws-vpc-experimental": map[string]string{
			"repository": "git@github.com:terraform-aws-modules/terraform-aws-vpc",
			"version":    "master",
		},
	} {
		testModuleLocation := path.Join(workingDirectory, "vendor/modules", moduleName+"__test")
		testcli.Run("git", "clone", "-b", cloneOptions["version"], cloneOptions["repository"], testModuleLocation)
		if !testcli.Success() {
			t.Fatalf("Expected to succeed, but failed: %q with message: %q", testcli.Error(), testcli.Stderr())
		}
		testcli.Run("diff", "--exclude=.git", "-r", path.Join(workingDirectory, "vendor/modules", moduleName), testModuleLocation)
		if !testcli.Success() {
			t.Fatalf("File difference found for %q, with failure: %q with message: %q", moduleName, testcli.Error(), testcli.Stderr())
		}
	}

	// Assert checked out correct version to non-default destination
	for dst, checkout := range map[string]map[string]map[string]string{
		"testdata/stackA/vendor/modules": map[string]map[string]string{
			"tf-aws-vpc-legacy": map[string]string{
				"repository": "git@github.com:terraform-aws-modules/terraform-aws-vpc",
				"version":    "v1.36.0",
			},
			"tf-aws-vpc-legacy2": map[string]string{
				"repository": "git@github.com:terraform-aws-modules/terraform-aws-vpc",
				"version":    "v1.35.0",
			},
		},
		"testdata/stackB/vendor/modules": map[string]map[string]string{
			"tf-aws-vpc-legacy2": map[string]string{
				"repository": "git@github.com:terraform-aws-modules/terraform-aws-vpc",
				"version":    "v1.35.0",
			},
		},
		"testdata/stackC/vendor/modules": map[string]map[string]string{
			"tf-aws-vpc-legacy2": map[string]string{
				"repository": "git@github.com:terraform-aws-modules/terraform-aws-vpc",
				"version":    "v1.35.0",
			},
		},
	} {
		for moduleName, cloneOptions := range checkout {
			testModuleLocation := path.Join(workingDirectory, dst, moduleName+"__test")
			os.RemoveAll(testModuleLocation)
			testcli.Run("git", "clone", "-b", cloneOptions["version"], cloneOptions["repository"], testModuleLocation)
			if !testcli.Success() {
				t.Fatalf("Expected to succeed, but failed: %q with message: %q", testcli.Error(), testcli.Stderr())
			}
			testcli.Run("diff", "--exclude=.git", "-r", path.Join(workingDirectory, dst, moduleName), testModuleLocation)
			if !testcli.Success() {
				t.Fatalf("File difference found for %q, with failure: %q with message: %q", moduleName, testcli.Error(), testcli.Stderr())
			}
		}
	}
}

func setup(t *testing.T) (current string, back func()) {
	folder, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	createTerrafile(t, folder)
	return folder, func() {
		assert.NoError(t, os.RemoveAll(folder))
	}
}

func createFile(t *testing.T, filename string, contents string) {
	assert.NoError(t, os.WriteFile(filename, []byte(contents), 0644))
}

func createTerrafile(t *testing.T, folder string) {
	var yaml = `tf-aws-vpc:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: "v1.46.0"
tf-aws-vpc-experimental:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: "master"
tf-aws-vpc-legacy:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: "v1.36.0"
  destination: 
    - "./testdata/stackA"
tf-aws-vpc-legacy2:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: "v1.35.0"
  destination:
    - "./testdata/stackA"
    - "./testdata/stackB"
    - "./testdata/stackC"
`
	createFile(t, path.Join(folder, "Terrafile"), yaml)
}
