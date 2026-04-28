#!/bin/bash

###############################################################
### Definizione variabili
SC_C_FILENAME="compara_sc.c"
STDIO_C_FILENAME="compara_stdio.c"
SC_EXE_FILENAME="compara_sc"
STDIO_EXE_FILENAME="compara_stdio"
SC_OUT_FILE="sc.out"
STDIO_OUT_FILE="stdio.out"

###############################################################
### Recupero argomenti da linea di comando
### La variabile speciale $# contiene il numero di argomenti passati allo script dalla linea di comando
if test $# -lt 2
then
	echo "Usage: `basename $0`  write_count  Buffer_size"  >&2
    exit -1
fi

### Gli argomenti passati allo script sono disponibili nelle variabili speciali $1, $2,..., $n
COUNT="$1"
SIZE="$2"

###############################################################
### Compilazione eseguibili
gcc -o $SC_EXE_FILENAME $SC_C_FILENAME
### La variabile speciale $? contiene il valore di ritorno dell'ultimo comando eseguito
if test $? -ne 0
then
	exit -1
fi
gcc -o $STDIO_EXE_FILENAME $STDIO_C_FILENAME
if test $? -ne 0
then
	exit -1
fi

###############################################################
### Confronto del consumo di CPU
echo -e "\nCOUNTING REAL AND CPU TIME..."
time ./$STDIO_EXE_FILENAME $STDIO_OUT_FILE $COUNT $SIZE

echo
time ./$SC_EXE_FILENAME $SC_OUT_FILE       $COUNT $SIZE

###############################################################
### Confronto del numero di chiamate a write(2)
echo -e "\n#####################################################"
echo "Std/IO library: counting write() calls..."
strace -e write ./$STDIO_EXE_FILENAME $STDIO_OUT_FILE   $COUNT $SIZE  2>&1 | grep ^write\(3 | wc -l

echo
echo "System Calls: counting write() calls..."
strace -e write ./$SC_EXE_FILENAME $SC_OUT_FILE   $COUNT $SIZE  2>&1 | grep ^write\(3 | wc -l

###############################################################
### Controllo che i due file di output siano uguali
echo -e "\n#####################################################"
ls -lh $SC_OUT_FILE $STDIO_OUT_FILE
ls -l $SC_OUT_FILE $STDIO_OUT_FILE
echo

diff $SC_OUT_FILE $STDIO_OUT_FILE
if test $? -eq 0
then
	echo "$SC_OUT_FILE e $STDIO_OUT_FILE sono uguali"
else
	echo "$SC_OUT_FILE e $STDIO_OUT_FILE sono diversi!"
fi

###############################################################
### Rimozione file di output
rm $SC_OUT_FILE $STDIO_OUT_FILE

exit 0
