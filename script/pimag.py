#! /usr/bin/env python3

import argparse
import subprocess
import os
import threading
import sys

def scan(id):
	cmd = "ping -c 1 -w 1 ensipc" + str(id) + " 2> /dev/null | grep rtt | wc -l"

	p = subprocess.Popen(cmd , shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
	out, err = p.communicate()

	if out.strip() == "1".encode():
		print("ensipc" + str(id), flush=True)
		return True
	else:
		return False

parser = argparse.ArgumentParser(description='Scripting for getting all ensimag PC up')
parser.add_argument("-a", "--min", dest="min", default="10",
                  help="Minimum number for Ensimag PC", type=int)
parser.add_argument("-b", "--max", dest="max", default="100",
                  help="Maximum number for Ensimag PC", type=int)
args = parser.parse_args()

for i in range(args.min, args.max):
    a = threading.Thread(None, scan, None, (i,), {})
    a.start()
