# Terrafile [![Build Status](https://circleci.com/gh/coretech/terrafile.svg?style=shield)](https://circleci.com/gh/coretech/terrafile)

Terrafile is a binary written in Go to systematically manage external modules from Github for use in Terraform. See this [article](http://bensnape.com/2016/01/14/terraform-design-patterns-the-terrafile/) for more information on how it was introduced in a Ruby rake task.

This is currently an experimental WIP.

An example Terrafile:
```
tf-aws-vpc:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    version: "v1.46.0"
tf-aws-vpc-experimental:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    version: "master"
```

## How to install

### macOS

```sh
brew tap coretech/terrafile && brew install terrafile
```

### Linux
Download your preferred flavor from the [releases](https://github.com/coretech/terrafile/releases/latest) page and install manually.
```sh
curl -L https://github.com/coretech/terrafile/releases/download/v0.2/terrafile_0.2_linux_amd64.tar.gz | tar xz -C /usr/local/bin
```