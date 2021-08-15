# Sync-Rep (Synchronous replication using standalone Redis)

## 2 Configurations
- (i) Sync-Rep (1): two VMs, one as the Redis leader and the other as a follower; and 
- (ii) Sync-Rep (2): three VMs, one as the Redis leader and the other two as followers.
- Note: We run clients on follower VMs to force the system to have one RTT so that it’s compatible with SMR-based approach.

## To Run

1. Start VMs
    - Start appropriate number of VMs according to Sync-Rep (1) or (2).
    - OS Assumption: Ubuntu 16.04

2. On each VM, complete the following steps:

    - Clone [Rabia](https://github.com/haochenpan/rabia) repository
        ```shell
        sudo su
        mkdir -p ~/go/src && cd ~/go/src
        git clone https://github.com/haochenpan/rabia.git
        ```

    - Install Rabia and its dependencies (Dependencies of Sync-Rep are included in Rabia's installation script.)
        ```shell
        cd ./rabia/deployment
        . ./install/install.sh
        ```

    - Clone Sync-Rep repository
        ```shell
        cd ~/go/src
        git clone https://github.com/YichengShen/redis-sync-rep.git
        cd redis-sync-rep
        ```

    - Configure IP of master VM
        - In `config.yaml`, change 'MasterIp' to the IP of your master VM.

    - Start Redis server: You could configure the current VM either as a master or a replica.
        - configure as master
            ```shell
            . ./deployment/startRedis/startServer.sh
            ```
        - configure as replica
            ```shell
            . ./deployment/startRedis/startServer.sh replica
            ```
        
    - Adjust parameters related to batching in `config.yaml`
        - Change ‘NClients’ and ‘ClientBatchSize’.
   
3. From one of the client VMs, run the main program
    ```shell
    go run main.go
    ```
    The program will run according to the configuration file and print out the results which will also be saved in the logs folder.  
