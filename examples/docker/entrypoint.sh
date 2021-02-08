#!/bin/sh
/probr/probr -varsfile=/probr/config.yml
node /probr/internal/view/index.js /probr/probr_output/cucumber
