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
Terrafile expects a file named `Terrafile` which will contain your terraform module dependencies in a yaml like format.

An example Terrafile:
```
tf-aws-vpc:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    version: "v1.46.0"
tf-aws-vpc-experimental:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    version: "master"
```

Terrafile config file in current directory and modules exported to ./vendor/modules
```sh
$ terrafile
INFO[0000] [*] Checking out v1.46.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc  
INFO[0000] [*] Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc  
```

Terrafile config file in custom directory
```sh
$ terrafile -f config/Terrafile
INFO[0000] [*] Checking out v1.46.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc  
INFO[0000] [*] Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc  
```

Terraform modules exported to custom directory
```sh
$ terrafile -p custom_directory
INFO[0000] [*] Checking out master of git@github.com:terraform-aws-modules/terraform-aws-vpc  
INFO[0001] [*] Checking out v1.46.0 of git@github.com:terraform-aws-modules/terraform-aws-vpc  
```

## TODO
* Break out the main logic into seperate commands (e.g. version, help, run)
* Update tests to include unit tests for broken out commands
* Add coverage tool and badge
* May be worth renaming Terrafile config file to something that won't be misinterpreted as the binary