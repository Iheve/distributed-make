#! /usr/bin/env python

import argparse
import subprocess
import os

parser = argparse.ArgumentParser(description='Scripting chain for DistributedMakefile')
parser.add_argument("-a", "--a", dest="min", default=10,
                  help="minimum range value for Ensimag PC ")
parser.add_argument("-b", "--b", dest="max", default=100,
                  help="max range value for Ensimag PC ")
parser.add_argument("-p", "--ports", dest="ports", default=4242,
                  help="List of ports to use on Ensimag PC")
parser.add_argument("-f", "--file", dest="file", default="hostfile",
                  help="Name of the output file")

args = parser.parse_args()
print(args)

script_repository = os.path.dirname(os.path.realpath(__file__))

cmd = "./pimag.sh " + str(args.min) + " " + str(args.max) + " > " + str(args.file)

print("Generating hostfile file")

p = subprocess.Popen(cmd , shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, cwd=script_repository)
p.wait()

ports = str(args.ports).split(",")

file_in = args.file
file_out = args.file + ".cfg"

cmd = "./generate_hostfile.py -p " + str(args.ports) + " -i " + file_in + " -o " + file_out


print("Generating hostfile.cfg file")

p = subprocess.Popen(cmd , shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, cwd=script_repository)
p.wait()