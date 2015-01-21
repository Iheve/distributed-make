#! /usr/bin/env python

import os
import subprocess
from collections import defaultdict

#nbthread = [1,2,4,8,16,32]
nbmachine = [1,2]
nbthread = [1,2]
hosts = ["ensipc150", "ensipc151", "ensipc144", "ensipc145", "ensipc149", "ensipc153", "ensipc142"]

#hosts = ["localhost"]
start = 101
fname = "mesures.csv"
ftmp = "tmp.csv"

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
