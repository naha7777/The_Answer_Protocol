#!/bin/bash
./load.sh &
sleep 2
python3 latency.py
wait
