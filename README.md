# netpol-visualizer

Initialized using operator-sdk.


API Design - Example
https://book.kubebuilder.io/cronjob-tutorial/api-design.html

Controller Design - Example:
https://book.kubebuilder.io/cronjob-tutorial/controller-implementation.html

API not used at the moment.

## Setup
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
kubectl port-forward -n neo4j svc/neo4j-neo4j 7474:7474 7687:7687
```

Go to http://localhost:7474 and enter the following parameters into the *Connect to Neo4J* window:
- Connect URL: **neo4j://localhost:7687**
- Username: **neo4j**
- Password: **secret**

In http://localhost:7474/browser > Settings > Graph Visualization, uncheck the setting "Connect result nodes".

Enter Cipher queries to analyze the network policy graph.

### Controller
Run the operator locally:
```shell script
make run LOCAL=1
```

Deploy the controller to the cluster:
```shell script
make docker-build deploy
```


## Example queries

Show accessible network paths:

```
MATCH (x:Pod)-[e1:INGRESS_ALLOWED*1]->(y:Pod)<-[e2:EGRESS_ALLOWED*1]-(x) RETURN x,e1,e2,y
```

Create new edges for each accessible network path:

```
MATCH (x:Pod)-[:INGRESS_ALLOWED*1]->(y:Pod) MATCH (x)-[:EGRESS_ALLOWED*1]->(y) CREATE (x)-[e:CAN_ACCESS]->(y) RETURN x,e,y
```

## TODO

- in namespaces where no NetworkPolicy is present, all traffic should be allowed (as defined in the K8S API spec)
- on change of NetworkPolicy or Pod, delete all nodes and edges
