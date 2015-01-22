#! /usr/bin/env python

import datetime
import os
import subprocess
import argparse
from collections import defaultdict

# Parsing args
parser = argparse.ArgumentParser(description='Scripting chain launching automatic tasks')
parser.add_argument("-m", "--makefile", dest="fmakefile", default="../makefiles/premier",
                  help="Relative path to Makefile directory")
args = parser.parse_args()

fmakefile = args.fmakefile

#nbthread = [1,2,4,8,16,32]
nbmachine = [1,2]
nbthread = [1,4]
hosts = ["localhost"]
#hosts = ["ensipc150", "ensipc151", "ensipc144", "ensipc145", "ensipc149", "ensipc153", "ensipc142"]

# Date format
now = datetime.datetime.now()
date_format = now.strftime("%Y%m%d_%H%M%S")

# Create subdirectory to save files
directory = "mesures"
if not os.path.exists(directory):
    os.makedirs(directory)

# File name to save mesures
fname = "mesures/mesures_" + str(date_format) + ".csv"
ftmp = "mesures/mesures_full_" + str(date_format) + ".csv"

# Path to script repository
script_repository = os.path.dirname(os.path.realpath(__file__))

result = defaultdict(list)


def moyenne(dict):
	result = []
	for k, v in dict:
		result.append([k,min(v)])

	return result

tmp_out = open(ftmp, "wb")
tmp_out.write("nthreads;temps;nmachine;nrealthreads\n")
tmp_out.close()

for y in nbmachine:
	list_hosts = hosts[:y]
	print("Calculating with " + str(y) + " computers")
	print("Will use computers : " + str(list_hosts))

	file_out = open("/tmp/hosts", "wb")
	for c in list_hosts:
		file_out.write(c)
		file_out.write("\n")
	file_out.close()

	for i in nbthread:
		print("Calculating for nb thread : " + str(i*y))

		for j in range(1,3):
			print("Loop iteration n" + str(j))

			# Make clean
			print("Executing : make clean")

			p = subprocess.Popen("make clean" , shell=True, stdout=subprocess.PIPE, 
										stderr=subprocess.PIPE, cwd=fmakefile)
			p.wait()

			# Launch the client
			full_cmd_client = "export TIMEFORMAT=%E; time $GOPATH/bin/client --hostfile /tmp/hosts --makefile Makefile -nbthread " + str(i)

			print("Executing : " + full_cmd_client)

			p = subprocess.Popen(full_cmd_client , shell=True, stdout=subprocess.PIPE, 
										stderr=subprocess.PIPE, cwd=fmakefile)
			out, err = p.communicate()

			print("Real time : " + str(err.split("\n")[-2]))
			print("NbThread used : " + out.strip())

			# Store data in a dictionnary
			result[out.strip()].append(float(err.split("\n")[-2]))

			# Write in tmp file
			tmp_out = open(ftmp, "a")
			tmp_out.write(out.strip())
			tmp_out.write(";")
			tmp_out.write(err.split("\n")[-2])
			tmp_out.write(";")
			tmp_out.write(str(y))
			tmp_out.write(";")
			tmp_out.write(str(i))
			tmp_out.write("\n")
			tmp_out.close()


tmp_out.close()

# Calculting average on all the values for nb threads
final = moyenne(result.items())
print(final)
file_out = open(fname, "wb")
file_out.write("nthreads;temps\n")
for c in final:
	file_out.write(c[0])
	file_out.write(";")
	file_out.write(str(c[1]))
	file_out.write("\n")
file_out.close()
