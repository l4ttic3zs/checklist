#!/bin/bash

podman build -t kubeapi .
podman tag kubeapi bmzsombi/checklist
podman push bmzsombi/checklist