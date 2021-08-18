# VK Jeopardy backend service

## Prepare remote server

* Create Ubuntu Linux VM
* Create inventory file. See infra/hosts_example.yml
* Run:

```shell
cd ansible
ansible-playbook playbook.yml -i <inventory file>
```

## Build And Deploy

Run:

```shell
GOOS=linux GOARCH=amd64 go build -o deploy/goj
cd ansible
ansible-playbook deploy_playbook.yml -i <inventory file> 
```

## Deploy frontend

```shell
cd ansible
./deploy.sh f  
```

## Deploy backend

```shell
cd ansible
./deploy.sh b  
```
