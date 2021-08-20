# fabric-ccs-builder

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


## Docker image 

Build a docker image embedding the `build`, `release`, and `detect` binaries into the /go/bin directory.  For kube-native environments, the 
fabric-ccs-builder image may be deployed as a sidecar in the peers.  When registered as an external builder in the peer / network configuration, 
this allows us to deploy external chaincode as a service while using the base fabric images.

```shell 
docker build \
  -t hyperledgendary/fabric-ccs-builder \
  .
```

