# Terrafile [![Build Status](https://circleci.com/gh/coretech/terrafile.svg?style=shield)](https://circleci.com/gh/coretech/terrafile)

Terrafile is a binary written in Go to systematically manage external modules from Github for use in Terraform. See this [article](http://bensnape.com/2016/01/14/terraform-design-patterns-the-terrafile/) for more information on how it was introduced in a Ruby rake task.

## How to install

### macOS

```sh
brew tap coretech/terrafile && brew install terrafile
```

### Linux
Download your preferred flavor from the [releases](https://github.com/coretech/terrafile/releases/latest) page and install manually.

For example:
```sh
curl -L https://github.com/coretech/terrafile/releases/download/v{VERSION}/terrafile_{VERSION}_Linux_x86_64.tar.gz | tar xz -C /usr/local/bin
```

## How to use
Terrafile expects a file named `Terrafile` which will contain your terraform module dependencies in a yaml format.

An example Terrafile:
```
tf-aws-vpc:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    version: "v1.46.0"
tf-aws-vpc-experimental:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    version: "master"
tf-aws-vpc-default:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
tf-aws-vpc-commit:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    commit: "01601169c00c68f37d5df8a80cc17c88f02c04d0"
tf-aws-vpc-path:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    path: "examples/simple-vpc"
```
By default, Terrafile will checkout the master branch of a module.

To checkout a tag or branch use the `version` key:
```
tf-aws-vpc:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    version: "v1.46.0"
```

To pin a module to a specific `commit`:
```
tf-aws-vpc-commit:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    commit: "01601169c00c68f37d5df8a80cc17c88f02c04d0"
```

Many organizations use a mono-repo approach with all their Terraform modules
defined in one repository. If you're in a position where all of your
Terraform modules are in a single repo but you only require some of them, you
can checkout a path from a Git repository.

To checkout a `path` from a repository:
```
tf-aws-vpc-path:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    path: "examples/simple-vpc"
```


Terrafile config file in current directory and modules exported to ./vendor/modules
```sh
$ terrafile
INFO[0000] [tf-aws-vpc-path] Checking out examples/simple-vpc from master of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0002] [tf-aws-vpc] Checking out v1.46.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0003] [tf-aws-vpc-experimental] Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0005] [tf-aws-vpc-default] Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0007] [tf-aws-vpc-commit] Checking out 01601169c00c68f37d5df8a80cc17c88f02c04d0 of git@github.com:terraform-aws-modules/terraform-aws-vpc
```

Terrafile config file in custom directory
```sh
$ terrafile -f config/Terrafile
INFO[0000] [tf-aws-vpc-path] Checking out examples/simple-vpc from master of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0002] [tf-aws-vpc] Checking out v1.46.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0003] [tf-aws-vpc-experimental] Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0005] [tf-aws-vpc-default] Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0007] [tf-aws-vpc-commit] Checking out 01601169c00c68f37d5df8a80cc17c88f02c04d0 of git@github.com:terraform-aws-modules/terraform-aws-vpc
```

Terraform modules exported to custom directory
```sh
$ terrafile -p custom_directory
INFO[0000] [tf-aws-vpc-path] Checking out examples/simple-vpc from master of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0002] [tf-aws-vpc] Checking out v1.46.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0003] [tf-aws-vpc-experimental] Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0005] [tf-aws-vpc-default] Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc
INFO[0007] [tf-aws-vpc-commit] Checking out 01601169c00c68f37d5df8a80cc17c88f02c04d0 of git@github.com:terraform-aws-modules/terraform-aws-vpc
```

## TODO
* Break out the main logic into seperate commands (e.g. version, help, run)
* Update tests to include unit tests for broken out commands
* Add coverage tool and badge
* May be worth renaming Terrafile config file to something that won't be misinterpreted as the binary
