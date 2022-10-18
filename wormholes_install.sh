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
if [ -d /wm/.wormholes/wormholes ]; then
    read -p "If you want to clear the historical data of the wormholes blockchain, please press y to clear it, if not, please press Enter directly:" xyz
    if [ "$xyz" = 'y' ];then
        cp /wm/.wormholes/wormholes/nodekey /wm/nodekey
        rm -rf /wm/.wormholes
        mkdir -p /wm/.wormholes/wormholes
        mv /wm/nodekey /wm/.wormholes/wormholes/nodekey
    else
        echo "not empty"
    fi
else
    read -p "Please enter your private key:" ky
fi

if [ -n "$ky" ]; then
    if [ ${#ky} -eq 64 ];then
        mkdir -p /wm/.wormholes/wormholes
        echo $ky > /wm/.wormholes/wormholes/nodekey
    elif [ ${#ky} -eq 66 ] && ([ ${ky:0:2} == "0x" ] || [ ${ky:0:2} == "0X" ]);then
        mkdir -p /wm/.wormholes/wormholes
        echo ${ky:2:64} > /wm/.wormholes/wormholes/nodekey
    else
        echo "the nodekey format is not correct"
        exit -1
    fi
fi
docker run -id -e KEY=$ky  -p 30303:30303 -p 8545:8545 -v /wm/.wormholes:/wm/.wormholes --name wormholes wormholestech/wormholes:v1


echo "Your private key is:"
sleep 6
docker exec -it wormholes /usr/bin/cat /wm/.wormholes/wormholes/nodekey
#cat /wm/.wormholes/geth/nodekey
