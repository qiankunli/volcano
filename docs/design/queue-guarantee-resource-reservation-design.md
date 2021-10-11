# Volcano Resource Reservation For Queue

@[qiankunli](https://github.com/qiankunli); Oct 11rd, 2021

## Motivation
As [issue 1101](https://github.com/volcano-sh/volcano/issues/1101) mentioned, Volcano should support resource reservation
for specified queue. Requirement detail as follows:
* Support reserving specified resources for specified queue
* We only Consider non-preemption reservation. 
* Support enable and disable resource reservation for specified queue dynamically without restarting Volcano.
* Support hard reservation resource specified and percentage reservation resource specified.

@[Thor-wl](https://github.com/Thor-wl) already provide a design doc [Volcano Resource Reservation For Queue](https://github.com/volcano-sh/volcano/blob/master/docs/design/queue-resource-reservation-design.md)
I do not implement all features above, supported feature are as follows:

* Support reserving specified resources for specified queue
* We only Consider non-preemption reservation.
* Support enable and disable resource reservation for specified queue dynamically without restarting Volcano.
* Support hard reservation resource specified

## Consideration
### Resource Request
* The reserved resource cannot be more than the total resource amount of cluster at all dimensions.
* If `capability` is set in a queue, the reserved resource must be no more than it at all dimensions.

### Safety
* Malicious application for large amount of resource reservation will cause jobs in other queue to block.

## Design
### API
```
apiVersion: scheduling.volcano.sh/v1beta1
kind: Queue
metadata:
  name: q1
spec:
  reclaimable: true
  weight: 1
  guarantee:             // reservation key word
    resource:            // specified reserving resource
      cpu: 2c
      memory: 4G
```

`guarantee.resource` list of reserving resource categories and amount. 

## Implementation

I think we can support `guarantee.resource` in proportion plugin whose code in `pkg/scheduler/plugins/proportion/proportion.go`.if there are queue1 and queue2 in cluster(memory=1000M)

|queue|spec.capacity.memory|spec.guarantee.resource.memory|job|
|---|---|---|---|
|queue1|nil|nil|job1|
|queue2|nil|200M|job2|

1. firstly, If `guarantee.resource` is set in a queue. it means all queue in cluster has a `capacity`. for example,if `queue2.spec.guarantee.resource.memory` is 200M , it means `queue1.spec.capacity.memory` is 800M. queue1 will never use the 200M memory by setting the capacity, even though `queue1.capacity` is nil.
   1. old logic:  `queue1.capacity=queue1.spec.capacity`
   2. new logic: `queue1.capacity=min(queue1.spec.capacity,cluster.resource - queue2.guarantee.resource)`
   3. use new `queue.capacity` replace the code `queue.spec.capacity` anywhere
2. `proportionPlugin.OnSessionOpen` will calculate `deserved resource` and `allocated resource` for a queue. if  `queue.deserved` > `queue.allocated + job.request`, it will allocate resource owned by queue for a job. we should change formula  of `queue.deserved`, so that queue2  reserve the 200M memory, even though there is no job in queue2.
   2. old logic: `queue2.request = job2.request` and `queue2.deserved = min(job2.request,queue2.capacity,queue2.weight-resource)`.
   3. new logic: we should modify `queue2.deserved` from old value,  `queue2.deserved = max(queue2.deserved,queue2.spec.guarantee.resource)`



