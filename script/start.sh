#!/bin/bash

./get_host.py -a $1 -b $2
mv hostfile.cfg /tmp
./taktuk -s -o connector -o status -f hostfile broadcast exec [ $GOPATH/bin/listener ]
