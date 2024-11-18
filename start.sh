#!/bin/bash
if [[ $DYNO == "web"* ]]; then
  ./bin/http 
elif [[ $DYNO == "worker"* ]]; then
  ./bin/consumer 
fi