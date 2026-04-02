#!/bin/bash

if [ $# -ge 2 ]
then
	# Utilizziamo "realpath" per ottenere il path assoluto delle due Named Pipe.
	# Questo perché se il client richiedesse di cambiare la directory corrente
	# il server potrebbe non trovare più le pipes
	FIFO_CMD=`realpath "$1"`
	FIFO_OUT=`realpath "$2"`
	
	if test ! -p "$FIFO_CMD"
	then
		mkfifo "$FIFO_CMD"  # creazione pipe
	fi
	if test ! -p "$FIFO_OUT"
	then
		mkfifo "$FIFO_OUT"
	fi	

	while test 1
	do
		CMD=`cat "$FIFO_CMD"`
		if test "${CMD,,}" = "exit"  # ${VAR,,} viene sostituito col valore della variabile convertita
		then                         # in minuscolo
			break
		fi
		#echo "Eseguo \"$CMD\"..."
		eval "$CMD" > "$FIFO_OUT" 2>&1  # vedi "help eval" (comando interno...)
	done

	exit 0
else
	echo "Usage: `basename $0`  fifo_cmd_name  fifo_out_name"  >&2
	exit -1
fi