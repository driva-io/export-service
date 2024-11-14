#!/bin/bash
if [[ $DYNO == "web"* ]]; then
  ./http 
elif [[ $DYNO == "worker"* ]]; then
  ./consumer 
fi