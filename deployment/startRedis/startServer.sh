#!/bin/bash

parse_yaml() {
    local prefix=$2
    local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
    sed -ne "s|^\($s\)\($w\)$s:$s\"\(.*\)\"$s\$|\1$fs\2$fs\3|p" \
            -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
    awk -F$fs '{
        indent = length($1)/2;
        vname[indent] = $2;
        for (i in vname) {if (i > indent) {delete vname[i]}}
        if (length($3) > 0) {
            vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
            printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
        }
    }'
}

# Read IP of master from configuration file
CONFIG_PATH="../../config.yaml"
eval $(parse_yaml $CONFIG_PATH "cfg_")
 
arg1=$1

# Check root
if [[ $(id -u) != 0 ]] ; then
    echo "Must be run as root" >&2
    exit 1
fi

# If master IP not added, then add it into redis.conf
if ! grep "$cfg_MasterIp" /etc/redis/redis.conf
then
    # Add cfg_MasterIp after "bind 127.0.0.1 ::1"
    sed -ie "s/^bind 127.0.0.1 ::1/& $cfg_MasterIp/g" /etc/redis/redis.conf
    echo "Added new Master IP"
else
    echo "Master IP already added"
fi

# If ran as replica, change redis.conf accordingly
if [ "$arg1" = "replica" ]
then
    if grep "# replicaof" /etc/redis/redis.conf
    then
        sed -i "s/.*# replicaof.*/replicaof $cfg_MasterIp 6379/" /etc/redis/redis.conf
        echo "Uncommented replicaof and wrote $cfg_MasterIp 6379"
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

redis-cli -h $cfg_MasterIp ping
