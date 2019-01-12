#!/bin/bash
set -eu
VALNUM=$(expr $1 - 1)

peers=()
for i in $(seq 0 $VALNUM); do
    sed "s/<NODENO>/$i/g" deployment.yaml.tpl > ./config/deployment${i}.yaml
    echo "generate deployment$i.yaml"
    peers+=(${i})
done
peers=$(IFS=','; echo "${peers[*]}")
sed "s/<VALIDATORS>/0,1/g" config.yaml.tpl > ./config/config.yaml
echo "generate config.yaml"
