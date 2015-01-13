#!/bin/bash

./taktuk -s -o connector -o status -f hostfile broadcast exec [ killall listener ]
