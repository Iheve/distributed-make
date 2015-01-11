#! /bin/bash

#scan les ensipc de a a b.
#augmenter STEP en cas de probleme de fork
a=$1
b=$2
STEP=2

#scan une plage d'ensipc du premier argument au deuxieme argument
function scan
{
	i=$1
	while test $i -le $2;
	do r=$(ping -c 1 -w 1 ensipc$i 2> /dev/null | grep rtt | wc -l);
		if [[ $r == 1 ]]
		then
			echo ensipc$i
		fi
		i=$(($i+1));
	done
}

#lance des scans en parallel
function parallelscan
{
	PIDS=""
	i=$a
	while test $i -le $b
	do
		scan $i $(($i+$STEP))&
		PIDS="$PIDS $!"
		i=$(($i+$STEP+1));
	done

	#on attend que tout ait ete termine avant de rendre la main
	for i in $PIDS
	do
		wait $i
		#echo $i
	done
}

parallelscan
