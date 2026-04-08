#!/bin/bash

echo "CIAO DALLO SCRIPT!"
echo "MY_VAR = $MY_VAR"

cat 0<&3 #cat legge dal file con fd 3

exit 0