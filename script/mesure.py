#! /usr/bin/env python

import os
import subprocess

#nbthread = [1,2,4,8,16,32]
nbthread = [1,2,3,4]
hosts = ["ensipc101", "ensipc103", "ensipc102", "ensipc104", "ensipc105", "ensipc106",
			"ensipc107", "ensipc108", "ensipc109",  "ensipc110", "ensipc111", "ensipc112"]
start = 101
fname = "mesures.csv"

fmakefile = "../makefiles/premier"

script_repository = os.path.dirname(os.path.realpath(__file__))
cmd_start = "./start.sh"
cmd_stop = "./stop.sh"


for i in nbthread:
	list_hosts = hosts[:i]

	print("Calculating for nb thread : " + str(i))
	print("Will using computers : " + str(list_hosts))

	file_out = open("/tmp/hosts", "wb")
	for c in list_hosts:
		file_out.write(c)
		file_out.write("\n")
	file_out.close()
	
	# Start.sh
	#full_cmd_start = cmd_start + " -a " + str(start) + " -b " + str(start + i - 1)
	#print(full_cmd_start)

	#p = subprocess.Popen(full_cmd_start , shell=True, stdout=subprocess.PIPE, 
	#							stderr=subprocess.PIPE, cwd=script_repository)
	#p.wait()

	for j in range(1,5):
		print("Calculous for " + str(j))
		# Make clean
		print("make clean")

		p = subprocess.Popen("make clean" , shell=True, stdout=subprocess.PIPE, 
									stderr=subprocess.PIPE, cwd=fmakefile)
		p.wait()

		# Launch the client
		full_cmd_client = "export TIMEFORMAT=%E; time $GOPATH/bin/client --hostfile /tmp/hosts --makefile Makefile-verysmall"
		
		print(full_cmd_client)

		p = subprocess.Popen(full_cmd_client , shell=True, stdout=subprocess.PIPE, 
									stderr=subprocess.PIPE, cwd=fmakefile)
		out, err = p.communicate()


		
		#p = subprocess.Popen("export TIMEFORMAT=%E; time make -f Makefile-small -j4" , shell=True, stdout=subprocess.PIPE, 
		#							stderr=subprocess.PIPE, cwd=fmakefile)
		#out, err = p.communicate()

		print("time : " + str(err.split("\n")[-2]))
		print(out)