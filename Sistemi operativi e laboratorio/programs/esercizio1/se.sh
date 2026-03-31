#!/bin/bash

if [ $# -ge 2 ] 
then
    FIFO_CMD=$1
    FIFO_OUT=$2
    if test ! -p $FIFO_CMD
    then
        mkfifo $FIFO_CMD
    fi
    if test ! -p $FIFO_OUT
    then
        mkfifo $FIFO_OUT
    fi    

    while test 1 
    do
        CMD=`cat "$FIFO_CMD"`
        if test "${CMD,,}" = "exit"
        then
            break
        fi
        echo "eseguo \"${CMD}\"..."
        eval "$CMD" > "$FIFO_OUT" 2>&1
    done
    exit 0
else
    echo "Usage: `basename "$0"` nome_fifo_cmd nome_fifo_out" 1>&2 
    exit 1
fi 