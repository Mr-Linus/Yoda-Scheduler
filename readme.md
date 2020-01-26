## Yoda-Scheduler
Yoda is a kubernetes scheduler based on [scheduling-framework](https://github.com/kubernetes/enhancements/blob/master/keps/sig-scheduling/20180409-scheduling-framework.md). By cooperation with [SCV Sniffer](https://github.com/NJUPT-ISL/SCV),
 it is schedules tasks according to GPU metrics.


[![Go Report Card](https://goreportcard.com/badge/github.com/NJUPT-ISL/Yoda-Scheduler)](https://goreportcard.com/report/github.com/NJUPT-ISL/Yoda-Scheduler)

### Get Started 
- Make sure SCV sniffer is deployed in kubernetes cluster: [SCV: Get-Started](https://github.com/NJUPT-ISL/SCV#get-started)

- Deploy Yoda Scheduler:
```shell
kubectl apply -f https://raw.githubusercontent.com/NJUPT-ISL/Yoda-Scheduler/master/deploy/deploy.yaml
```

- Check the Yoda Scheduler Status:
```shell
kubectl get pods -n kube-system 
```

- Deploy a sample pod using Yoda:
```shell
kubectl apply -f https://raw.githubusercontent.com/NJUPT-ISL/Yoda-Scheduler/master/deploy/test-deployment.yaml
```

- Check the sample pod Status:
```shell
kubectl get pods 
```
