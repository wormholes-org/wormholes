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
#docker rmi wormholes:v1 > /dev/null 2>&1
#while :
#do
    #if [ -z "$ky" ];then
        #echo "密钥不能为空"
        #continue
	#fi
	#break
#done
if [ -d /wm/.wormholes/keystore ]; then
    read -p "是否要清空wormholes区块链历史数据,清空请按y,不清空请直接回车：" xyz
    if [ "$xyz" = 'y' ];then
        rm -rf /wm/.wormholes
	read -p "请输入您的私钥：" ky
    else
        echo "不清空"
    fi
fi
#mkdir -p /wm
mkdir -p /wm/.wormholes/wormholes
if [ -n "$ky" ]; then
    echo $ky > /wm/.wormholes/wormholes/nodekey
fi
#docker run -id -e KEY=$ky -p 30303:30303 -p 8545:8545 -v /var/lib/wormholes:/root/.wormholes --name wormholes wormholestech/wormholes:v1
docker run -id -e KEY=$ky  -p 30303:30303 -p 8545:8545 -v /wm/.wormholes:/wm/.wormholes --name wormholes wormholestech/wormholes:v1


echo "您的私钥是:" 
sleep 6
docker exec -it wormholes /usr/bin/cat /wm/.wormholes/wormholes/nodekey
#cat /wm/.wormholes/geth/nodekey
