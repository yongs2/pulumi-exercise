# [Docs / Clouds / Kubernetes / Get started / Next steps](https://www.pulumi.com/docs/clouds/kubernetes/get-started/next-steps/)

## [How-to Guides](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/)

### [kubernetes.helm.sh/v3.Release](https://www.pulumi.com/registry/packages/kubernetes/api-docs/helm/v3/release/)

### [kubernetes.helm.sh/v3.Chart](https://www.pulumi.com/registry/packages/kubernetes/api-docs/helm/v3/chart/)

```sh
cd /workspace
mkdir -p 05.learn-pulumi-helm && cd 05.learn-pulumi-helm
ls -la ~/.pulumi/templates/
pulumi new helm-kubernetes-yaml --name learn-pulumi-helm -y

pulumi up --yes
pulumi stack output

helm list -A | grep nginx

pulumi destroy -y
```
