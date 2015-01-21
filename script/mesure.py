#! /usr/bin/env python

import os
import subprocess
from collections import defaultdict

#nbthread = [1,2,4,8,16,32]
nbmachine = [1]
nbthread = [1,2,3,4]
#hosts = ["ensipc101", "ensipc103", "ensipc102", "ensipc104", "ensipc105", "ensipc106",
#			"ensipc107", "ensipc108", "ensipc109",  "ensipc110", "ensipc111", "ensipc112"]

hosts = ["localhost"]
start = 101
fname = "mesures.csv"

fmakefile = "../makefiles/premier"

script_repository = os.path.dirname(os.path.realpath(__file__))
cmd_start = "./start.sh"
cmd_stop = "./stop.sh"

result = defaultdict(list)

def moyenne(dict):
	result = []
	for k, v in dict:
		result.append([k,min(v)])

	return result

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

		for j in range(1,11):
			print("Calculous for " + str(j))

			# Make clean
			print("make clean")

			p = subprocess.Popen("make clean" , shell=True, stdout=subprocess.PIPE, 
										stderr=subprocess.PIPE, cwd=fmakefile)
			p.wait()

			# Launch the client
			full_cmd_client = "export TIMEFORMAT=%E; time $GOPATH/bin/client --hostfile /tmp/hosts --makefile Makefile-verysmall -nbthread " + str(i)

			print(full_cmd_client)

			p = subprocess.Popen(full_cmd_client , shell=True, stdout=subprocess.PIPE, 
										stderr=subprocess.PIPE, cwd=fmakefile)
			out, err = p.communicate()

			#p = subprocess.Popen("export TIMEFORMAT=%E; time make -f Makefile-small -j4" , shell=True, stdout=subprocess.PIPE, 
			#							stderr=subprocess.PIPE, cwd=fmakefile)
			#out, err = p.communicate()

			print("time : " + str(err.split("\n")[-2]))
			print("nb thread reel : " + out.strip())

			result[out.strip()].append(float(err.split("\n")[-2]))

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
