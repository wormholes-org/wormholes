# Wormholes Chain

The WormholesChain solves the blockchain trilemma, which entails a necessary tradeoff between scalability, security, and decentralization, by building the
technology to achieve the ideal balance between these three metrics, creating a highly scalable and secure blockchain system that doesn’t sacrifice
decentralization.

[![Gitter](https://badges.gitter.im/wormholes-org/Internal-test-miner.svg)](https://gitter.im/wormholes-org/Internal-test-miner?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## The Approach

The significant step before spinning up your node is choosing your approach. Based on requirements and many potential possibilities,
you must select the client implementation (of both execution and consensus clients), the environment (hardware, system), and the
parameters for client settings.

To decide whether to run the software on your hardware or in the cloud depending on your demands.

You can use the startup script to start your node after preparing the environment.

When the node is running and syncing, you are ready to use it, but make sure to keep an eye on its maintenance.

## Environment and Hardware

Wormholes clients are able to run on consumer-grade computers and do not require any special hardware, such as mining machines.
Therefore, you have more options for deploying the node based on your demands. let us think about running a node on both a local
physical machine and a cloud server:

### Hardware

Wormholes clients can run on your computer, laptop, server, or even a single-board computer. Although running clients on
different devices are possible, it had better use a dedicated machine to enhance its performance and underpin the security,
which can minimize the impact on your computer.

Hardware requirements differ by the client but generally are not that high since the node just needs to stay synced.

Do not confuse it with mining, which requires much more computing power. However , sync time and performance do improve with more
powerful hardware.

### Minimum requirements

-  CPU: Main frequency 2.9GHz, 4 cores or above CPU.
-  Memory: Capacity 8GB or more.
-  Hard Disk: Capacity 500GB or more.
-  Network bandwidth: 6M uplink and downlink peer-to-peer rate or higher

Before installing the client, please ensure your computer has enough resources to run it. You can find the minimum and recommended requirements below.

## Spin-up Your Own Wormholes Node

Running your node benefits you a lot, including opening more possibilities, and helping to support the ecosystem. This page will guide you through spinning
up your node and taking part in validating Wormholes transactions.

The actual client setup can be done by using the automatic launcher or manually.For ordinary users, we recommend you use a startup script, which guides you through the installation and automates the client setup process. However,
if you have experience with the terminal, the manual setup steps should be easy to follow.

### Docker Clients Setup

#### Install docker

You must ensure that docker has been installed. If it is not, please install docker first. For details, please refer to the [Rookie Tutorial](https://www.runoob.com/docker/ubuntu-docker-install.html). Or ignore it if already installed.

#### Prepare the JavaScript

Taking the Linux system as an example, the script Wormholes_install.sh to run the node can be named accordingly.

```
#!/bin/bash
#check docker cmd
which docker >/dev/null 2>&1
if  [ $? -ne 0 ] ; then
     echo "docker not found, please install first!"
     echo "ubuntu:sudo apt install docker.io -y"
     echo "centos:yum install  -y docker-ce "
     echo "fedora:sudo dnf  install -y docker-ce"
     exit
fi
#check docker service
docker ps > /dev/null 2>&1
if [ $? -ne 0 ] ; then

     echo "docker service is not running! you can use command start it:"
     echo "sudo service docker start"
     exit
fi

docker stop wormholes > /dev/null 2>&1
docker rm wormholes > /dev/null 2>&1
docker rmi wormholestech/wormholes:v1 > /dev/null 2>&1

if [ -d /wm/.wormholes/keystore ]; then
   read -p "Whether to clear the Wormholes blockchain data history, if yes, press the “y” button, and if not, click “enter.”：" xyz
   if [ "$xyz" = 'y' ]; then
         rm -rf /wm/.wormholes
              read -p "Enter your private key：" ky
   else
         echo "Do not clear"
   fi
else
   read -p "Please import your private key：" ky
fi

mkdir -p /wm/.wormholes/wormholes
if [ -n "$ky" ]; then
   echo $ky > /wm/.wormholes/wormholes/nodekey
fi

docker run -id -e KEY=$ky  -p 30303:30303 -p 8545:8545 -v /wm/.wormholes:/wm/.wormholes --name wormholes wormholestech/wormholes:v1

echo "Your private key is:"
sleep 6
docker exec -it wormholes /usr/bin/cat /wm/.wormholes/wormholes/nodekey
```

#### Run the node

When using the script to start the node, you must enter the private key of the account used for pledge prepared earlierWhen using the script to start the
node, you must enter the private key of the account used for pledge prepared earlier.

``` 
bash ./wormholes_install.sh 
```

### Manual clients setup

The actual client setup can be done by using the automatic launcher or manually.

For ordinary users, we recommend you use a startup script, which guides you through the installation and automates the client setup process. However, if
you have experience with the terminal, the manual setup steps should be easy to follow.

#### Startup parameters

- Start ***Wormholes*** in fast sync mode (default, can be changed withthe ***--syncmode*** flag),causing it to download more data in exchange for
avoiding processing the entire history of the Wormholes Chain network, which is very CPU intensive.

- Start up***Wormholes's*** built-in interactive JavaScript,(via the trailing ***console*** subcommand) through which you can interact using ***web3***
  [methods](https://learnblockchain.cn/docs/web3.js/getting-started.html)(note: the ***web3*** version bundled within ***Wormholes*** is very old, and
  not up to date with official docs), as well as ***Wormholes's*** own [management APIs](https://www.wormholes.com/docs/management/) .
  This tool is optional and if you leave it out you can always attach to an already running ``Wormholes`` instance with ***Wormholes attach*** .

#### Full nodes functions

-  Stores the full blockchain history on disk and can answer the data request from the network.
-  Receives and validates the new blocks and transactions.
-  Verifies the states of every account.

#### Start ordinary node

1. Download the binary, config and genesis files from [release](https://github.com/wormholes-org/wormholes), or compile the binary by ``make wormholes``.

2. Start your full node.

   ````
      # Ordinary nodes need to be started in full mode
      ./wormholes --devnet --syncmode=full

   ````

#### Start validator node

1. Download the binary, config and genesis files from [release](https://github.com/wormholes-org/wormholes), or compile the binary by ``make wormholes``.


2. Prepare a script to start the node, and name it run_node (If the system is windows, you need to add the file suffix .bat), or something else. Note that
the script should be in the same directory as the main Wormholes program.

   If you are using the windows system, the reference is as follows:

   ````
	@echo off
	set rootPath=%~dp0
	set nodePath=%rootPath%.wormholes
	if exist %nodePath% (
	     rd /s/q %nodePath%
	)

	if "%1" == "" (
	   echo "Please pass in the private key of the account to be pledged."
		exit -1
	) else (
	   md %nodePath%\geth
	   echo %1 > %nodePath%\geth\nodekey
	   wormholes.exe --devnet --datadir %nodePath% --mine --syncmode=full
	)
   
   ````

   If you are using a Linux system, the reference is as follows:

   ````
      
	#!/bin/bash
	# write private key to file
	if [ -d .wormholes ]; then
	   rm -rf .wormholes
	fi

	if [[ $# -gt 0 ]] ; then
	   mkdir -p .wormholes/wormholes
	   echo $1>.wormholes/wormholes/nodekey
	else
	   echo "Please pass in the private key of the account to be pledged."
	   return -1
	fi
	./wormholes --devnet --datadir .wormholes --syncmode=full
   
   ````

3. Start node.

   You can start a node based on the startup parameters, or you can start your own node using a startup script. If you use scripts, you need to select
   scripts for different environments according to different system environments. When running the startup script, you must pass in the private key of the
   account to be pledged, which is the private key saved in step 1.  the reference is as follows:

   ````

	#linux system
	#Runtime parameter -- private
	sh ./run_node 94b796b1b11893561c34cf000f23ecf3b39067bb198b9ec9f7b1a79646114680

	#windows system
	#Go to the directory where the startup script is located in the CMD terminal
	./run_node.bat 94b796b1b11893561c34cf000f23ecf3b39067bb198b9ec9f7b1a79646114680

   ````
