# [Docs / Clouds / Kubernetes / Get started / Next steps](https://www.pulumi.com/docs/clouds/kubernetes/get-started/next-steps/)

## [How-to Guides](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/)

### [Kubernetes: Create, Update, and Destroy](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/configmap-rollout/)

- Create a Project

```sh
cd /workspace
mkdir -p 10.k8s-configmap-rollout-go && cd 10.k8s-configmap-rollout-go
pulumi new kubernetes-go --name k8s-configmap-rollout-go -y --force
pulumi stack ls
```

- [refer example code](https://github.com/pulumi/examples/tree/master/kubernetes-ts-configmap-rollout)
- [api-doc of configmap](https://www.pulumi.com/registry/packages/kubernetes/api-docs/core/v1/configmap/)

- Deploying

```sh
# cluster IP
pulumi config set isMinikube true
pulumi up -y

# LoadBalancer IP
pulumi config set isMinikube false
pulumi up -y
```

- check configmap

```sh
kubectl get configmap $(pulumi stack output nginxConfigName) -o yaml
```

- Get the IP address once the chart is deployed and ready

```sh
# cluster IP
pulumi stack output frontendIP
kubectl port-forward svc/$(pulumi stack output svcName) 8765:80 --address='0.0.0.0'
curl -v http://localhost:8765
# redirect pulumi.com

# change and apply
sed -i "s/pulumi.github.io/google.com/g" default.conf
pulumi preview --diff
pulumi up -y

kubectl port-forward svc/$(pulumi stack output svcName) 8765:80 --address='0.0.0.0'
curl -v http://localhost:8765
# redirect google.com
```

- Clean up

```sh
sed -i "s/google.com/pulumi.github.io/g" default.conf
pulumi destroy -y

# Clean up procedure in exceptional cases
kubectl get deployment --no-headers | awk '{print $1}' | xargs kubectl delete deployment
kubectl get statefulset --no-headers | awk '{print $1}' | xargs kubectl delete statefulset
```
