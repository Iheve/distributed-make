# distributed-make

[![Build Status](https://travis-ci.org/Iheve/distributed-make.svg?branch=master)](https://travis-ci.org/Iheve/distributed-make)

## Prerequisite

You should consider adding `GOPATH=~/go; export $GOPATH` to your .bashrc

## Install

Run `go get github.com/Iheve/distributed-make/...`.

If the previous command did not work, you can try the following:
* Check that your $GOPATH is set. (If not, run `mkdir ~/go; GOPATH=~/go; export $GOPATH`)
* Run `mkdir -p $GOPATH/src/github.com/Iheve`
* Run `cd $GOPATH/src/github.com/Iheve`
* Run `git clone git@github.com:Iheve/distributed-make.git`

## Build

Run `go install github.com/Iheve/distributed-make/...`


## Configuration

### Hostfile for the client

The client needs a list of the servers (listeners). The list looks like this:
```
ensipc101:4242
ensipc100
ensipc102:4243
```
Note that if the default port is 4242.

The default hostfile is `hostfile.cfg`. However, it is convenient to put it in
`/tmp/hostfile.cfg` and use the flag `--hostfile /tmp/hostfile.cfg`



### Hostfile for taktuk

We use taktuk to deploy the listeners on several computers.
Taktuk needs a hostfile too. It looks like this:
```
ensipc100
ensipc101
ensipc102
```

Note that if you are running all the listeners on port 4242, you can use the
same hostfile for taktuk and the client.

### Generating hostfiles

Some scripts will help you to generate config files:

* Use the get_host.py script in folder script/
* Basic use is something like:
```
./get_host.py -f /tmp/hostfile
```
You will generate two files : /tmp/hostfile and /tmp/hostfile.cfg

* You can choose the range of Ensimag PC:
```
./get_host.py -a 20 -b 100
```
* You can choose a list of port for the client connection:
```
./get_host.py -p 4242,4343,4444
```
The list of ports need to be separate with comma.

* You can combine all of these options in order to get what you want

```
./get_host.py -h
```
Can help you ;-)

## Run

### On a single machine

* Don't forget to keep the binaries up to date : `go install github.com/Iheve/distributed-make/...`
* Launch a server with `$GOPATH/bin/listener`
* Create the hostfile.cfg with `echo localhost > /tmp/hostfile.cfg`
* Launch the client with `$GOPATH/bin/client --hostfile /tmp/hostfile.cfg`

You can explore the options by using the --help flag.

### Deploy on several computers at the Ensimag

* Go to the `script/` folder
* Run `./start.sh 100 200` (it will run the listener on ensipc100-ensipc200)
* In an other terminal, launch the client with `$GOPATH/bin/client --hostfile /tmp/hostfile.cfg` in
    the folder where the makefile is

To stop the servers, run `./stop.sh`

## Taktuk deployment

For more fancy deployment...
* Deploy with taktuk:
```
./taktuk -s -o connector -o status -f /tmp/hostfile broadcast exec [ ~/go/bin/listener ]
```
