# [Docs / Clouds / Kubernetes / Get started / Next steps](https://www.pulumi.com/docs/clouds/kubernetes/get-started/next-steps/)

## [How-to Guides](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/)

### [Kubernetes Stateless App Deployment](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/stateless-app/)

- Create a Project

```sh
cd /workspace
mkdir -p 08.k8s-nginx-yaml && cd 08.k8s-nginx-yaml
pulumi new kubernetes-yaml --name k8s-nginx -y
pulumi stack ls
```

- Config a Project

```sh
pulumi config set replicas 1
pulumi config
pulumi config -j

pulumi up -y
```

- check deployment

```sh
kubectl describe deployment $(pulumi stack output nginx)
kubectl get pods -l app=nginx
```

- Update nginx version

```sh
sed -i "s/nginx:1.7.9/nginx:1.8/g" Pulumi.yaml
pulumi up -y
kubectl get pods -l app=nginx
```

- Scale our application by increasing the replica count

```sh
pulumi config set replicas 3
pulumi up -y
kubectl get pods -l app=nginx
```

- Clean up

```sh
pulumi destroy -y
kubectl get pods -l app=nginx
```

- Convert to golang

```sh
pulumi convert --language go --out ../08.k8s-nginx-go
```
