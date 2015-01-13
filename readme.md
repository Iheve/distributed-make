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

## Deploy
* `go install github.com/Iheve/distributed-make/listener`
* List hosts in `/tmp/hosts`, something like:
```
ensipc100
ensipc101
ensipc102
```
* deploy with taktuk:
```
./taktuk -s -o connector -o status -f /tmp/hosts broadcast exec [ ~/go/bin/listener ]
```
* List listeners in `/tmp/hostfile`, something like:
```
ensipc100:4242
ensipc101:4242
ensipc102:4242
```
* Run the client with `$GOHOME/bin/client --hostfile /tmp/hostfile`

## Generate config host file
* Use the get_host.py script in folder script/
* Basic use is something like:
```
./get_host.py -f /tmp/hostfile
```
You will generate two files : /tmp/hostfile and /tmp/hostfile.cfg

* You can choose the range of EnsimagPC:
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
