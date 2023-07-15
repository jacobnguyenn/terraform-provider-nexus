<img src="assets/logo.png" alt="logo" width="auto" height="600" />

# Terraform provider Nexus

![codeql workflow](https://github.com/datadrivers/terraform-provider-nexus/actions/workflows/codeql-analysis.yml/badge.svg)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/datadrivers/terraform-provider-nexus)](https://goreportcard.com/report/github.com/datadrivers/terraform-provider-nexus)

- [Terraform provider Nexus](#terraform-provider-nexus)
  - [Introduction](#introduction)
  - [Usage](#usage)
    - [Provider config](#provider-config)
  - [Development](#development)
    - [Build](#build)
    - [Testing](#testing)
    - [Create documentation](#create-documentation)
  - [Author](#author)

## Introduction
This repo is forked from https://github.com/datadrivers/terraform-provider-nexus.

What's included:
- Bug fixes:
  - [provider crash when negative_cache block is not included](https://github.com/datadrivers/terraform-provider-nexus/issues/352)
  - [order in repository group member is not respected](https://github.com/datadrivers/terraform-provider-nexus/issues/339)
- Feature:
  - support creating [cleanup policy](https://github.com/nduyphuong/terraform-provider-nexus/pull/5)



- Terraform provider to configure Sonatype Nexus using its API.

Implemented and tested with Sonatype Nexus `3.43.0-01`.

## Usage

### Provider config

```hcl
terraform {
  required_version = ">= 0.14"

  required_providers {
    nexus = {
      source = "nduyphuong/nexus"
    }
  }
}
provider "nexus" {
  insecure = true
  password = "admin123"
  url      = "https://127.0.0.1:8080"
  username = "admin"
}
```

## Development

### Build

There is a [makefile](./GNUmakefile) to build the provider and place it in repos root dir.

```sh
make
```

To use the local build version you need tell terraform where to look for it via a terraform config override.

Create `dev.tfrc` in your terraform code folder (f.e. in [dev.tfrc](./examples/local-development/dev.tfrc)):

```hcl
# dev.tfrc
provider_installation {

  # Use /home/developer/tmp/terraform-nexus as an overridden package directory
  # for the datadrivers/nexus provider. This disables the version and checksum
  # verifications for this provider and forces Terraform to look for the
  # nexus provider plugin in the given directory.
  # relative path also works, but no variable or ~ evaluation
  dev_overrides {
    nduyphuong/nexus = "../../"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Tell your shell environment to use override file:

```bash
export TF_CLI_CONFIG_FILE=dev.tfrc
```

Now run your terraform commands (`plan` or `apply`), `init` is ***not*** required.

```bash
# start local nexus
make start-services
# run local terraform code
cd examples/local-development
terraform plan
terraform apply
```

### Testing

**NOTE**: For testing Nexus Pro features, place the `license.lic` in `scripts/`.

For testing start a local Docker containers using make

```shell
make start-services
```

This will start a Docker and MinIO containers and expose ports 8081 and 9000.

Now start the tests

```shell
make testacc
```

or skipped tests:

```shell
SKIP_S3_TESTS=1 make testacc
SKIP_AZURE_TESTS=1 make testacc
SKIP_PRO_TESTS=1 make testacc
```

#### Debug the provider with [delve](https://github.com/go-delve/delve)
```bash
 dlv exec --accept-multiclient --continue --headless ./terraform-provider-nexus -- -debuggable
 export TF_REATTACH_PROVIDERS=.....
```

#### To debug tests

Set env variable `TF_LOG=DEBUG` to see additional output.

Use `printState()` function to discover terraform state (and resource props) during test.

Debug configurations are also available for VS Code.

### Create documentation

To generate the terraform documentation from go files, you can run

```shell
make docs
```

## Author

Alot of thanks go to [Datadrivers GmbH](https://www.datadrivers.de)
