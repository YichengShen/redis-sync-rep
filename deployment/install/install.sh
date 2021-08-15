# OS assumption: Ubuntu 16.04

sync_rep_folder=~/go/src/redis-sync-rep          # the path to the folder cloned from Github
redis_folder=~                                   # the path where Redis is installed
go_tar=go1.15.8.linux-amd64.tar.gz               # the version of Golang to be downloaded in install_go

redis_ver=redis-6.2.2   # version of Redis

# Copies the SSH public key to the "authorized_keys" file
function install_key() {
    mkdir -p ~/.ssh/
    cat "${sync_rep_folder}"/deployment/install/id_ed25519.pub >>~/.ssh/authorized_keys
    chmod 400 "${sync_rep_folder}"/deployment/install/id_ed25519
}

# Install Redis from source (currently in use)
function install_redis_from_source() {
    cd ${redis_folder}
    sudo apt install -y tar make
    wget https://download.redis.io/releases/${redis_ver}.tar.gz
    tar xzf ${redis_ver}.tar.gz
    rm ${redis_ver}.tar.gz
    cd ${redis_ver}
    make
    cd $sync_rep_folder
}

# Install Redis using apt-get
function install_redis_apt() {
    echo | sudo apt-add-repository ppa:chris-lea/redis-server
    sudo apt-get update
    sudo apt-get install --assume-yes redis-server
}

# Installs a version of Golang
function install_go() {
    wget -q https://golang.org/dl/${go_tar}
    sudo tar -C /usr/local -xzf ${go_tar}
    rm ${go_tar}
    echo 'export PATH=${PATH}:/usr/local/go/bin' >>~/.bashrc
    echo 'export GOPATH=~/go' >>~/.bashrc
    source ~/.bashrc
    go version
}

# Installs Golang packages
function install_go_deps() {
    go get github.com/go-redis/redis/v8   # go-redis
    go get -u github.com/rs/zerolog/log
}

install_key
install_redis_from_source
install_go