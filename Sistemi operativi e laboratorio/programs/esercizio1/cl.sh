#!/bin/bash

if [ $# -ge 2] 
then
    FIFO_CMD=$1
    FIFO_OUT=$2

    while test 1 
    do
        read -p "insesrisci comando:" CMD
        echo "$CMD" > "FIFO_CMD"
        if test "${CMD,,}" = "exit"
        then
            break
        fi
        cat < "$FIFO_OUT"
        exit 0
    done

else
    echo "Usage:`basename $0` nome_fifo_cmd nome_fifo_out" 1>&2 
    exit 1
fi 