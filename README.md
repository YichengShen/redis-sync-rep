# Redis Testing

## TO RUN

1. Clone the repository
```shell

sudo su && cd ~
mkdir -p ~/go/src
cd ~/go/src
git clone https://github.com/YichengShen/redis-testing.git
cd redis-testing

```

2. Install
```shell
. ./deployment/install/install.sh
```

3. Start Redis server
    - run as master
        ```shell
        . ./deployment/run/run.sh
        ```
    - run as replica
        ```shell
        . ./deployment/run/run.sh replica
        ```
    