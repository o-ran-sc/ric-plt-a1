# Continuous Integration

This directory contains the configuration for the Jenkins CI system to
build a docker image with the software in this repository.  Each
docker build requires at least the following files:

* Dockerfile - instructions for assembling the image
* container-tag.yaml - string to use as the tag, usually a version number

Additional files required to build a Docker image, for example a
script for the CMD or ENTRYPOINT, can also be stored here.

This directory is used as the DOCKER_ROOT by the CI job.  All files
should be in this directory if only a single image must be built.  If
multiple images must be built, create additional subdirectories (e.g.,
"testimage") for each one so that they can be used as a DOCKER_ROOT.
