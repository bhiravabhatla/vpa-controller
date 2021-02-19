## VPA  Controller

Kubernetes operator to create and maintain verticalPodAutoscalers(VPA).


#### Run Locally

Make sure you have installed below tools:

* Kind ( brew install kind ) [https://kind.sigs.k8s.io/]
* Go (brew install go)
* Kubectl
* controller-gen (go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5)

To run the code in a local kind cluster -

```shell
make kind-setup
```

To clean up -

```shell
make kind-delete
```