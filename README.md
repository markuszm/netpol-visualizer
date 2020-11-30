# netpol-visualizer

Initialized using operator-sdk.


API Design - Example
https://book.kubebuilder.io/cronjob-tutorial/api-design.html

Controller Design - Example:
https://book.kubebuilder.io/cronjob-tutorial/controller-implementation.html

## Install
Install the following tools:
- helm
- helmfile
- minikube

## Setup

Run the commands:
```shell script
minikube start
kubectl config use-context minikube
```

Set up Docker and CRDs:
```shell script
make install
minikube docker-env --shell=<your-shell>
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

Go to http://localhost:7474 and enter the following parameters into the *Connect to Neo4J* window:
- Connect URL: **neo4j://localhost:7687**
- Username: **neo4j**
- Password: **secret**

### Controller
Run the operator locally:
```shell script
make run LOCAL=1
```

Deploy the controller to the cluster:
```shell script
make docker-build deploy
```
