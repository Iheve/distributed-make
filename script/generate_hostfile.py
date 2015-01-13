#! /usr/bin/env python3

import sys
import argparse

parser = argparse.ArgumentParser(description='Scripting chain for DistributedMakefile')
parser.add_argument("-i", "--in", dest="hostfile_in", default="hostfile",
                  help="Name of the output file")
parser.add_argument("-o", "--out", dest="hostfile_out", default="hostfile.cfg",
                  help="Name of the output file")
parser.add_argument("-p", "--ports", dest="ports", default=4242,
                  help="List of ports to use on Ensimag PC : 4242,4343,4344")

args = parser.parse_args()
print(args)

list_ports = str(args.ports).split(",")

print("List of ports " + str(list_ports))

file_in = open(args.hostfile_in, "rb")
file_out = open(args.hostfile_out, "wb")

for ligne in file_in:
	for port in list_ports:
		hostname = ligne.strip()
		file_out.write(hostname)
		file_out.write(bytes(":", 'UTF-8'))
		file_out.write(bytes(port, 'UTF-8'))
		file_out.write(bytes("\n", 'UTF-8'))

file_in.close()
file_out.close()



