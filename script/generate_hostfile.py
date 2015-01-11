#! /usr/bin/env python

import sys

list_ports = sys.argv[1:]
HOST_FILE_IN = "hostfile"
HOST_FILE_OUT = "hostfile.cfg"

print("List of args " + str(list_ports))

file_in = open(HOST_FILE_IN, "rb")
file_out = open(HOST_FILE_OUT, "wb")

for ligne in file_in:
	for port in list_ports:
		hostname = ligne.strip()
		file_out.write(hostname + bytes(":", 'UTF-8') +
			bytes(port, 'UTF-8') + bytes('\n', 'UTF-8'))

file_in.close()
file_out.close()



