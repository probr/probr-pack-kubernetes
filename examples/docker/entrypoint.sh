#!/bin/sh

./probr -outputType=IO -outputDir=./testoutput
node internal/view/index.js ./testoutput
