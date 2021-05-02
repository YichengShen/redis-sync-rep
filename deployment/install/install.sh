# OS assumption: Ubuntu 16.04

redis_folder=~/go/src/redis-testing           # the path to the folder cloned from Github
go_tar=go1.15.8.linux-amd64.tar.gz            # the version of Golang to be downloaded in install_go

# Copies the SSH public key to the "authorized_keys" file
function install_key() {
    mkdir -p ~/.ssh/
    cat "${redis_folder}"/deployment/install/id_ed25519.pub >>~/.ssh/authorized_keys
    chmod 400 "${redis_folder}"/deployment/install/id_ed25519
}

# Install Redis
function install_redis() {
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
}

install_key
install_redis
install_go