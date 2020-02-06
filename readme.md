# Yoda-Scheduler

![Status](https://github.com/NJUPT-ISL/Yoda-Scheduler/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/NJUPT-ISL/Yoda-Scheduler)](https://goreportcard.com/report/github.com/NJUPT-ISL/Yoda-Scheduler)

![Yoda](./img/yoda.png)
----

Yoda is a kubernetes scheduler based on [scheduling-framework](https://github.com/kubernetes/enhancements/blob/master/keps/sig-scheduling/20180409-scheduling-framework.md). By cooperation with [SCV Sniffer](https://github.com/NJUPT-ISL/SCV),
 it is schedules tasks according to GPU metrics.
 

 
## Get Started 

- Make sure kubernetes cluster version is `1.17+` and SCV sniffer is deployed in kubernetes cluster: [SCV: Get-Started](https://github.com/NJUPT-ISL/SCV#get-started)

- Deploy Yoda Scheduler:
```shell
kubectl apply -f https://raw.githubusercontent.com/NJUPT-ISL/Yoda-Scheduler/master/deploy/deploy.yaml
```

- Check the Yoda Scheduler Status:
```shell
kubectl get pods -n kube-system 
```
## Deploy a Pod using Yoda
- Deploy a sample pod using Yoda:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test
  labels:
    app: test
spec:
  schedulerName: yoda-scheduler
  containers:
    - image: nginx
      imagePullPolicy: IfNotPresent
      name: nginx
      ports:
        - containerPort: 80
```

- Check the sample pod Status:
```shell
kubectl get pods 
```
## How to develop Yoda
- Compile yoda-scheduler:
```shell
make local
```
- Build the docker image:
```shell
make build
```
- Clean the Build file
```shell
make clean
```