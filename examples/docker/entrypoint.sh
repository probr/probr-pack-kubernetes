#!/bin/sh
/probr/probr -varsfile=/probr/config.yml
node /probr/internal/view/index.js /probr/cucumber_output
