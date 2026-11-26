#!/bin/bash
for i in $(seq 1 20); do
  (printf 'CONNECT load%s\n' $i; sleep 1
   for j in $(seq 1 200); do echo who; sleep 0.07; done
   sleep 1) | timeout 30 nc localhost 8080 > /dev/null &
done
wait
