# distributed-election
Distributed election algorithms. Project for the exam "Sistemi Distribuiti e Cloud Computing". 

## Overview
This application, made in the context of the "*Sistemi distribuiti e Cloud Computing*" (9 CFU = 9 ETCS) course, followed at *University of Rome Tor Vergata*, aims to implement two *distributed election* algorithms, i.e.:
- H. Garcia-Molina's *bully* algorithm (1982)
- G. N. Fredrickson and N. A. Lynch algorithm, based on a ring topology (1987)

The project's functionalities are managed by some **nodes**, that elect a coordinator between them, and a **service registry**, that manages the node's network addresses.

Further information about this project are available into the [`docs`](https://github.com/massimostanzione/distributed-election/tree/main/docs) folder.

## Quickstart
You can either launch a [local run](#local-run) or see the [deployment](#deployment) in action.

### Local Run
Some *ready-to-use* scripts are available into the [`examples`](https://github.com/massimostanzione/distributed-election/tree/main/examples) folder. Follow these instructions to run them:

0. Install Go by following [these instructions](https://go.dev/doc/install), then install gRPC and protocol-buffer if you don't have it in your machine:
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
unzip protoc-3.15.8-linux-x86_64.zip -d $HOME/.local
export PATH="$PATH:$HOME/.local/bin"
```

1. Download the repository and apply privileges:
```
git clone https://github.com/massimostanzione/distributed-election
sudo chmod -R 770 distributed-election
cd distributed-election
```
2. Build the binaries for *nodes* and *service registry*:
```
cd scripts
./setup.sh
```
*(Please run `./teardown.sh` if something goes wrong.)*

3. Just launch the algorithm you want:
```
cd ../examples

# to run a bully algorithm demo:
./bully.sh

# to run a Fredrickson-Lynch algorithm demo:
./fredricksonlynch.sh
```
... then just follow the instructions and enjoy!

### Deployment
**Notice:** to execute this you have to have an active *Amazon EC2* instance.
```
cd ../deployments/ansible

# use vim, nano, or whatever you prefer
vim hosts.ini
```

Substitute the following items:
- Replace `HOST-IP-ADDRESS` with your *Amazon EC2* instance outbound IP address;
- Replace `AMAZON_EC2_USERNAME` with your *Amazon EC2* instance username (default: `ec2-user`);
- Replace `PATH/TO/KEYS.pem` with your `.pem` file containing the private key you use to access your *Amazon EC2* instance.

Now just deploy:
```
# install ansible, if you don't have it in your machine:
python3 -m pip install --user ansible

# start the deployment
ansible-playbook -v -i hosts.ini deploy.yaml
```
*(if it is not your first Ansible run, the flag `--flush-cache` may be useful)*

Now, to see what is running, connect to your *Amazon EC2* instance:

```
ssh -i "PATH/TO/KEYS.pem" AMAZON_EC2_USERNAME@HOST-IP-ADDRESS

# Verify if containers are up and running
# ***NOTICE:*** Containers can take several (about 3, empirically) minutes before
# being made up and running! Be patient...
docker ps

# Inspect logs for all the containers
cd distributed-election/deployments
docker container logs

# or, to see log just for one container, take its name from 'docker ps' and run, from any position:
docker container logs CONTAINER_NAME
```

## Keep it calm
If you already enjoyed the [quickstart](#quickstart), you can now launch a more *reasoned* execution:

Assuming you already followed the [setup](#setup) section, move to the `distributed-election` folder and run:
```
cd bin
```
here you can find the executable `node` and `serviceregistry`. You can run them in separate terminal windows or place them into different nodes of a network and just run with
```
./node
```
or
```
./serviceregistry
```

Please take a look at the `--help` flag to see all the possible customization (further more of them, with the whole set of parameters that can be set, can be found into the [`configs`](https://github.com/massimostanzione/distributed-election/tree/main/configs) folder.
