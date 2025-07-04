#!/bin/bash

podman build -t kubeapi . -f deploy/docker/Dockerfile
podman tag kubeapi bmzsombi/checklist
podman push bmzsombi/checklist