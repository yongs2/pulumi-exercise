- [Kubernetes: Installation & Configuration](https://www.pulumi.com/registry/packages/kubernetes/installation-configuration/)

```sh
export KUBECONFIG="/workspace/.k8s-tb.config"
kubectl config view
kubectl get node -o wide
```

- [Pulumi & Kubernetes: Create project](https://www.pulumi.com/docs/clouds/kubernetes/get-started/create-project/)

```sh
# Pulumi & Kubernetes: Create project
export PULUMI_CONFIG_PASSPHRASE=
cd /workspace
mkdir -p 01.quickstart-go && cd 01.quickstart-go
pulumi new kubernetes-go
pulumi stack select dev
pulumi stack ls
```

- [Pulumi & Kubernetes: Review project](https://www.pulumi.com/docs/clouds/kubernetes/get-started/review-project/)

```sh
cat main.go
```

- [Pulumi & Kubernetes: Deploy stack](https://www.pulumi.com/docs/clouds/kubernetes/get-started/deploy-stack/)
  - go run main.go 를 실행하면 exit status 32 로 처리되므로 반드시 pulumi up 으로 실행해야 함
  - 예제와 다르게 github.com/pulumi/pulumi-kubernetes/sdk/v4 으로 생성되므로 일부 소스 수정 필요

```sh
pulumi up -y
pulumi stack output
```

- [Pulumi & Kubernetes: Modify program](https://www.pulumi.com/docs/clouds/kubernetes/get-started/modify-program/)
  - if No LoadBalancer, isMinikube set true

```sh
pulumi config set isMinikube true
```

- [Pulumi & Kubernetes: Deploy changes](https://www.pulumi.com/docs/clouds/kubernetes/get-started/deploy-changes/)

```sh
pulumi up -y
pulumi stack output
```

- [Pulumi & Kubernetes: Destroy stack](https://www.pulumi.com/docs/clouds/kubernetes/get-started/destroy-stack/)

```sh
pulumi destroy -y
```
