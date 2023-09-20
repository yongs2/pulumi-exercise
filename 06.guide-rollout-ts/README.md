# [Docs / Clouds / Kubernetes / Get started / Next steps](https://www.pulumi.com/docs/clouds/kubernetes/get-started/next-steps/)

## [How-to Guides](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/)

### [Gated Kubernetes Deployments with Prometheus](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/p8s-rollout/)

- Initailize

```sh
cd /workspace
mkdir -p 06.guide-rollout-ts && cd 06.guide-rollout-ts
pulumi new typescript --name guide-rollout-ts -y
```

- [download the example code](https://github.com/pulumi/examples/blob/master/kubernetes-ts-staged-rollout-with-prometheus)

```sh
# download the example code
npm install

helm repo add prometheus https://prometheus-community.github.io/helm-charts
helm repo update

pulumi up -y
```

- clean up

```sh
pulumi destroy -y

# Clean up procedure in exceptional cases
kubectl get pvc --no-headers | awk '{print $1}' | xargs kubectl delete pvc
kubectl delete deployment.apps/p8s-prometheus-server
kubectl delete statefulset.apps/p8s-alertmanager
kubectl get svc | grep p8s | awk '{print $1}' | xargs kubectl delete svc
```
