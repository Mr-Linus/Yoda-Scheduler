## Yoda-Scheduler
Yoda is a scheduler based on GPU metrics. By cooperation with [SCV Sniffer](https://github.com/NJUPT-ISL/SCV),
 It is schedules tasks according to GPU performance.

### Get Started 
- Make sure SCV sniffer is deployed in kubernetes cluster: [SCV:Get Started](https://github.com/NJUPT-ISL/SCV#get-started)

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

```