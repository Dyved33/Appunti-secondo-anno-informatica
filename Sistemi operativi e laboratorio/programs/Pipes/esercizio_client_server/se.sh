#!/bin/bash

if [ $# -ge 2 ] # Se il numero di argomenti è maggiore o uguale a 2
then
	# Utilizziamo "realpath" per ottenere il path assoluto delle due Named Pipe.
	# Questo perché se il client richiedesse di cambiare la directory corrente
	# il server potrebbe non trovare più le pipes
	FIFO_CMD=`realpath "$1"`
	FIFO_OUT=`realpath "$2"`
	
	if test ! -p "$FIFO_CMD" # -p verifica se si tratti di una named pipe
	then
		mkfifo "$FIFO_CMD"  # creazione named pipe
	fi
	if test ! -p "$FIFO_OUT"
	then
		mkfifo "$FIFO_OUT"
	fi	

	while test 1 # sempre vero
	do
		CMD=`cat "$FIFO_CMD"` # legge un comando dalla pipe (bloccante finché non arriva qualcosa)
		if test "${CMD,,}" = "exit"  # ${VAR,,} viene sostituito col valore della variabile convertita
		then                         # in minuscolo
			break
		fi
		#echo "Eseguo \"$CMD\"..."
		eval "$CMD" > "$FIFO_OUT" 2>&1  # vedi "help eval" (comando interno...)
		# eval esegue il comando come shell, l'output viene rediretto sulla pipe di output, così come l'eventuale messaggio di errore (2>&1)
	done

	exit 0
else
	echo "Usage: `basename $0`  fifo_cmd_name  fifo_out_name"  >&2
	exit -1
fi