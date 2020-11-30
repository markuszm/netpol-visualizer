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
Deploy or update the helm releases:
```shell script
helmfile apply
```

...
