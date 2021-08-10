# Redis Synchronous Replication

## TO RUN

1. Start VMs
    - The number of VMs depends on your setup. The minimum is 2 VMs (one being the master and the other being a replica).

2. On each VM, complete the following steps:

    - Clone the repository
    ```shell

    sudo su && cd ~
    mkdir -p ~/go/src
    cd ~/go/src
    git clone https://github.com/YichengShen/redis-sync-rep.git
    cd redis-sync-rep

    ```

    - Install Redis, Go, and dependancies
    ```shell
    . ./deployment/install/install.sh
    ```

    - Start Redis server
        - configure the current VM as a master
            ```shell
            . ./deployment/startRedis/startServer.sh
            ```
        - configure the current VM as a replica
            ```shell
            . ./deployment/startRedis/startServer.sh replica
            ```
        Note: For the minimum requirement, you need 2 VMs. You configure one to be the master and the other to be a replica using the commands above.
        
    - Adjust parameters in `config.yaml`
   
3. From one of the client VMs, run the main program
    ```shell
    go run main.go
    ```
    The program will run according to the configuration file and print out the results which will also be saved in the logs folder.  
