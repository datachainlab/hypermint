# Testnet

## Getting started

```
# Generate k8s yaml files, and build docker image for testnet.
$ make VALS_NUM=4 gen-config

# Create a testnet with generated config.
$ make create

# Destory testnet
$ make destroy
```
