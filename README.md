# PodChaosMonkey

PodChaosMonkey is a basic chaos engineering tool that randomly deletes a pod matching a given namespace and annotations at a specified interval.
It runs a kubernetes deployment and can be configured using the following environment variables in `./helm/pod-chaos-monkey/values.yaml`:

```
env:
  - name: "PODCHAOSMONKEY_ANNOTATIONS"
    value: "podchaosmonkey=true"
  - name: "PODCHAOSMONKEY_INTERVAL_SECONDS"
    value: "60"
  - name: "PODCHAOSMONKEY_NAMESPACE"
    value: "workloads"
  - name: "PODCHAOSMONKEY_GRACE_PERIOD_SECONDS"
    value: "10"
```

## Build

In order to build the project run `make` in this directory.
The only dependencies are `make` and `docker`. Golang doesn't need to be installed on the host as the binary is built on docker itself. Look at `Dockerfile` for more details.
The makefile will build an image with the following tag: `podchaosmonkey:latest`.

## Install
```
helm install pod-chaos-monkey ./helm/pod-chaos-monkey --namespace pod-chaos-monkey --create-namespace
```

The above command is just an example - any namespace can be used to install PodChaosMonkey - RBAC resources will be created by the chart, accordingly.

## Test
Unit tests are being run as part of the build.
The helm charts available under `./helm` can be used to simualte a real-life scenario.

```
helm install test-workload ./helm/test-workload --namespace workloads --create-namespace
```

will start a k8s deployment running nginx in the namespace (`workloads`) monitored by default by PodChaosMonkey. The deployment also has the annotation that makes PodChaosMonkey consider it for deletion:

```
podAnnotations:
  podchaosmonkey: "true"
```

Feel free to spin multiple instances of the chart or remove the annotation to test the behaviour.

## TODO / ideas for improvement
- fix hardcoded namespace in cluster role binding
- think about how to structure rbac stuff in charts
- add more options for filtering pods - age, state, exlusion, etc.
- support more complex scheduling options - only run on certain days of the week, or during a specified time window during weekdays
- prometheus metrics endpoint
- readiness/liveness probes
- testing beyond the "happy path"
- unit tests for random generator
