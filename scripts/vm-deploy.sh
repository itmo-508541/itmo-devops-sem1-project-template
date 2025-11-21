#!/bin/bash

docker compose -f docker-compose.deploy.yml config > ./build/deploy.yml

while ! nc -z -v -w30 $1 22; do
  sleep 1
done

ssh-keyscan runner@$1 >> ~/.ssh/known_hosts
scp -P 22 -r ./build/deploy.yml runner@$1:/tmp/deploy.yml

ssh runner@$1 "if ! command -v docker; then sudo apt update && sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin ; fi ;
  cd /srv ;
  docker compose stop || true ;
  sudo mv /tmp/deploy.yml docker-compose.yml ;
  docker compose pull -q ;
  docker compose up -d"
