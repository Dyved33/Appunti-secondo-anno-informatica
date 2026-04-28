#!/bin/bash

RET=0

# La variabile speciale "$#" contiene il numero degli argomenti
# passati allo script
if test $# -ge 1
then
	LABEL="${1}" 	# Quando ne referenziamo il valore possiamo racchiudere il nome di una variabile tra graffe
					# Permette di evitare ambiguità
	
	gcc -o "$LABEL" "$LABEL".c
	RET=$? # "$?" contiene il valore di ritorno dell'ultimo comando eseguito

	if test $RET -eq 0
	then
		"./$LABEL" & #   &  per l'esecuzione in background. Si è inserito a titolo di esempio
		
		wait 	# comando interno che attende la terminazione dei child dello script
				# Si è inserito a titolo di esempio a seguito dell'esecuzione in background
		
		echo "Pid dello script = $$"
		echo "Pid della bash   = $PPID"
	fi
else
	# 1>&2  duplica il descrittore 2 (stderr) ed il duplicato è 1. Ciò che era aperto con 1 (stdout) viene chiuso.
	# In altre parole  1>&2 manda stdout su ciò che è aperto col descrittore 2
	echo "Usage: `basename $0` nome_file"   1>&2
	RET=1
fi

exit $RET