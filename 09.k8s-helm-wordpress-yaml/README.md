# [Docs / Clouds / Kubernetes / Get started / Next steps](https://www.pulumi.com/docs/clouds/kubernetes/get-started/next-steps/)

## [How-to Guides](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/)

### [Kubernetes WordPress Helm Chart](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/wordpress-chart/)

- Create a Project

```sh
cd /workspace
mkdir -p 09.k8s-helm-wordpress-yaml && cd 09.k8s-helm-wordpress-yaml
pulumi new kubernetes-yaml --name k8s-helm-wordpress -y
pulumi stack ls
```

- [refer example code](https://github.com/pulumi/examples/blob/master/kubernetes-ts-helm-wordpress)
- Helm Chart resources are not supported in YAML, consider using the Helm Release resource instead
  - type: kubernetes:helm.sh/v3:Chart not SUPPORTED

- Deploying

```sh
pulumi up -y
```

- Get the IP address once the chart is deployed and ready

```sh
helm list
helm get values $(pulumi stack output name)
wordpressIP=$(kubectl get svc $(pulumi stack output name) -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
curl -v http://${wordpressIP}
curl -sL http://${wordpressIP} | grep "<title>"
```

- Clean up

```sh
pulumi destroy -y

# Clean up procedure in exceptional cases
kubectl get pvc --no-headers | awk '{print $1}' | xargs kubectl delete pvc
kubectl get deployment --no-headers | awk '{print $1}' | xargs kubectl delete deployment
kubectl get statefulset --no-headers | awk '{print $1}' | xargs kubectl delete statefulset
```
