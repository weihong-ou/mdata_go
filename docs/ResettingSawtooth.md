# Resetting the Sawtooth Test Network
After shutting down all instances of a test network, I find that the network can no longer reach consensus when rebooted. Since I do not need the network up all the time, just when I test, I find it simpler to rebuild the network when I restart all the nodes.

## Delete Existing Sawtooth Data
```
sudo su -
rm -r /var/lib/sawtooth/*
exit
```

## Generate new genesis block
```
sudo sawset genesis -k /etc/sawtooth/keys/validator.priv -o config-genesis.batch &&\
cd /tmp &&

sudo -u sawtooth sawset proposal create -k /etc/sawtooth/keys/validator.priv \
sawtooth.consensus.algorithm.name=pbft \
sawtooth.consensus.algorithm.version=0.1 \
sawtooth.consensus.pbft.peers=['"'$(paste ~/fleet_keys/*.pub -d , | sed s/,/\",\"/g)'"'] \
sawtooth.consensus.pbft.view_change_timeout=4000 \
sawtooth.consensus.pbft.message_timeout=10 \
sawtooth.consensus.pbft.max_log_size=1000 \
-o config.batch &&

sudo mv config.batch ~/ &&\
cd &&\
sudo -u sawtooth sawadm genesis config-genesis.batch config.batch
```