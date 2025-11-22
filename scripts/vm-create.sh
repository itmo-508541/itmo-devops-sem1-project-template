#!/bin/bash

yc compute instance create \
  --name itmo508541_vm \
  --description "sem1-project-vm" \
  --zone ru-central1-d \
  --cores 2 \
  --memory 2 \
  --preemptible \
  --create-boot-disk name=itmo508541-disk,type=network-ssd,size=10,image-folder-id=standard-images,image-family=ubuntu-2404-lts \
  --network-interface subnet-name=itmo508541-subnet,nat-ip-version=ipv4 \
  --metadata-from-file user-data=./scripts/user-data.yaml \
  --format json  > $1 || yc compute instance start --name itmo508541_vm --format=json > $1 || yc compute instance get --name itmo508541_vm --format=json > $1
