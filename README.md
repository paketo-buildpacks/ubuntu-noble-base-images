# Ubuntu noble base images

## Paketo Noble Base Stack

### What is this stack for?

Ideal for:

- Java apps and .NET Core apps
- Go apps that require some C libraries
- Node.js/Python/Ruby/etc. apps **without** many native extensions

### What's in the build and run images of this stack?

This stack's build and run images are based on Ubuntu Noble Numbat.

- To see the **list of all packages installed** in the build or run image for a given release, see the `noble-base-stack-{version}-build-receipt.cyclonedx.json` and `noble-base-stack-{version}-run-receipt.cyclonedx.json` attached to each [release](https://github.com/paketo-buildpacks/noble-base-stack/releases). For a quick overview of the packages you can expect to find, see the [stack descriptor file](stack/stack.toml).

## Paketo Noble Tiny Stack

### What is this stack for?

Ideal for:

- most Golang apps
- Java [GraalVM Native Images](https://www.graalvm.org/docs/reference-manual/native-image/)

### What's in the build and run images of this stack?

This stack's build image is based on Ubuntu Noble Numbat. Its run image does not include a Linux distribution.

- To see the **list of all packages installed** in the build or run image for a given release,
  see the `noble-tiny-stack-{version}-build-receipt.cyclonedx.json` and `noble-tiny-stack-{version}-run-receipt.cyclonedx.json` attached to each [release](https://github.com/paketo-buildpacks/noble-tiny-stack/releases). For a quick overview of the packages you can expect to find, see the [stack descriptor file](stack/stack.toml).

## Paketo Noble Static Stack

Stack for statically-linked binaries for Ubuntu 24.04: Noble Numbat

## What is a stack?

See Paketo's [stacks documentation](https://paketo.io/docs/concepts/stacks/).

## How can I contribute?

Contribute changes to this stack via a Pull Request. Depending on the proposed changes, you may need to [submit an RFC](https://github.com/paketo-buildpacks/rfcs) first.

## How do I test the stack locally?

Run [`scripts/test.sh`](scripts/test.sh).

## How do I generate package receipts

To generate a package receipt based on existing `build.oci` and `run.oci` archives, use [`scripts/receipts.sh`](scripts/receipts.sh).
