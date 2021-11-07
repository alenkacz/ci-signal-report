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
git clone git@github.com/alenkacz/ci-signal-report.git <folder>
cd <folder>
GITHUB_AUTH_TOKEN=xxx go run main.go
```

### Flags

- `-h` info about the flags
- `-short` shortens the report output (This reduces the report to `New/Not Yet Started` and `In Flight` issues on github.)
- `-emoji-off` report does not print emojis (see example output with emojis)
- `-v XXX` specify a k8s release version that should be added to the testgrid report. Where the XXX can be like `1.22`, the report statistics get extended for the chosen version. To specify multiple version use `-v "1.22, 1.21"`
- `-json` prints in json format

Example

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
$GITHUB_AUTH_TOKEN=XXX go run ./cmd/ci-reporter.go
GITHUB REPORT

#105965 volume metrics tests failure after removal of storage_operation_status_count [sig/storage]
- https://github.com/kubernetes/kubernetes/issues/105965
- Created 2021-10-28, Updated 2021-10-28, Comments: 3
- kind/failing-test
#105242 [Failing test][sig-storage] ci-kubernetes-e2e-gci-gce-serial [sig/storage]
- https://github.com/kubernetes/kubernetes/issues/105242
- 🔴Created 2021-09-24, ✨Updated 2021-11-05, Comments: 11
- priority/important-soon kind/failing-test milestone v1.23
#105675 HPA CPU e2e tests are failing [sig/autoscaling]
- https://github.com/kubernetes/kubernetes/issues/105675
- Created 2021-10-14, Updated 2021-10-14, Comments: 2
- kind/failing-test
#105580 [Failing Job] periodic-conformance-main-k8s-main [sig/cluster]
- https://github.com/kubernetes/kubernetes/issues/105580
- Created 2021-10-08, Updated 2021-11-01, Comments: 3
- priority/important-soon kind/failing-test milestone v1.23
#106139 Failure test: Volume metrics Ephemeral should create prometheus metrics for volume provisioning and attach/detach [sig/node]
- https://github.com/kubernetes/kubernetes/issues/106139
- ✨Created 2021-11-04, Updated 2021-11-04, Comments: 1
- kind/failing-test
#97783 Device manager for Windows passes when run on cluster that does not have a GPU but cuases cascading errors [sig/windows]
- https://github.com/kubernetes/kubernetes/issues/97783
- 🔴Created 2021-01-07, Updated 2021-10-14, Comments: 17
- kind/flake
#105677 HPA Custom metrics tests are failing [sig/autoscaling sig/testing]
- https://github.com/kubernetes/kubernetes/issues/105677
- Created 2021-10-14, Updated 2021-10-14, Comments: 2
- kind/failing-test
#99877 [flaky test] Test_Run_Positive_VolumeAttachMountUnmountDetach [sig/storage]
- https://github.com/kubernetes/kubernetes/issues/99877
- 🔴Created 2021-03-06, 🔴Updated 2021-10-05, Comments: 6
- kind/bug kind/flake
#89178 The Multi-AZ spreading test makes an incorrect assumption about which zones are candidates [sig/scheduling sig/testing]
- https://github.com/kubernetes/kubernetes/issues/89178
- 🔴Created 2020-03-17, Updated 2021-10-28, Comments: 18
- kind/flake
#105336 supported CNIs should have stable networking for Linux and Windows pods  [sig/windows]
- https://github.com/kubernetes/kubernetes/issues/105336
- 🔴Created 2021-09-29, 🔴Updated 2021-09-29, Comments: 1
- kind/flake
#104956 kubernetes-e2e-gce-stable1-latest-upgrade-master-parallel failed for oidc-discovery-test pod [sig/testing]
- https://github.com/kubernetes/kubernetes/issues/104956
- 🔴Created 2021-09-13, 🔴Updated 2021-09-14, Comments: 3
- kind/failing-test
#106008 [Failing test][sig-cloud-provider] gce-cos-master-serial [sig/cloud]
- https://github.com/kubernetes/kubernetes/issues/106008
- Created 2021-10-29, Updated 2021-11-01, Comments: 3
- priority/important-soon kind/failing-test
#98574 [Flaky Test] ci-kubernetes-e2e-aks-engine-azure-master-windows-containerd CNI failed to parse json [sig/windows]
- https://github.com/kubernetes/kubernetes/issues/98574
- 🔴Created 2021-01-29, Updated 2021-10-27, Comments: 13
- kind/flake
#106031 [Flaky Test] gce-ubuntu-master-containerd [sig/cloud]
- https://github.com/kubernetes/kubernetes/issues/106031
- Created 2021-10-31, Updated 2021-10-31, Comments: 3
- kind/flake
#104173 [Flaking test] DATA RACE in TestVolumeUnmountAndDetachControllerEnabled [sig/node]
- https://github.com/kubernetes/kubernetes/issues/104173
- 🔴Created 2021-08-05, 🔴Updated 2021-09-01, Comments: 8
- kind/bug kind/flake kind/failing-test


TESTGRID REPORT


🔥 Tests in Master-Blocking
- 18 jobs total
- 15 jobs passing
- 3 jobs flaky
- 0 jobs failing


FAILING & FLAKY JOBS:
FLAKY 🔵 verify-master
- https://testgrid.k8s.io/sig-release-master-blocking#verify-master
- 8 of 9 passed recently
FLAKY 🔵 gci-gce-ingress
- https://testgrid.k8s.io/sig-release-master-blocking#gci-gce-ingress
- 9 of 10 passed recently
FLAKY 🔵🔵🔵 gce-cos-master-default
- https://testgrid.k8s.io/sig-release-master-blocking#gce-cos-master-default
- 3 of 9 passed recently


💡 Tests in Master-Informing
- 23 jobs total
- 10 jobs passing
- 12 jobs flaky
- 1 jobs failing


FAILING & FLAKY JOBS:
FLAKY ✨ post-release-push-image-vulndash
- https://testgrid.k8s.io/sig-release-master-informing#post-release-push-image-vulndash
- 1 of 2 passed recently
FLAKY 🔵 kubeadm-kinder-latest
- https://testgrid.k8s.io/sig-release-master-informing#kubeadm-kinder-latest
- 9 of 10 passed recently
FLAKY ✨ post-kubernetes-push-image-etcd
- https://testgrid.k8s.io/sig-release-master-informing#post-kubernetes-push-image-etcd
- 2 of 3 passed recently
FLAKY 🔵 capg-conformance-v1alpha4-k8s-master
- https://testgrid.k8s.io/sig-release-master-informing#capg-conformance-v1alpha4-k8s-master
- 9 of 10 passed recently
FLAKY 🔵 ci-crio-cgroupv1-node-e2e-conformance
- https://testgrid.k8s.io/sig-release-master-informing#ci-crio-cgroupv1-node-e2e-conformance
- 8 of 9 passed recently
FAILING 🔴🔴🔴 gce-cos-master-serial
- https://testgrid.k8s.io/sig-release-master-informing#gce-cos-master-serial
- Sig's involved [sig-storage]
- Currently 6 test are failing
- 0 of 9 passed recently
FLAKY 🔵🔵 gce-cos-master-slow
- https://testgrid.k8s.io/sig-release-master-informing#gce-cos-master-slow
- 7 of 9 passed recently
FLAKY ✨ post-release-push-image-kube-cross
- https://testgrid.k8s.io/sig-release-master-informing#post-release-push-image-kube-cross
- 0 of 1 passed recently
FLAKY 🔵 aks-engine-windows-containerd-master
- https://testgrid.k8s.io/sig-release-master-informing#aks-engine-windows-containerd-master
- 9 of 9 passed recently
FLAKY 🔵 periodic-conformance-main-k8s-main
- https://testgrid.k8s.io/sig-release-master-informing#periodic-conformance-main-k8s-main
- 9 of 10 passed recently
FLAKY ✨ post-release-push-image-debian-base
- https://testgrid.k8s.io/sig-release-master-informing#post-release-push-image-debian-base
- 0 of 1 passed recently
FLAKY ✨ post-release-push-image-setcap
- https://testgrid.k8s.io/sig-release-master-informing#post-release-push-image-setcap
- 0 of 1 passed recently
FLAKY ✨ post-release-push-image-go-runner
- https://testgrid.k8s.io/sig-release-master-informing#post-release-push-image-go-runner
- 1 of 2 passed recently

```
