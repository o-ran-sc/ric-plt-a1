# Continuous Integration

This directory contains configuration for the Jenkins CI system to
build a docker image with the software in this repository.  Each
docker build requires at least the following files:

* Dockerfile - instructions for assembling the image
* container-tag.yaml - string to use as the tag, usually a version number

Additional files required to build a Docker image, for example a
script for the CMD or ENTRYPOINT, might also be stored here.
