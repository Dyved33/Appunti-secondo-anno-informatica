#!/bin/bash

if test $# -ge 2
then
	FIFO_CMD="$1"
	FIFO_OUT="$2" 
	
	####### Mostriamo tre diverse sintassi di "if" equivalenti
	#if test -p "$FIFO_CMD" -a -p "$FIFO_OUT"
	#if [ -p "$FIFO_CMD" -a -p "$FIFO_OUT" ]
	if [[ -p "$FIFO_CMD" && -p "$FIFO_OUT" ]]
	then
		while test true
		do
			read -p "Inserisci il comando: " CMD
			echo "$CMD" > "$FIFO_CMD"
			if test "${CMD,,}" = "exit"
			then
				break
			fi
			cat "$FIFO_OUT"
		done
		
		exit 0
	else
		echo "Fifo non trovata/e!" >&2
		exit 1
	fi
else
	echo "Usage: `basename $0`  fifo_cmd_name  fifo_out_name"  >&2
	exit 1
fi
