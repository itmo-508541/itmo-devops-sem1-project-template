#!/bin/bash

while ! nc -z -v -w30 $1 22; do
  sleep 1
done

ssh-keyscan $1 >> ~/.ssh/known_hosts
ssh $1 "curl -fsSL https://get.docker.com -o get-docker.sh ;
  sudo bash ./get-docker.sh ;
  sudo chown runner:runner /srv ;
  rm /srv/deploy.yml || true"

mkdir --parent ./build
docker compose -f docker-compose.deploy.yml config > ./build/deploy.yml
scp -P 22 -r ./build/deploy.yml $1:/srv/deploy.yml

ssh $1 -o StrictHostKeyChecking=no "cd /srv ;
  docker compose stop || true ;
  rm docker-compose.yml || true ;
  mv deploy.yml docker-compose.yml ;
  docker compose pull ;
  docker compose up -d"
