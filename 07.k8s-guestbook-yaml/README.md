# [Docs / Clouds / Kubernetes / Get started / Next steps](https://www.pulumi.com/docs/clouds/kubernetes/get-started/next-steps/)

## [How-to Guides](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/)

### [Guestbook App with Redis and Nginx](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/guestbook/)

- Initialization procedure for implementing pulumi in YAML

```sh
cd /workspace
mkdir -p 07.k8s-guestbook-yaml && cd 07.k8s-guestbook-yaml
pulumi new kubernetes-yaml --name k8s-guestbook-yaml -y
pulumi stack ls
```

- [See the go example source](https://github.com/pulumi/examples/tree/master/kubernetes-go-guestbook)

- Deploying

```sh
pulumi up -y
```

- Viewing the Guestbook

```sh
# cluster IP
pulumi stack output frontendIP
kubectl port-forward svc/frontend 8765:80
curl -v http://localhost:8765
```

- Clean up

```sh
pulumi destroy -y
```

- Appendix
  - [Pulumi YAML does not have any for loop or conditional logic](https://github.com/pulumi/pulumi/discussions/12560)
  - Convert to Other Pulumi Languages

  ```sh
  pulumi convert --language go --out ./k8s-guestbook-go
  ```
