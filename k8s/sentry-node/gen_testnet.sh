#!/bin/bash
set -eu
VERSION=$(cat ../../version)
VALNUM=$1

vals=()
nvals=()
for i in $(seq 0 $(expr $VALNUM - 1)); do
    SENTRYNO=$(expr $i + $VALNUM)
    sed -e "s/<VERSION>/${VERSION}/g" -e "s/<NODENO>/$i/g" -e "s/<SENTRYNO>/${SENTRYNO}/g" deployment.yaml.tpl > ./config/deployment${i}.yaml
    echo "generate deployment$i.yaml"
    vals+=(${i})
    nvals+=($SENTRYNO)
done
vals=$(IFS=','; echo "${vals[*]}")
nvals=$(IFS=','; echo "${nvals[*]}")
sed -e "s/<VALIDATORS>/${vals}/g" -e "s/<SENTRYNODES>/${nvals}/g" config.yaml.tpl > ./config/config.yaml
echo "generate config.yaml"
