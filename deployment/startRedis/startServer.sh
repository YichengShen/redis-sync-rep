#!/bin/bash

# IP of master
master_ip="10\.142\.0\.58"

arg1=$1

# Check root
if [[ $(id -u) != 0 ]] ; then
    echo "Must be run as root" >&2
    exit 1
fi

# If master IP not added, then add it into redis.conf
if ! grep "$master_ip" /etc/redis/redis.conf
then
    # Add master_ip after "bind 127.0.0.1 ::1"
    sed -ie "s/^bind 127.0.0.1 ::1/& $master_ip/g" /etc/redis/redis.conf
    echo "Added new Master IP"
else
    echo "Master IP already added"
fi

# If ran as replica, change redis.conf accordingly
if [ "$arg1" = "replica" ]
then
    if grep "# replicaof" /etc/redis/redis.conf
    then
        sed -i "s/.*# replicaof.*/replicaof $master_ip 6379/" /etc/redis/redis.conf
        echo "Uncommented replicaof and wrote master_ip 6379"
    else
        echo "replicaof already uncommented"
    fi
fi

# restart Redis server
systemctl restart redis-server.service

# give redis-server a second to wake up
sleep 1

# open up access to the Redis port
ufw allow 6379

redis-cli -h 10.142.0.58 ping
