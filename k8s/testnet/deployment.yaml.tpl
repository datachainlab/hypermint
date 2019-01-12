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
        - name: hm
          image: bluele/hypermint-testnet:0.1.0
          env:
          - name: NODENO
            value: "<NODENO>"
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
              for a in "${VALIDATORS_NUM[@]}"; do
                NODEID=$(/hmd tendermint show-node-id --home /mytestnet/node${a}/hmd)
                peers+=("${NODEID}@hm-validator-${a}:26656")
              done
              peers=$(IFS=','; echo "${peers[*]}")

              /hmd start --home=/mytestnet/node${NODENO}/hmd \
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
            - mountPath: /mytestnet/node<NODENO>/hmd/data
              name: hmdir<NODENO>
  volumeClaimTemplates:
    - metadata:
        name: hmdir<NODENO>
      spec:
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi
