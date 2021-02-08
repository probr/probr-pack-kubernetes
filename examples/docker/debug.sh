#!/bin/sh
/probr/probr -varsfile=/probr/config.yml -loglevel=DEBUG
node /probr/internal/view/index.js /probr/cucumber_output
