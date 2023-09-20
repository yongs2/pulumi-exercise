# [Docs / Clouds / Kubernetes / Get started / Next steps](https://www.pulumi.com/docs/clouds/kubernetes/get-started/next-steps/)

## [How-to Guides](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/)

### [Kubernetes WordPress Helm Chart](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/wordpress-chart/)

- Create a Project

```sh
cd /workspace
mkdir -p 09.k8s-helm-wordpress-go && cd 09.k8s-helm-wordpress-go
pulumi new kubernetes-go --name k8s-helm-wordpress -y
pulumi stack select dev
pulumi stack ls
```

- [refer example code](https://github.com/pulumi/examples/blob/master/kubernetes-ts-helm-wordpress)

- Deploying

```sh
gofmt -w -s *.go; pulumi up -y
```

- Get the IP address once the chart is deployed and ready

```sh
curl -v http://$(pulumi stack output frontendIP)
curl -sL http://$(pulumi stack output frontendIP) | grep "<title>"
```

- Clean up

```sh
pulumi destroy -y

# Clean up procedure in exceptional cases
kubectl get pvc --no-headers | awk '{print $1}' | xargs kubectl delete pvc
kubectl get deployment --no-headers | awk '{print $1}' | xargs kubectl delete deployment
kubectl get statefulset --no-headers | awk '{print $1}' | xargs kubectl delete statefulset
```
