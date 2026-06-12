#!/bin/bash

chmod 777 server.c client.c executor.c # concede ai file tutti i permessi di lettura, scrittura ed esecuzione

gcc -o server server.c
gcc -o client client.c
gcc -o executor executor.c

echo "Compilazione avvenuta con successo"

./server