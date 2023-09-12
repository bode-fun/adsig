# AD-Sig

## Description

AD-Sig is a CLI tool to automate the process of generating email signatures for
users via the Active Directory.

## Attention

**Currently, this project is not actively developed, because I have no use-case for it.**

## Quick Start

```sh
$ git clone git@github.com:bode-fun/adsig.git
$ cd adsig
$ go build ./cmd/adsig
$ cp adsig.example.yml adsig.yml # and modify it
$ ./adsig --help
```

Change the configuration file and templates to fit your needs.
Since this project is in development, the templates will need to be in a folder 
called `templates` relative to the current working directory.
The same is true for the config file.

Once, I start the development again, I will use the default OS-specific config and data locations.

For more information on configuration, take a look at the [example file](./adsig.example.yml) or the [relevant code](./config/config.go).

To read more about plans, take a look at the [TODO](./TODO.md) file.

## Contribution

Currently, I do not take contributions, because I am likely to change the licensing
to the [Fair Source License](https://fair.io).

## Usage

Since Germany has no concept of public domain, this code is still under my copyright.
I hereby forbid the usage and modification of the code and from it resulting executables.
