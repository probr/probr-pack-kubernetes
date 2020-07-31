# WIP - DO NOT MERGE

I'm having trouble getting the container to `CMD ['./run.sh']` even though I've validated that the file is present... it just keeps saying the file isn't found. 

Based on my reading, this may be because I'm building the container from my windows CLI and my line endings are funking up the build. Since that's my best guess, I did attempt to convert the line endings, but that effort didn't produce any results thus far.

To rule out the above theory, could someone attempt to build this Dockerfile and find out whether it tips over when the container is run? Will keep poking this, but I'm out of working theories for the moment.
