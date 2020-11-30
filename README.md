# netpol-visualizer

Initialized using operator-sdk.


API Design - Example
https://book.kubebuilder.io/cronjob-tutorial/api-design.html

Controller Design - Example:
https://book.kubebuilder.io/cronjob-tutorial/controller-implementation.html

## Setup
Install the following tools:
- helm
- helmfile
- minikube

Run the commands:
```shell script
minikube start
kubectl config use-context minikube
```

## Run

### Helmfile

Deploy or update the helm releases:
```shell script
helmfile apply
```

### Neo4J
Get access to running Neo4J instance:
```shell script
kubectl port-forward svc/neo4j-neo4j 7474:7474 7687:7687
```

Log in at http://localhost:7474 with username *neo4j* and password *secret*.
...
