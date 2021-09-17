#!/usr/bin/env bash

if [[ $1 == "i" ]]; then
  echo "Preparing infrastructure"
  ansible-playbook playbook.yml -i hosts.yml
elif [[ $1 == "f" ]]; then
  echo "Rebuilding frontend"
  cd ../../goj-frontend || exit
  npm run build || exit
  echo "Deploying frontend"
  cd ../goj/ansible || exit
  ansible-playbook deploy_frn_playbook.yml -i hosts.yml
elif [[ $1 == "b" ]]; then
  echo "Rebuilding backend"
  cd ..
  GOOS=linux GOARCH=amd64 go build -o deploy/goj || exit
  cd ansible || exit
  echo "Deploying backend"
  ansible-playbook deploy_playbook.yml -i hosts.yml
else
  echo "Use [f]rontend or [b]ackend option"
fi