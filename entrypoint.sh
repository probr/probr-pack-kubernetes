#!/bin/sh

./probr
cd internal/view
mkdir ./testoutput
node index.js ./testoutput
