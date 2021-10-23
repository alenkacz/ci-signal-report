# CI signal report

You can get the current overview for CI signal report by running

```bash
GITHUB_AUTH_TOKEN=xxx go run ./cmd/ci-reporter.go
```

It needs a GitHub token to be able to query the project board for CI signal. For some reason even though those boards are available for public view, the APIs require auth. See [this documentation](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line) to set up your access token.

## Prerequisites

- GoLang >=1.16

## Run the report

```bash
git clone git@github.com:leonardpahlke/ci-signal-report.git <folder>
cd <folder>
GITHUB_AUTH_TOKEN=xxx go run main.go
```

### Other version statistics

By adding `RELEASE_VERSION=xxx` where the XXX can be like `1.23`, the report statistics get extended for the chosen version.

```bash
GITHUB_AUTH_TOKEN=xxx RELEASE_VERSION=xxx go run ./cmd/ci-reporter.go
```

### Short report

You can also output a short version of the report with the flag `-short`. This reduces the report to `New/Not Yet Started` and `In Flight` issues.

```bash
GITHUB_AUTH_TOKEN=xxx go run ./cmd/ci-reporter.go -short
```

## Rate limits

GitHub API has rate limits, to see how much you have used you can query like this (replace User with your GH user and Token with your Auth Token):

```bash
curl \
  -u GIT_HUB-USER:GIT_HUB_TOKEN -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/rate_limit & curl \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/rate_limit
```

## Example output

```bash
GITHUB_AUTH_TOKEN=yourFavoriteGitHubTokenLivesHere RELEASE_VERSION=1.21 go run ./cmd/ci-reporter.go -short

----------
New/Not Yet Started
SIG Cloud-Provider
#100230 https://github.com/kubernetes/kubernetes/issues/100230 [Flaky Test] [sig-cloud-provider-gcp] Nodes [Disruptive] Resize [Slow] should be able to delete nodes


----------
In flight
SIG Release
#1693 https://github.com/kubernetes/release/issues/1693 push-build.sh container image pushes should precede staging GCS artifacts and writing version markers

SIG Network
#93740 https://github.com/kubernetes/kubernetes/issues/93740 [Flaky Test][sig-network] Loadbalancing: L7 GCE [Slow] [Feature:Ingress] should conform to Ingress spec

SIG Api-Machinery
#100112 https://github.com/kubernetes/kubernetes/issues/100112 [flaky test] k8s.io/kubernetes/pkg/registry/core/endpoint/storage.TestWatch

SIG Apps
#100314 https://github.com/kubernetes/kubernetes/issues/100314 [Flaky Test] Kubernetes e2e suite: [sig-apps] Deployment iterative rollouts should eventually progress

SIG Scalability
#97071 https://github.com/kubernetes/kubernetes/issues/97071 [Flaky test] [sig-storage] In-tree Volumes [Driver: gcepd] [Testpattern: Pre-provisioned PV (xfs)][Slow] volumes should store data


----------
Observing
SIG Api-Machinery
#97312 https://github.com/kubernetes/kubernetes/issues/97312 [Flaky Test] go_test: //staging/src/k8s.io/apiserver/pkg/server/go_default_test:run_2_of_2

SIG Instrumentation
#95366 https://github.com/kubernetes/kubernetes/issues/95366 [Flaky Test] [sig-instrumentation] MetricsGrabber should grab all metrics from a Scheduler

SIG Node
#99437 https://github.com/kubernetes/kubernetes/issues/99437 [Flake][sig-node] Pods should run through the lifecycle of Pods and PodStatus

SIG
#102300 https://github.com/kubernetes/kubernetes/issues/102300  kubernetes.up e2e-up.sh

SIG Apps
#98501 https://github.com/kubernetes/kubernetes/issues/98501 [Flaky Test] [sig-apps] CronJob should delete failed finished jobs with limit of one job

SIG Scalability
#103688 https://github.com/kubernetes/kubernetes/issues/103688  [sig-node] ci-kubernetes-e2e-gce-scale-correctness

SIG cluster-lifecycle
#102345 https://github.com/kubernetes/kubernetes/issues/102345 [sig-cluster-lifecycle] ci-kubernetes-e2e-kubeadm-kinder-latest-on-1-21

SIG CLI
#98854 https://github.com/kubernetes/kubernetes/issues/98854 [Flaky Test] [[sig-cli] Kubectl client Simple pod should return command exit codes

SIG Scheduling
#98857 https://github.com/kubernetes/kubernetes/issues/98857 [Flaky Test]  [sig-scheduling] SchedulerPredicates [Serial] validates resource limits of pods that are allowed to run [Conformance]

SIG Network
#102006 https://github.com/kubernetes/kubernetes/issues/102006 L7 GCE Ingress failing tests

SIG Storage
#102077 https://github.com/kubernetes/kubernetes/issues/102077 e2e failures: all pdcsi tests

SIG Windows
#101906 https://github.com/kubernetes/kubernetes/issues/101906 [Flaking Test] [sig-node] ConfigMap should be consumable via the environment

SIG Cloud-Provider
#102904 https://github.com/kubernetes/kubernetes/issues/102904 "dial timeout" flakes in e2e tests with konnectivity proxy enabled


----------
Resolved
SIG Scalability
#100621 https://github.com/kubernetes/kubernetes/issues/100621 gce-master-scale-performance is failing

SIG Apps
#100551 https://github.com/kubernetes/kubernetes/issues/100551  k8s.io/kubernetes/test/integration/job.TestMain (ci-kubernetes-integration-master)

SIG Api-Machinery
#100787 https://github.com/kubernetes/kubernetes/issues/100787 [sig-api-machinery] AdmissionWebhook - should mutate custom resource with pruning [Conformance]

SIG
#21185 https://github.com/kubernetes/test-infra/issues/21185 ci-cluster-api-provider-gcp-make-conformance-v1alpha3-k8s-ci-artifacts failing to get kubeconfig

SIG Node
#100252 https://github.com/kubernetes/kubernetes/issues/100252  [sig-node] Container Runtime blackbox test when running a container with a new image should be able to pull from private registry with secret [NodeConformance] (ci-kubernetes-e2e-windows-gce-20h2)

SIG Storage
#97040 https://github.com/kubernetes/kubernetes/issues/97040 [Flaky Test] k8s.io/kubernetes/test/integration/storageversion.TestStorageVersionBootstrap 

SIG Scheduling
#103655 https://github.com/kubernetes/kubernetes/issues/103655  [sig-scheduling] Pod should avoid nodes that have avoidPod annotation (ci-kubernetes-e2e-gci-gce-serial)

SIG Network
#100132 https://github.com/kubernetes/kubernetes/issues/100132 [Flaky test] [sig-network] Services should be able to update service

Failures in Master-Blocking
	18 jobs total
	12 are passing
	5 are flaking
	1 are failing
	0 are stale


Failures in Master-Informing
	23 jobs total
	9 are passing
	12 are flaking
	2 are failing
	0 are stale


```
