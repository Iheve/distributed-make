#!/bin/bash

./get_host.py -a 10 -b 300
mv hostfile.cfg /tmp
./taktuk -s -o connector -o status -f hostfile broadcast exec [ $GOPATH/bin/listener ]
