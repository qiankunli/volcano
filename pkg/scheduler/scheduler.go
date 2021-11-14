/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scheduler

import (
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"sync"
	"time"
	"volcano.sh/volcano/pkg/scheduler/conf"
	"volcano.sh/volcano/pkg/scheduler/plugins"

	schedcache "volcano.sh/volcano/pkg/scheduler/cache"
	"volcano.sh/volcano/pkg/scheduler/framework"
	"volcano.sh/volcano/pkg/scheduler/metrics"
)

// Scheduler watches for new unscheduled pods for volcano. It attempts to find
// nodes that they fit on and writes bindings back to the api server.
type Scheduler struct {
	cache           schedcache.Cache
	schedulerConfig *conf.SchedulerConfig
	schedulePeriod  time.Duration
	mutex           sync.Mutex
}

// NewScheduler returns a scheduler
func NewScheduler(
	config *rest.Config,
	schedulerName string,
	schedulerConf string,
	period time.Duration,
	defaultQueue string,
) (*Scheduler, error) {

	schedulerConfig, err := conf.NewSchedulerConf(schedulerConf)
	if err != nil {
		return nil, err
	}
	scheduler := &Scheduler{
		schedulerConfig: schedulerConfig,
		schedulePeriod:  period,
		cache:           schedcache.New(config, schedulerName, schedulerConfig, defaultQueue),
	}
	handler := func(sc *conf.SchedulerConfig) {
		plugins.ApplyPluginConf(sc.Plugins)
	}
	schedulerConfig.AddConfigChangeHandler(handler)
	return scheduler, nil
}

// Run runs the Scheduler
func (pc *Scheduler) Run(stopCh <-chan struct{}) {
	go pc.schedulerConfig.WatchSchedulerConf(stopCh)
	// Start cache for policy.
	go pc.cache.Run(stopCh)
	pc.cache.WaitForCacheSync(stopCh)
	klog.V(2).Infof("scheduler completes Initialization and start to run")
	go wait.Until(pc.runOnce, pc.schedulePeriod, stopCh)
}

func (pc *Scheduler) runOnce() {
	klog.V(4).Infof("Start scheduling ...")
	scheduleStartTime := time.Now()
	defer klog.V(4).Infof("End scheduling ...")

	pc.mutex.Lock()
	actionNames := pc.schedulerConfig.Actions
	plugins := pc.schedulerConfig.Plugins
	configurations := pc.schedulerConfig.Configurations
	pc.mutex.Unlock()

	ssn := framework.OpenSession(pc.cache, plugins, configurations)
	defer framework.CloseSession(ssn)

	for _, actionName := range actionNames {
		if action, found := framework.GetAction(actionName); found {
			actionStartTime := time.Now()
			action.Execute(ssn)
			metrics.UpdateActionDuration(action.Name(), metrics.Duration(actionStartTime))
		}
	}
	metrics.UpdateE2eDuration(metrics.Duration(scheduleStartTime))
}
