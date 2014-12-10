# distributed-make

## Prerequisite

You should consider adding `GOPATH=~/go; export $GOPATH` to your .bashrc

## Install

* Check that your $GOPATH is set. (If not, run `mkdir ~/go; GOPATH=~/go; export $GOPATH`)
* Run `mkdir -p $GOPATH/src/github.com/Iheve`
* Run `cd $GOPATH/src/github.com/Iheve`
* Run `git clone git@github.com:Iheve/distributed-make.git`

## Build

* `go install github.com/Iheve/distributed-make/listener`
* `go install github.com/Iheve/distributed-make/client`

## Run

* launch the server with `$GOPATH/bin/listener`
* launch the client with `$GOPATH/bin/client`
