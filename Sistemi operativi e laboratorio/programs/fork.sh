LABEL="fork" #non mettere spazi tra le parole, altrimenti da errore

gcc -o $LABEL $LABEL.c
RET=$? #ciò che gcc ritorna alla shell ce lo ho in $?

if test $RET -eq 0 #-eq significa equal e si usa per i valori numerici
then

    ./$LABEL

    echo "Pid dello script = $$"
    echo "Pid della bash = $BASHPID"
fi #fine if

exit $RET 