#!/bin/bash

if test $# -ge 2
then
	FIFO_CMD="$1"
	FIFO_OUT="$2" 
	
	# Mostriamo tre diverse sintassi di "if" equivalenti
	# if test -p "$FIFO_CMD" -a -p "$FIFO_OUT"
	# if [ -p "$FIFO_CMD" -a -p "$FIFO_OUT" ]
	if [[ -p "$FIFO_CMD" && -p "$FIFO_OUT" ]] # verifica che entrambe siano named pipe
	then
		while test true
		do
			read -p "Inserisci il comando: " CMD # legge l'input dell'utente da testiera e con -p mostra un messaggio di prompt
			echo "$CMD" > "$FIFO_CMD" # scrive il comando sulla pipe di input, così da farlo leggere al server
			if test "${CMD,,}" = "exit" # ${CMD,,} viene sostituito col valore della variabile convertita in minuscolo
			then
				break
			fi
			cat "$FIFO_OUT" # legge l'output del comando dalla pipe di output e lo mostra a video (bloccante finché non arriva qualcosa)
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
