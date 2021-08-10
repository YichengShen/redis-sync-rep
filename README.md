# Redis Synchronous Replication

## TO RUN

1. Clone the repository
```shell

sudo su && cd ~
mkdir -p ~/go/src
cd ~/go/src
git clone https://github.com/YichengShen/redis-sync-rep.git
cd redis-sync-rep

```

2. Install
```shell
. ./deployment/install/install.sh
```

3. Start Redis server
    - run as master
        ```shell
        . ./deployment/startRedis/startServer.sh
        ```
    - run as replica
        ```shell
        . ./deployment/startRedis/startServer.sh replica
        ```
   
4. Run Clients
```shell
go run main.go
```
