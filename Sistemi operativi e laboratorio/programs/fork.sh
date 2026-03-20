#!/bin/bash

RET=0
LABEL="${1}" #non mettere spazi tra le parole, altrimenti da errore

if test $# -ge 1
then

    gcc -o $LABEL $LABEL.c
    RET=$? #ciò che gcc ritorna alla shell ce lo ho in $?

    if test $RET -eq 0 #-eq significa equal e si usa per i valori numerici
    then

        ./$LABEL

        echo "Pid dello script = $$"
        echo "Pid della bash = $BASHPID"
    fi #fine if
else
    echo "Usage: `basename $0` nome_file" >&2
    RET=1
fi

exit $RET 