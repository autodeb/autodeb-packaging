# autodeb

[![godoc reference](https://godoc.org/salsa.debian.org/autodeb-team/autodeb?status.svg)](http://godoc.org/salsa.debian.org/autodeb-team/autodeb)
[![pipeline status](https://salsa.debian.org/autodeb-team/autodeb/badges/master/pipeline.svg)](https://salsa.debian.org/autodeb-team/autodeb/commits/master)
[![go report card](https://goreportcard.com/badge/salsa.debian.org/autodeb-team/autodeb)](https://goreportcard.com/report/salsa.debian.org/autodeb-team/autodeb)
[![Build Status](https://travis-ci.com/autodeb/autodeb.svg?branch=master)](https://travis-ci.com/autodeb/autodeb)
[![codecov](https://codecov.io/gh/autodeb/autodeb/branch/master/graph/badge.svg)](https://codecov.io/gh/autodeb/autodeb)

autodeb tries to automatically update Debian packages to newer upstream versions or to backport them.

autodeb is the concretization Lucas Nussbaum's GSOC 2018 proposed project titled "Automatic Packages for Everything (backports, new upstream versions, etc.)". The [project proposal](https://wiki.debian.org/SummerOfCode2018/Projects) can be found in the Debian Wiki. The project was [officially accepted](https://summerofcode.withgoogle.com/projects/#5560246244737024).

## Getting in touch

You may chat with us at [#autodeb on irc.debian.org](irc://irc.debian.org:6667/autodeb) (or via [webchat](https://webchat.oftc.net/?channels=autodeb)). If you've found something that is clearly a bug, feel free to report it in the [issue tracker](https://salsa.debian.org/autodeb-team/autodeb/issues).

## Documentation

 - [Getting started](#getting-started): everything you need to know to build autodeb
 - dependency graphs (helpful to understand the software architecture):
   + [cmd/autodeb-server](https://autodeb-team.pages.debian.net/autodeb/dependency-graph-autodeb-server.svg): full dependency graph for autodeb-server
   + [cmd/autodeb-worker](https://autodeb-team.pages.debian.net/autodeb/dependency-graph-autodeb-worker.svg): full dependency graph for autodeb-worker
   + [internal/server](https://autodeb-team.pages.debian.net/autodeb/dependency-graph-server.svg): dependency graph for the server package (easier to read)
   + [internal/worker](https://autodeb-team.pages.debian.net/autodeb/dependency-graph-worker.svg): dependency graph for the worker package (easier to read)
 - [godoc](https://godoc.org/salsa.debian.org/autodeb-team/autodeb): code documentation
 - [Packaging (other repository)](https://salsa.debian.org/autodeb-team/autodeb-packaging/blob/master/debian/README.md): Debian packaging for autodeb
 - [Infrastructure (other repository)](https://salsa.debian.org/autodeb-team/infrastructure): Contains the ansible scripts for auto.debian.net
 - [Wiki](https://salsa.debian.org/autodeb-team/autodeb/wikis/home): everything else

## Available executables

- ``list-packages-with-newer-upstream-versions``: lists source packages that have newer upstream versions available

- ``autodeb-server``: This is the server component of the system. It provides a web interface, a REST API and dput-compatible interface.

- ``autodeb-worker``: This is the worker component of the system. It retrieves jobs from the main server and executes them.

## Getting started

### 1. Setup Go

Note that you might want to get a recent version of the go compiler from a backports repository.

```shell
$ apt-get install golang-go git make
$ export GOPATH=~/go
$ go get -u golang.org/x/lint/golint
```

### 2. Clone the project

```shell
$ mkdir -p $GOPATH/src/salsa.debian.org/autodeb-team/
$ git clone https://salsa.debian.org/autodeb-team/autodeb.git $GOPATH/src/salsa.debian.org/autodeb-team/autodeb
$ cd $GOPATH/src/salsa.debian.org/autodeb-team/autodeb
```

### 3. Build the project

```shell
$ make get-deps
$ make
```

### 4. Run any of the scripts

Note that runtime dependencies of the scripts include:
 - devscripts
 - sbuild

```shell
$ ./list-packages-with-newer-upstream-versions
$ ./autodeb-server
$ ./autodeb-worker
```
