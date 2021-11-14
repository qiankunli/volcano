package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog"
	"path/filepath"
	"strings"
	"sync"
	"volcano.sh/volcano/pkg/filewatcher"
)

var defaultSchedulerConf = `
actions: "enqueue, allocate, backfill"
tiers:
- plugins:
  - name: priority
  - name: gang
  - name: conformance
- plugins:
  - name: overcommit
  - name: drf
  - name: predicates
  - name: proportion
  - name: nodeorder
`

type SchedulerConfig struct {
	once          sync.Once
	schedulerConf string
	fileWatcher   filewatcher.FileWatcher
	// should be private
	Plugins        []Tier
	Configurations []Configuration
	Actions        []string
	NodeSelector   map[string]string

	handlers []ConfigChangeHandler
	mutex    sync.Mutex
}
type ConfigChangeHandler func(sc *SchedulerConfig)

func NewSchedulerConf(schedulerConf string) (*SchedulerConfig, error) {
	var watcher filewatcher.FileWatcher
	if schedulerConf != "" {
		var err error
		path := filepath.Dir(schedulerConf)
		watcher, err = filewatcher.NewFileWatcher(path)
		if err != nil {
			return nil, fmt.Errorf("failed creating filewatcher for %s: %v", schedulerConf, err)
		}
	}

	schedulerConfig := &SchedulerConfig{
		schedulerConf: schedulerConf,
		fileWatcher:   watcher,
		handlers:      make([]ConfigChangeHandler, 0),
	}

	schedulerConfig.LoadSchedulerConf()
	return schedulerConfig, nil
}

func (sc *SchedulerConfig) LoadSchedulerConf() {
	var err error
	sc.once.Do(func() {
		sc.Actions, sc.Plugins, sc.Configurations, sc.NodeSelector, err = unmarshalSchedulerConf(defaultSchedulerConf)
		if err != nil {
			klog.Errorf("unmarshal scheduler config %s failed: %v", defaultSchedulerConf, err)
			panic("invalid default configuration")
		}
	})

	var config string
	if len(sc.schedulerConf) != 0 {
		if config, err = readSchedulerConf(sc.schedulerConf); err != nil {
			klog.Errorf("Failed to read scheduler configuration '%s', using previous configuration: %v",
				sc.schedulerConf, err)
			return
		}
	}

	actions, plugins, configurations, nodeSelector, err := unmarshalSchedulerConf(config)
	if err != nil {
		klog.Errorf("scheduler config %s is invalid: %v", config, err)
		return
	}
	sc.mutex.Lock()
	// If it is valid, use the new configuration
	sc.Actions = actions
	sc.Plugins = plugins
	sc.Configurations = configurations
	sc.NodeSelector = nodeSelector
	sc.mutex.Unlock()
}

func (sc *SchedulerConfig) WatchSchedulerConf(stopCh <-chan struct{}) {
	if sc.fileWatcher == nil {
		return
	}
	eventCh := sc.fileWatcher.Events()
	errCh := sc.fileWatcher.Errors()
	for {
		select {
		case event, ok := <-eventCh:
			if !ok {
				return
			}
			klog.V(4).Infof("watch %s event: %v", sc.schedulerConf, event)
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				sc.LoadSchedulerConf()
				if len(sc.handlers) > 0 {
					for _, handler := range sc.handlers {
						handler(sc)
					}
				}
			}
		case err, ok := <-errCh:
			if !ok {
				return
			}
			klog.Infof("watch %s error: %v", sc.schedulerConf, err)
		case <-stopCh:
			return
		}
	}
}

func (sc *SchedulerConfig) AddConfigChangeHandler(handler ConfigChangeHandler) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.handlers = append(sc.handlers, handler)
}

func unmarshalSchedulerConf(confStr string) ([]string, []Tier, []Configuration, map[string]string, error) {
	schedulerConf := &SchedulerConfiguration{}

	if err := yaml.Unmarshal([]byte(confStr), schedulerConf); err != nil {
		return nil, nil, nil, nil, err
	}
	// Set default settings for each plugin if not set
	for _, tier := range schedulerConf.Tiers {
		// drf with hierarchy enabled
		hdrf := false
		// proportion enabled
		proportion := false
		for j := range tier.Plugins {
			if tier.Plugins[j].Name == "drf" &&
				tier.Plugins[j].EnabledHierarchy != nil &&
				*tier.Plugins[j].EnabledHierarchy {
				hdrf = true
			}
			if tier.Plugins[j].Name == "proportion" {
				proportion = true
			}
		}
		if hdrf && proportion {
			return nil, nil, nil, nil, fmt.Errorf("proportion and drf with hierarchy enabled conflicts")
		}
	}

	var actions []string
	actionNames := strings.Split(schedulerConf.Actions, ",")
	for _, actionName := range actionNames {
		actionName = strings.TrimSpace(actionName)
		if len(actionName) > 0 {
			actions = append(actions, actionName)
		}
	}

	if len(schedulerConf.NodeSelector) == 0 {
		return actions, schedulerConf.Tiers, schedulerConf.Configurations, nil, nil
	}
	nodeSelectorLabels := make(map[string]string)
	for labelName, labelValue := range schedulerConf.NodeSelector {
		labelValues := strings.Split(labelValue, ",")
		for _, lv := range labelValues {
			lv = strings.TrimSpace(lv)
			if len(lv) > 0 {
				key := labelName + ":" + lv
				nodeSelectorLabels[key] = ""
			}
		}
	}
	return actions, schedulerConf.Tiers, schedulerConf.Configurations, nodeSelectorLabels, nil
}

func readSchedulerConf(confPath string) (string, error) {
	dat, err := ioutil.ReadFile(confPath)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}
