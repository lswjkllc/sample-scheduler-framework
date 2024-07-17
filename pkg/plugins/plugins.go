package plugins

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	fruntime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
)

// 插件名称
const Name = "sample-plugin"

type Args struct {
	FavoriteColor  string `json:"favorite_color,omitempty"`
	FavoriteNumber int    `json:"favorite_number,omitempty"`
	ThanksTo       string `json:"thanks_to,omitempty"`
}

type Sample struct {
	args   *Args
	handle framework.Handle
}

func (s *Sample) Name() string {
	return Name
}

func (s *Sample) PreFilter(ctx context.Context, pod *v1.Pod) *framework.Status {
	klog.V(3).Infof("prefilter pod: %v", pod.Name)
	return framework.NewStatus(framework.Success, "")
}

func (s *Sample) Filter(ctx context.Context, pod *v1.Pod, nodeName string) *framework.Status {
	klog.V(3).Infof("filter pod: %v, node: %v", pod.Name, nodeName)
	return framework.NewStatus(framework.Success, "")
}

func (s *Sample) PreBind(ctx context.Context, pod *v1.Pod, nodeName string) *framework.Status {
	if nodeInfo, err := s.handle.SnapshotSharedLister().NodeInfos().Get(nodeName); err != nil {
		return framework.NewStatus(
			framework.Error,
			fmt.Sprintf("prebind get node[%s] info error: %v", nodeName, err))
	} else {
		klog.V(3).Infof("prebind node info: %+v", nodeInfo.Node())
		return framework.NewStatus(framework.Success, "")
	}
}

// type PluginFactory = func(configuration *runtime.Unknown, f FrameworkHandle) (Plugin, error)
func New(plArgs runtime.Object, h framework.Handle) (framework.Plugin, error) {
	args := &Args{}
	if err := fruntime.DecodeInto(plArgs, args); err != nil {
		return nil, err
	}
	klog.V(3).Infof("get plugin config args: %+v", args)
	return &Sample{
		args:   args,
		handle: h,
	}, nil
}
