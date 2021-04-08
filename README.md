# CI signal report

You can get the current overview for CI signal report by running

```
GITHUB_AUTH_TOKEN=xxx RELEASE_VERSION=xxx USER=xxx go run report.go
```
Where the RELEASE_VERSION can be like `1.21` and USER is your GitHub user.

It needs github token to be able to query the project board for CI signal. For some reason even though those boards are available for public view, the APIs require auth. See [this documentation](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line) to set up your access token.

## Prerequisites
- GoLang >=1.16

## Running
```
git clone git@github.com:alenkacz/ci-signal-report.git <folder>
cd <folder>
GITHUB_AUTH_TOKEN=xxx RELEASE_VERSION=xxx USER=xxx go run report.go
```

## Ratelimits
GitHub API has rate limits, to see how much you have used you can query like this (replace User with your GH user and Token with your Auth Token):
```
curl \
  -u USER:TOKEN -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/rate_limit
```


## Example output

```
Resolved
SIG cluster-lifecycle
#80434 https://api.github.com/repos/kubernetes/kubernetes/issues/80434  Errors bringing up kube-proxy in CI

In flight
SIG testing
#79662 https://api.github.com/repos/kubernetes/kubernetes/issues/79662   Nodes resize test failing in master-blocking

SIG cluster-lifecycle
#78907 https://api.github.com/repos/kubernetes/kubernetes/issues/78907 [Flaky Tests] task-06-upgrade is failing on master-informing

SIG scheduling
#74931 https://api.github.com/repos/kubernetes/kubernetes/issues/74931 Scheduler TestPreemptionRaces is flaky

SIG apps
#79740 https://api.github.com/repos/kubernetes/kubernetes/issues/79740  Test Deployment deployment should support rollback is failing on master informing

SIG cli
#79533 https://api.github.com/repos/kubernetes/kubernetes/issues/79533  Kubectl client Conformance test failing

New/Not Yet Started
SIG network
#80719 https://api.github.com/repos/kubernetes/kubernetes/issues/80719 [sig-network] Services should only allow access from service loadbalancer source ranges [Slow]

SIG storage
#80717 https://api.github.com/repos/kubernetes/kubernetes/issues/80717  [sig-storage] CSI Volumes [Driver: csi-hostpath] Snapshot Tests

SIG release
#80715 https://api.github.com/repos/kubernetes/kubernetes/issues/80715 [ Failing test ] build-packages-debs, build-packages-rpms

SIG network
#77538 https://api.github.com/repos/kubernetes/kubernetes/issues/77538 Some [sig-network] tests don't work on "private" clusters

SIG node
#74917 https://api.github.com/repos/kubernetes/kubernetes/issues/74917 [test failed] "regular resource usage tracking resource tracking for 100 pods per node" for containerd

SIG cluster-lifecycle
#78901 https://api.github.com/repos/kubernetes/kubernetes/issues/78901 [Flaky Tests] Flaky reboot tests

SIG cluster-lifecycle
#74893 https://api.github.com/repos/kubernetes/kubernetes/issues/74893 [Failing test] Upgrade [Feature:Upgrade] cluster upgrade should maintain a functioning cluster [Feature:ClusterUpgrade]

SIG auth
#75563 https://api.github.com/repos/kubernetes/kubernetes/issues/75563 Cleanup Advanced Audit testing

Failures in Master-Blocking
    14 jobs total
    9 are passing
    2 are flaking
    3 are failing
    0 are stale


Failures in Master-Informing
    14 jobs total
    7 are passing
    4 are flaking
    3 are failing
    0 are stale
 ```
