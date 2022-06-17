# fabric-ccs-builder

**Archived:** This project has now been archived. A [chaincode as a service (CCaaS) builder](https://github.com/hyperledger/fabric/tree/main/ccaas_builder) is now available in the main fabric repository.

## Context

The existing model for chaincode within Fabric, is based around the peer orchestrating the creation and starting of a docker image. Given a set of code, the peer would create a suitable docker container with this code and then start the container. The peer would supply and TLS Certificates, information about the identitiy of the chaincode and the address of the peer. The chaincode library then connects, via a long running gRPC connection, to the peer to 'register'. Transaction requests are then sent from the peer to the chaincode. 

'out-of-the-box' Fabric can understand and create docker images for Java, Node and Go chaincode. The 'Chaincode Builder' interface within Fabric allows you to define a 'builder' that would allow the user to configure a set of scripts or binaries that know how to build and run your choice of chaincode. With this model, the chaincode is in a client role, and needs to connect to the peer to register. The builder code also needs to be able to control where the chaincode runs.

With 'as-service' approach, the roles are swapped. The chaincode that is installed on the peer is a definition of where the chaincode is (host/port/tls certs etc), and any CouchDB indexes needed. If this is say a k8s enviroment, the chaincode needs to be started 'as-a-service' (sometimes called as-a-server).  When the peer needs to send a transaction to the chaincode, the peer connects (in a client role) to the chaincode (in a server role). From this point onwards the communication between the peer/chaincode is indentical to all other cases.

### References

- [Chaincode Builder Model](https://hyperledger-fabric.readthedocs.io/en/release-2.2/cc_launcher.html)

- [Chaincode as a Service](https://hyperledger-fabric.readthedocs.io/en/release-2.2/cc_service.html)


## This repo

The documentation (above) contains references and examples of how you can use bash scripts to setup this type of environment. However the usual docker images for Hyperledger Fabric don't contain BASH.

Golang code is more readily and reliably run the peer, and this repo contains the go binaries to do exactly this. 


## Docker Image

Build a docker image embedding the `build`, `release`, and `detect` binaries into the /go/bin directory: 

```shell 
docker build \
  -t hyperledgendary/fabric-ccs-builder \
  .
```

## Peer Configuration

For Fabric networks running on Kubernetes, the ccs-builder image may be used by a sidecar or init container to load the external builder routines into pods running the `hyperledger/fabric-peer`.  For example, the registration of an external builder in core.yaml: 

```yaml
    externalBuilders:
      - path: /var/hyperledger/fabric/chaincode/ccs-builder
        name: ccs-builder
        propagateEnvironment:
          - HOME
          - CORE_PEER_ID
          - CORE_PEER_LOCALMSPID
```

may be fulfilled by copying the binaries off of this image into the peer container at init time: 

```yaml 
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: org1-peer1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: org1-peer1
  template:
    metadata:
      labels:
        app: org1-peer1
    spec:
      containers:
        - name: main
          image: hyperledger/fabric-peer:{{FABRIC_VERSION}}
          imagePullPolicy: IfNotPresent
          envFrom:
            - configMapRef:
                name: org1-peer1-config
          ports:
            - containerPort: 7051
            - containerPort: 7052
            - containerPort: 9443
          volumeMounts:
            - name: fabric-volume
              mountPath: /var/hyperledger
            - name: fabric-config
              mountPath: /var/hyperledger/fabric/config
            - name: ccs-builder
              mountPath: /var/hyperledger/fabric/chaincode/ccs-builder/bin

      # load the external chaincode builder into the peer image prior to peer launch.
      initContainers:
        - name: fabric-ccs-builder
          image: hyperledgendary/fabric-ccs-builder
          imagePullPolicy: IfNotPresent
          command: [sh, -c]
          args: ["cp /go/bin/* /var/hyperledger/fabric/chaincode/ccs-builder/bin/"]
          volumeMounts:
            - name: ccs-builder
              mountPath: /var/hyperledger/fabric/chaincode/ccs-builder/bin

      volumes:
        - name: fabric-volume
          persistentVolumeClaim:
            claimName: fabric-org1
        - name: fabric-config
          configMap:
            name: org1-config
        - name: ccs-builder
          emptyDir: {}

```


