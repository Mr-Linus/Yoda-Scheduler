package collection

import (
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/filter"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"sync"
)

const Workers int = 4

type Data struct {
	Value int64
}

func (s *Data) Clone() framework.StateData {
	c := &Data{
		Value: s.Value,
	}
	return c
}

var Sum = []string{"Cores","FreeMemory","Bandwidth","MemoryClock","MemorySum","Number","Memory"}


func CollectMaxValue(value string,state *framework.CycleState,nodes []*v1.Node,filteredNodesStatuses framework.NodeToStatusMap) *framework.Status {
	Max := Data{Value: 0}
	for _,n := range nodes{
		if filteredNodesStatuses[n.GetName()].IsSuccess(){
			if filter.StrToInt64(n.Labels["scv/"+value]) > Max.Value{
				Max.Value = filter.StrToInt64(n.Labels["scv/FreeMemory"])
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


func ParallelCollection(workers int,state *framework.CycleState,nodes []*v1.Node,filteredNodesStatuses framework.NodeToStatusMap) *framework.Status{
	var (
		stop <-chan struct{}
		mx sync.RWMutex
		msg = ""
	)
	pieces := len(Sum)
	toProcess := make(chan string, pieces)
	for _,v := range Sum{
		toProcess <- v
	}
	close(toProcess)
	if pieces > workers {
		workers = pieces
	}
	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0;i <= workers; i++{
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