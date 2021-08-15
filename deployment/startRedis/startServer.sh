#!/bin/bash

# Check root
if [[ $(id -u) != 0 ]] ; then
    echo "Must be run as root" >&2
    exit 1
fi

# Take the first argument to identify master or replica
arg1=$1

# Search for redis path in /root (assumes redis is installed under /root)
REDIS_PATH=$(find /root -maxdepth 1 -type d -name "redis-*.*.*")

# Read yaml
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
CONFIG_PATH="config.yaml"
eval $(parse_yaml $CONFIG_PATH "cfg_")

# Start server
if [ "$arg1" = "replica" ]
then
    $REDIS_PATH/src/redis-server --port $cfg_ReplicaPort --replicaof $cfg_MasterIp $cfg_MasterPort --appendonly no --save "" --daemonize yes
else 
    $REDIS_PATH/src/redis-server --port $cfg_MasterPort --appendonly no --save "" --daemonize yes
fi

$REDIS_PATH/src/redis-cli config set protected-mode no

# Success if you see PONG
$REDIS_PATH/src/redis-cli -h $cfg_MasterIp ping 
