#!/bin/bash

# IP of master
master_ip="10\.142\.0\.58"

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

# restart Redis server
systemctl restart redis-server.service

# give redis-server a second to wake up
sleep 1

# open up access to the Redis port
ufw allow 6379

redis-cli -h 10.142.0.58 ping
