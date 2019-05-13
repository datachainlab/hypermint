apiVersion: v1
kind: Service
metadata:
  name: hm-validator-<NODENO>
spec:
  type: NodePort
  selector:
    app: hm-validator-<NODENO>
  ports:
  - port: 26656
    name: p2p
  - port: 26657
    name: rpc

---

apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: hm-validator-<NODENO>
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hm-validator-<NODENO>
  serviceName: hm-validator-<NODENO>
  template:
    metadata:
      labels:
        app: hm-validator-<NODENO>
    spec:
      containers:
        - name: hm-validator
          image: bluele/hypermint-testnet-sentry:<VERSION>
          command:
            - bash
            - "-c"
            - |
              set -ex
              sentryid=$(/hmd tendermint show-node-id --home /mytestnet/node<SENTRYNO>/hmd)
              peer="${sentryid}@127.0.0.1:26656"
              state_file=/mytestnet/node<NODENO>/hmd/data/priv_validator_state.json
              if [ ! -e $state_file ]; then
              cat << EOF > $state_file
              {
                "height": "0",
                "round": "0",
                "step": 0
              }
              EOF
              fi
              TM_PARAMS="p2p.allow_duplicate_ip=true,p2p.addr_book_strict=false" /hmd start --home=/mytestnet/node<NODENO>/hmd \
              --log_level="*:info" \
              --p2p.persistent_peers=$peer \
              --p2p.laddr="tcp://0.0.0.0:36656" \
              --rpc.laddr="tcp://0.0.0.0:36657" \
              --address="tcp://0.0.0.0:36658"
          volumeMounts:
            - mountPath: /mytestnet/node<NODENO>/hmd/data
              name: hmdir-validator-<NODENO>

        - name: hm-sentry
          image: bluele/hypermint-testnet-sentry:<VERSION>
          env:
          - name: VALIDATORS
            valueFrom:
              configMapKeyRef:
                name: hm-config
                key: validators
          command:
            - bash
            - "-c"
            - |
              set -ex
              IFS=',' read -ra VALIDATORS_NUM <<< "$VALIDATORS"
              peers=()

              for n in "${VALIDATORS_NUM[@]}"; do
                nodeid=$(/hmd tendermint show-node-id --home /mytestnet/node$(expr ${n} + 4)/hmd)
                peers+=("${nodeid}@hm-validator-${n}:26656")
              done

              validatorid=$(/hmd tendermint show-node-id --home /mytestnet/node<NODENO>/hmd)
              peers+=("${validatorid}@127.0.0.1:36656")
              peers=$(IFS=','; echo "${peers[*]}")

              state_file=/mytestnet/node<SENTRYNO>/hmd/data/priv_validator_state.json
              if [ ! -e $state_file ]; then
              cat << EOF > $state_file
              {
                "height": "0",
                "round": "0",
                "step": 0
              }
              EOF
              fi
              TM_PARAMS="p2p.allow_duplicate_ip=true,p2p.addr_book_strict=false,p2p.private_peer_ids=${validatorid}" /hmd start --home=/mytestnet/node<SENTRYNO>/hmd \
              --log_level="*:info" \
              --p2p.persistent_peers="$peers" \
              --p2p.laddr="tcp://0.0.0.0:26656" \
              --rpc.laddr="tcp://0.0.0.0:26657" \
              --address="tcp://0.0.0.0:26658"
          ports:
            - containerPort: 26656
              name: p2p
            - containerPort: 26657
              name: rpc
          volumeMounts:
            - mountPath: /mytestnet/node<SENTRYNO>/hmd/data
              name: hmdir-sentry-<SENTRYNO>

  volumeClaimTemplates:
    - metadata:
        name: hmdir-validator-<NODENO>
      spec:
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi
    - metadata:
        name: hmdir-sentry-<SENTRYNO>
      spec:
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi
