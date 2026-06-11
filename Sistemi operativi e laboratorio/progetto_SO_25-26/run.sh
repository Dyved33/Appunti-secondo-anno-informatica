#!/bin/bash

chmod 777 server.c client.c executor.c

gcc -o server server.c
gcc -o client client.c
gcc -o executor executor.c

echo "Compilazione avvenuta con successo"

./server