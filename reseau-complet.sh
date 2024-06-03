#!/bin/bash

# Vérifie si un argument est passé et qu'il s'agit d'un nombre
if [ $# -eq 0 ]; then
    echo "Usage: $0 <number_of_applications>"
    exit 1
fi

NUM_SITES=$1

# CRÉATION DES FIFOS/DES LIENS
for i in $(seq 1 $NUM_SITES); do
  mkfifo /tmp/in_A$i /tmp/out_A$i
  mkfifo /tmp/in_C$i /tmp/out_C$i
done

# LANCEMENT DES SITES
for i in $(seq 1 $NUM_SITES); do
  go run app-base -n A$i -m g -port $((4443 + i)) < /tmp/in_A$i > /tmp/out_A$i &
  go run app-control -n C$i -nbsites $NUM_SITES < /tmp/in_C$i > /tmp/out_C$i &
  open -a "Google Chrome" http://localhost:63340/pixel-war/app-base/frontend/index.html & #valable pour macOS
  #ici : ajouter les alternatives de lancement selon les OS
done

# CONFIGURER LES CONNEXIONS ENTRE LES SITES
for i in $(seq 1 $NUM_SITES); do
  cat /tmp/out_A$i > /tmp/in_C$i &
  cat /tmp/out_C$i | tee /tmp/in_A$i $(for ((j=1;j<=$NUM_SITES;j++)); do if [ $j -ne $i ]; then echo "/tmp/in_C$j"; fi done) > /dev/null &
done