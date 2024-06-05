#!/bin/bash

NUM_SITES=3

# CRÉATION DES FIFOS/DES LIENS
for i in $(seq 1 $NUM_SITES); do
  mkfifo /tmp/in_A$i /tmp/out_A$i
  mkfifo /tmp/in_C$i /tmp/out_C$i
  mkfifo /tmp/in_N$i /tmp/out_N$i
done

# LANCEMENT DES SITES
for i in $(seq 1 $NUM_SITES); do
  go run app-base -n A$i -m g -port $((4443 + i)) < /tmp/in_A$i > /tmp/out_A$i &
  go run app-control -n C$i -nbsites $NUM_SITES < /tmp/in_C$i > /tmp/out_C$i &
  go run app-net -n N$i < /tmp/in_N$i > /tmp/out_N$i &
  open -a "Google Chrome" http://localhost:63340/pixel-war/app-base/frontend/index.html & #valable pour macOS
  #ici : ajouter les alternatives de lancement selon les OS
done

# CONFIGURER LES CONNEXIONS ENTRE LES SITES
for i in $(seq 1 $NUM_SITES); do
  next=$((i % NUM_SITES + 1))
  cat /tmp/out_A$i > /tmp/in_C$i &
  cat /tmp/out_C$i | tee /tmp/in_A$i > /tmp/in_N$i &
  cat /tmp/out_N$i | tee /tmp/in_C$i > /tmp/in_N$next &

done
