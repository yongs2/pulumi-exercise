# [Docs / Clouds / Kubernetes / Get started / Next steps](https://www.pulumi.com/docs/clouds/kubernetes/get-started/next-steps/)

## [How-to Guides](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/)

### [Guestbook App with Redis and Nginx](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/guestbook/)

- Initailize

```sh
cd /workspace
mkdir -p 07.k8s-guestbook-go && cd 07.k8s-guestbook-go
pulumi new kubernetes-go --name k8s-guestbook -y
pulumi stack init dev
pulumi stack ls
```

- [download the example code](https://github.com/pulumi/examples/tree/master/kubernetes-go-guestbook)

- Deploying

```sh
pulumi config set useLoadBalancer false
pulumi config -j

pulumi up -y
```

- Viewing the Guestbook

```sh
# cluster IP, useLoadBalancer false
pulumi stack output frontendIP
kubectl port-forward svc/frontend 8765:80
curl -v http://localhost:8765

# LoadBalancer, useLoadBalancer true
curl -v $(pulumi stack output frontendIP)
```

- Modify frontend replicas

```sh
sed -i "s/frontendReplicas := 1/frontendReplicas := 3/g" main.go
pulumi up -y --skip-preview
```

- Clean up

```sh
pulumi destroy -y
```
