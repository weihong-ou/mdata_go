# Purpose
This repository is meant to stimulate ideas on how to store public data records on Hyperledger Sawtooth. 
It's a slightly modified copy of xo_go transaction processor, part of github.com/hyperledger/sawtooth-sdk-go

# RFC
Please read and review the RFC located [here](docs/RFC.md).

# System Requirements
1. OS Packages
    ```
    sudo apt install zip curl python3 python3-pip
    ```

2. Install protobuf compilers (make sure to get a 3.x.x version from [here](https://github.com/protocolbuffers/protobuf/releases))
    ```
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
    unzip protoc-3.7.1-linux-x86_64.zip -d protoc3
    sudo mv protoc3/bin/* /usr/local/bin/
    sudo mv protoc3/include/* /usr/local/include/
    ```

2. Install python's grpcio-tools library
    ```
    sudo su - 
    python3 -m pip install grpcio-tools
    ```

3. Install go dependencies
    ```
    go get -u google.golang.org/grpc && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get github.com/satori/go.uuid && \
    go get github.com/pebbe/zmq4 && \
    go get github.com/golang/mock/gomock && \
    go install github.com/golang/mock/mockgen && \
    go get github.com/hyperledger/sawtooth-sdk-go && \
    cd $GOPATH/src/github.com/hyperledger/sawtooth-sdk-go && \
    go generate && \
    cd && \
    go get github.com/jessevdk/go-flags && \
    go get github.com/stretchr/testify/mock && \
    go get github.com/btcsuite/btcd/btcec && \
    go get gopkg.in/yaml.v2
    ```

#Restting the Sawtooth Test Network
After shutting down all instances of a test network, I find that the network can no longer reach consensus when rebooted. Since I do not need the network up all the time, just when I test, I find it simpler to rebuild the network when I restart all the nodes.

## Delete Existing Sawtooth Data
```
sudo su -
rm -r /var/lib/sawtooth/*
```

## Generate new genesis block
```
sudo sawset genesis -k /etc/sawtooth/keys/validator.priv -o config-genesis.batch \
cd /tmp \
sudo -u sawtooth sawset proposal create -k /etc/sawtooth/keys/validator.priv \
sawtooth.consensus.algorithm.name=pbft \
sawtooth.consensus.algorithm.version=0.1 \
sawtooth.consensus.pbft.peers=['"'$(paste ~/fleet_keys/*.pub -d , | sed s/,/\",\"/g)'"'] \
sawtooth.consensus.pbft.view_change_timeout=4000 \
sawtooth.consensus.pbft.message_timeout=10 \
sawtooth.consensus.pbft.max_log_size=1000 \
-o config.batch \ 
sudo mv config.batch ~/ \
cd \
sudo -u sawtooth sawadm genesis config-genesis.batch config.batch
```