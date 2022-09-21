#!/bin/sh
sshpass -f osakapass.txt scp -P 5678 $1 15.152.100.177:/var/www/html/upload/
