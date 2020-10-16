#!/bin/sh

./probr -outputType=IO -outputDir=./cucumber_output
node internal/view/index.js ./cucumber_output
