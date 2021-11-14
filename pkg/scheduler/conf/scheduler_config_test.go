/*
Copyright 2019 The Kubernetes Authors.

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

package conf

import (
	"reflect"
	"testing"
)

func TestLoadSchedulerConf(t *testing.T) {
	configuration := `
actions: "enqueue, allocate, backfill"
tiers:
- plugins:
  - name: priority
  - name: gang
  - name: conformance
- plugins:
  - name: drf
  - name: predicates
  - name: proportion
  - name: nodeorder
nodeSelector:
  app: app1
  demo: demo1,demo2
`

	trueValue := true
	expectedTiers := []Tier{
		{
			Plugins: []PluginOption{
				{
					Name:                  "priority",
					EnabledJobOrder:       &trueValue,
					EnabledNamespaceOrder: &trueValue,
					EnabledJobReady:       &trueValue,
					EnabledJobPipelined:   &trueValue,
					EnabledTaskOrder:      &trueValue,
					EnabledPreemptable:    &trueValue,
					EnabledReclaimable:    &trueValue,
					EnabledQueueOrder:     &trueValue,
					EnabledPredicate:      &trueValue,
					EnabledBestNode:       &trueValue,
					EnabledNodeOrder:      &trueValue,
					EnabledTargetJob:      &trueValue,
					EnabledReservedNodes:  &trueValue,
					EnabledJobEnqueued:    &trueValue,
					EnabledVictim:         &trueValue,
					EnabledJobStarving:    &trueValue,
				},
				{
					Name:                  "gang",
					EnabledJobOrder:       &trueValue,
					EnabledNamespaceOrder: &trueValue,
					EnabledJobReady:       &trueValue,
					EnabledJobPipelined:   &trueValue,
					EnabledTaskOrder:      &trueValue,
					EnabledPreemptable:    &trueValue,
					EnabledReclaimable:    &trueValue,
					EnabledQueueOrder:     &trueValue,
					EnabledPredicate:      &trueValue,
					EnabledBestNode:       &trueValue,
					EnabledNodeOrder:      &trueValue,
					EnabledTargetJob:      &trueValue,
					EnabledReservedNodes:  &trueValue,
					EnabledJobEnqueued:    &trueValue,
					EnabledVictim:         &trueValue,
					EnabledJobStarving:    &trueValue,
				},
				{
					Name:                  "conformance",
					EnabledJobOrder:       &trueValue,
					EnabledNamespaceOrder: &trueValue,
					EnabledJobReady:       &trueValue,
					EnabledJobPipelined:   &trueValue,
					EnabledTaskOrder:      &trueValue,
					EnabledPreemptable:    &trueValue,
					EnabledReclaimable:    &trueValue,
					EnabledQueueOrder:     &trueValue,
					EnabledPredicate:      &trueValue,
					EnabledBestNode:       &trueValue,
					EnabledNodeOrder:      &trueValue,
					EnabledTargetJob:      &trueValue,
					EnabledReservedNodes:  &trueValue,
					EnabledJobEnqueued:    &trueValue,
					EnabledVictim:         &trueValue,
					EnabledJobStarving:    &trueValue,
				},
			},
		},
		{
			Plugins: []PluginOption{
				{
					Name:                  "drf",
					EnabledJobOrder:       &trueValue,
					EnabledNamespaceOrder: &trueValue,
					EnabledJobReady:       &trueValue,
					EnabledJobPipelined:   &trueValue,
					EnabledTaskOrder:      &trueValue,
					EnabledPreemptable:    &trueValue,
					EnabledReclaimable:    &trueValue,
					EnabledQueueOrder:     &trueValue,
					EnabledPredicate:      &trueValue,
					EnabledBestNode:       &trueValue,
					EnabledNodeOrder:      &trueValue,
					EnabledTargetJob:      &trueValue,
					EnabledReservedNodes:  &trueValue,
					EnabledJobEnqueued:    &trueValue,
					EnabledVictim:         &trueValue,
					EnabledJobStarving:    &trueValue,
				},
				{
					Name:                  "predicates",
					EnabledJobOrder:       &trueValue,
					EnabledNamespaceOrder: &trueValue,
					EnabledJobReady:       &trueValue,
					EnabledJobPipelined:   &trueValue,
					EnabledTaskOrder:      &trueValue,
					EnabledPreemptable:    &trueValue,
					EnabledReclaimable:    &trueValue,
					EnabledQueueOrder:     &trueValue,
					EnabledPredicate:      &trueValue,
					EnabledBestNode:       &trueValue,
					EnabledNodeOrder:      &trueValue,
					EnabledTargetJob:      &trueValue,
					EnabledReservedNodes:  &trueValue,
					EnabledJobEnqueued:    &trueValue,
					EnabledVictim:         &trueValue,
					EnabledJobStarving:    &trueValue,
				},
				{
					Name:                  "proportion",
					EnabledJobOrder:       &trueValue,
					EnabledNamespaceOrder: &trueValue,
					EnabledJobReady:       &trueValue,
					EnabledJobPipelined:   &trueValue,
					EnabledTaskOrder:      &trueValue,
					EnabledPreemptable:    &trueValue,
					EnabledReclaimable:    &trueValue,
					EnabledQueueOrder:     &trueValue,
					EnabledPredicate:      &trueValue,
					EnabledBestNode:       &trueValue,
					EnabledNodeOrder:      &trueValue,
					EnabledTargetJob:      &trueValue,
					EnabledReservedNodes:  &trueValue,
					EnabledJobEnqueued:    &trueValue,
					EnabledVictim:         &trueValue,
					EnabledJobStarving:    &trueValue,
				},
				{
					Name:                  "nodeorder",
					EnabledJobOrder:       &trueValue,
					EnabledNamespaceOrder: &trueValue,
					EnabledJobReady:       &trueValue,
					EnabledJobPipelined:   &trueValue,
					EnabledTaskOrder:      &trueValue,
					EnabledPreemptable:    &trueValue,
					EnabledReclaimable:    &trueValue,
					EnabledQueueOrder:     &trueValue,
					EnabledPredicate:      &trueValue,
					EnabledBestNode:       &trueValue,
					EnabledNodeOrder:      &trueValue,
					EnabledTargetJob:      &trueValue,
					EnabledReservedNodes:  &trueValue,
					EnabledJobEnqueued:    &trueValue,
					EnabledVictim:         &trueValue,
					EnabledJobStarving:    &trueValue,
				},
			},
		},
	}
	expectedNodeSelector := map[string]string{
		"app:app1":   "",
		"demo:demo1": "",
		"demo:demo2": "",
	}

	var expectedConfigurations []Configuration

	_, tiers, configurations, nodeSelector, err := unmarshalSchedulerConf(configuration)
	if err != nil {
		t.Errorf("Failed to load scheduler configuration: %v", err)
	}
	if !reflect.DeepEqual(tiers, expectedTiers) {
		t.Errorf("Failed to set default settings for plugins, expected: %+v, got %+v",
			expectedTiers, tiers)
	}
	if !reflect.DeepEqual(configurations, expectedConfigurations) {
		t.Errorf("Wrong configuration, expected: %+v, got %+v",
			expectedConfigurations, configurations)
	}
	if !reflect.DeepEqual(nodeSelector, expectedNodeSelector) {
		t.Errorf("Wrong nodeSelector, expected: %+v, got %+v",
			expectedConfigurations, configurations)
	}
}
