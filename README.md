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
GITHUB_AUTH_TOKEN=XXX go run ./cmd/ci-reporter.go -short

----------
🤔 New/Not Yet StartedSIG Cloud-Provider
#100230 https://github.com/kubernetes/kubernetes/issues/100230 [Flaky Test] [sig-cloud-provider-gcp] Nodes [Disruptive] Resize [Slow] should be able to delete nodes


----------
🛫 In flightSIG Network
#93740 https://github.com/kubernetes/kubernetes/issues/93740 [Flaky Test][sig-network] Loadbalancing: L7 GCE [Slow] [Feature:Ingress] should conform to Ingress spec

SIG Api-Machinery
#100112 https://github.com/kubernetes/kubernetes/issues/100112 [flaky test] k8s.io/kubernetes/pkg/registry/core/endpoint/storage.TestWatch
#100760 https://github.com/kubernetes/kubernetes/issues/100760 [Flaking-test] Kubernetes e2e suite.[sig-api-machinery] AdmissionWebhook [Privileged:ClusterAdmin] listing validating webhooks should work [Conformance]

SIG Apps
#100314 https://github.com/kubernetes/kubernetes/issues/100314 [Flaky Test] Kubernetes e2e suite: [sig-apps] Deployment iterative rollouts should eventually progress
#98180 https://github.com/kubernetes/kubernetes/issues/98180 [Flaky Test] [sig-apps] Deployment should run the lifecycle of a Deployment

SIG Scalability
#97071 https://github.com/kubernetes/kubernetes/issues/97071 [Flaky test] [sig-storage] In-tree Volumes [Driver: gcepd] [Testpattern: Pre-provisioned PV (xfs)][Slow] volumes should store data
#103742 https://github.com/kubernetes/kubernetes/issues/103742 [Flaking Test] [sig-scalability] restarting konnectivity-agent (ci-kubernetes-e2e-gci-gce-scalability)

SIG Release
#1693 https://github.com/kubernetes/release/issues/1693 push-build.sh container image pushes should precede staging GCS artifacts and writing version markers

⛔ Tests in Master-Blocking
	18 jobs total
	12 are passing
	5 are flaking
	1 are failing
	0 are stale


💡 Tests in Master-Informing
	23 jobs total
	9 are passing
	12 are flaking
	2 are failing
	0 are stale

```
