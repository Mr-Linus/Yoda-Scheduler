package collection

import (
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/filter"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"sync"
)

type Data struct {
	Value uint
}

func (d *Data) Clone() framework.StateData {
	clone := Data{Value:d.Value}
	return &clone
}

func CollectMaxValue(value string,state *framework.CycleState,nodes []*v1.Node,filteredNodesStatuses framework.NodeToStatusMap) *framework.Status {
	Max := Data{Value:0}
	for _,n := range nodes{
		if filteredNodesStatuses[n.GetName()].IsSuccess(){
			if filter.StrToUInt(n.Labels["scv/"+value]) > Max.Value{
				Max.Value = filter.StrToUInt(n.Labels["scv/FreeMemory"])
			}
		}
	}
	if Max.Value == 0{
		return framework.NewStatus(framework.Error," The max "+value+" of the nodes is 0")
	}
	state.Lock()
	state.Write(framework.StateKey("Max"+value), &Max)
	state.Lock()
	return framework.NewStatus(framework.Success,"")
}


func ParallelCollection(state *framework.CycleState,nodes []*v1.Node,filteredNodesStatuses framework.NodeToStatusMap) *framework.Status{
	var (
		stop <-chan struct{}
		mx sync.RWMutex
		msg = ""
	)
	sum := []string{"Cores","FreeMemory","Bandwidth","MemoryClock"}
	pieces := len(sum)
	toProcess := make(chan string, pieces)
	for _,v := range sum{
		toProcess <- v
	}
	close(toProcess)
	wg := sync.WaitGroup{}
	wg.Add(len(sum))
	for i := 0;i <= pieces; i++{
		go func() {
			defer wg.Done()
			for value := range toProcess{
				select {
					case <- stop:
						return
					default:
						if re := CollectMaxValue(value,state,nodes,filteredNodesStatuses);re.IsSuccess(){
							klog.V(3).Infof(re.Message())
							mx.Lock()
							msg += re.Message()
							mx.Unlock()
						}
				}
			}
		}()
	}
	wg.Wait()
	if msg != ""{
		return framework.NewStatus(framework.Error,msg)
	}
	return framework.NewStatus(framework.Success,"")
}