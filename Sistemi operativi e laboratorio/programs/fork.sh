#!/bin/bash #mi dice quale shell (in questo caso la Bash) per eseguire lo script

RET=0 #contiene il codice di uscita dello script (0 va tutto bene altrimenti ci sono errori)
LABEL="${1}" #non mettere spazi tra le parole, altrimenti da errore
#$1 indica il primo argomento passato allo script e LABEL assume come valoroe tale argomento
if test $# -ge 1 #se viene passato almeno un argomento (-ge = >=)
then

    gcc -o $LABEL $LABEL.c #compila
    RET=$? #ciò che gcc ritorna alla shell ce lo ho in $? (0 o valore diverso da 0)

    if test $RET -eq 0 #-eq significa equal e si usa per i valori numerici
    then #se la compilazione è andata a buion fine

        ./$LABEL

        echo "Pid dello script = $$" # PID del processo principale dello script
        echo "Pid della bash = $BASHPID" # PID della shell corrente
    fi #fine if
else
    echo "Usage: `basename $0` nome_file" >&2 # stampa su stderr il nome dello script senza path
    RET=1
fi

exit $RET 

#N.B. il comando test mi serve per valutare una condizione e restituire vero o falso
