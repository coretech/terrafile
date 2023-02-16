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

	defer println(testcli.Stdout())
	defer println(testcli.Stderr())

	defer func() {
		assert.NoError(t, os.RemoveAll("testdata/"))
	}()

	if !testcli.Success() {
		t.Fatalf("Expected to succeed, but failed: %q with message: %q", testcli.Error(), testcli.Stderr())
	}
	// Assert output
	for _, output := range []string{
		"Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"Checking out v1.46.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"Checking out v3.2.0 of git@github.com:terraform-aws-modules/terraform-aws-vpn-gateway",
		"Checking out v3.6.1 of git@github.com:terraform-aws-modules/terraform-aws-s3-bucket",
		"Checking out v5.11.1 of git@github.com:terraform-aws-modules/terraform-aws-iam",
	} {
		assert.Contains(t, testcli.Stdout(), output)
	}

	// Assert output does not contain the following
	for _, output := range []string{
		"[*] Removing artifacts from ./vendor/modules",
		"[*] Removing artifacts from testdata/networking/vendor/modules",
		"[*] Removing artifacts from testdata/some-other-stack/vendor/modules",
		"[*] Removing artifacts from testdata/iam/vendor/modules",
		"[*] Removing artifacts from testdata/onboarding/vendor/modules",
	} {
		assert.NotContains(t, testcli.Stdout(), output)
	}

	// Assert folder exist with default destination
	for _, moduleName := range []string{
		"tf-aws-vpc",
		"tf-aws-vpc-experimental",
	} {
		assert.DirExists(t, path.Join(workingDirectory, "vendor/modules", moduleName))
	}

	// Assert folder exist with non-default destinations
	for _, moduleName := range []string{
		"testdata/networking/vendor/modules/tf-aws-vpn-gateway/",
		"testdata/networking/vendor/modules/tf-aws-s3-bucket/",
		"testdata/iam/vendor/modules/tf-aws-iam/",

		// Symlinks are not Dirs. But contents will be tested later on
		// "testdata/onboarding/vendor/modules/tf-aws-s3-bucket/",
		// "testdata/some-other-stack/vendor/modules/tf-aws-s3-bucket/",
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

	// Assert files exist with non-default destinations
	for _, moduleName := range []string{
		"testdata/networking/vendor/modules/tf-aws-vpn-gateway/main.tf",
		"testdata/networking/vendor/modules/tf-aws-s3-bucket/main.tf",

		// terraform-aws-modules/terraform-aws-iam doesn't have main.tf, as it represents set of modules
		// However, some terraform-aws-modules/terraform-aws-iam/modules have, e.g.:
		"testdata/iam/vendor/modules/tf-aws-iam/modules/iam-account/main.tf",
		"testdata/onboarding/vendor/modules/tf-aws-s3-bucket/main.tf",
		"testdata/some-other-stack/vendor/modules/tf-aws-s3-bucket/main.tf",
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

	// Assert checked out correct version to non-default destinations
	for dst, checkout := range map[string]map[string]map[string]string{
		"testdata/networking/vendor/modules": {
			"tf-aws-s3-bucket": {
				"repository": "git@github.com:terraform-aws-modules/terraform-aws-s3-bucket",
				"version":    "v3.6.1",
			},
			"tf-aws-vpn-gateway": {
				"repository": "git@github.com:terraform-aws-modules/terraform-aws-vpn-gateway",
				"version":    "v3.2.0",
			},
		},
		"testdata/iam/vendor/modules": {
			"tf-aws-iam": {
				"repository": "git@github.com:terraform-aws-modules/terraform-aws-iam",
				"version":    "v5.11.1",
			},
		},
		"testdata/onboarding/vendor/modules": {
			"tf-aws-s3-bucket": {
				"repository": "git@github.com:terraform-aws-modules/terraform-aws-s3-bucket",
				"version":    "v3.6.1",
			},
		},
		"testdata/some-other-stack/vendor/modules": {
			"tf-aws-s3-bucket": {
				"repository": "git@github.com:terraform-aws-modules/terraform-aws-s3-bucket",
				"version":    "v3.6.1",
			},
		},
	} {
		for moduleName, cloneOptions := range checkout {
			testModuleLocation := path.Join(workingDirectory, dst, moduleName+"__test")
			_ = os.RemoveAll(testModuleLocation)
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

func TestTerraformWithTerrafileClean(t *testing.T) {
	folder, back := setup(t)
	defer back()

	// Run original Terrafile with path
	testcli.Run(terrafileBinaryPath, "-f", fmt.Sprint(folder, "/Terrafile"))

	defer println(testcli.Stdout())
	defer println(testcli.Stderr())

	defer func() {
		assert.NoError(t, os.RemoveAll("testdata/"))
	}()

	// Create Terrafile2.yaml
	createTerrafile2(t, folder)

	// Run terrafile with path to Terrafile2 and clean flag
	testcli.Run(terrafileBinaryPath, "-f", fmt.Sprint(folder, "/Terrafile2"), "-c")

	// Print out output of second run
	defer println(testcli.Stdout())
	defer println(testcli.Stderr())

	if !testcli.Success() {
		t.Fatalf("Expected to succeed, but failed: %q with message: %q", testcli.Error(), testcli.Stderr())
	}

	// Assert output
	for _, output := range []string{
		"[*] Removing artifacts from ./vendor/modules",
		"[*] Removing artifacts from testdata/networking/vendor/modules",
		"[*] Removing artifacts from testdata/some-other-stack/vendor/modules",
		"[*] Removing artifacts from testdata/iam/vendor/modules",
		"[*] Removing artifacts from testdata/onboarding/vendor/modules",
	} {
		assert.Contains(t, testcli.Stdout(), output)
	}

	// Assert files exist with default destination
	for _, moduleName := range []string{
		"tf-aws-vpc/main.tf",
	} {
		assert.FileExists(t, path.Join(workingDirectory, "vendor/modules", moduleName))
	}

	// Assert files does not exist with default destination
	for _, moduleName := range []string{
		"tf-aws-vpc-experimental/main.tf",
	} {
		assert.NoFileExists(t, path.Join(workingDirectory, "vendor/modules", moduleName))
	}

	// Assert files exist with non-default destinations
	for _, moduleName := range []string{
		"testdata/networking/vendor/modules/tf-aws-vpn-gateway/main.tf",

		// terraform-aws-modules/terraform-aws-iam doesn't have main.tf, as it represents set of modules
		// However, some terraform-aws-modules/terraform-aws-iam/modules have, e.g.:
		"testdata/iam/vendor/modules/tf-aws-iam/modules/iam-account/main.tf",
		"testdata/onboarding/vendor/modules/tf-aws-s3-bucket/main.tf",
		"testdata/some-other-stack/vendor/modules/tf-aws-vpn-gateway/main.tf",
	} {
		assert.FileExists(t, path.Join(workingDirectory, moduleName))
	}

	// Assert files do not exist with non-default destinations
	for _, moduleName := range []string{
		"testdata/networking/vendor/modules/tf-aws-s3-bucket/main.tf",
		"testdata/some-other-stack/vendor/modules/tf-aws-s3-bucket/main.tf",
	} {
		assert.NoFileExists(t, path.Join(workingDirectory, moduleName))
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
tf-aws-vpn-gateway:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpn-gateway"
  version: "v3.2.0"
  destinations:
    - testdata/networking
tf-aws-iam:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-iam"
  version: "v5.11.1"
  destinations:
    - testdata/iam
tf-aws-s3-bucket:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-s3-bucket"
  version: "v3.6.1"
  destinations:
    - testdata/networking
    - testdata/onboarding
    - testdata/some-other-stack
`
	createFile(t, path.Join(folder, "Terrafile"), yaml)
}

func createTerrafile2(t *testing.T, folder string) {
	var yaml = `tf-aws-vpc:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: "v1.46.0"
tf-aws-vpn-gateway:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpn-gateway"
  version: "v3.2.0"
  destinations:
    - testdata/networking
    - testdata/some-other-stack
tf-aws-iam:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-iam"
  version: "v5.11.1"
  destinations:
    - testdata/iam
tf-aws-s3-bucket:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-s3-bucket"
  version: "v3.6.1"
  destinations:
    - testdata/onboarding
`
	createFile(t, path.Join(folder, "Terrafile2"), yaml)
}
