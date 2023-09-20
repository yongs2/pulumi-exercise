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
mkdir -p 01.quickstart && cd 01.quickstart
ls -la ~/.pulumi/templates/
pulumi new kubernetes-yaml

# Pulumi & Kubernetes: Deploy stack
pulumi up --yes
pulumi stack output

# Pulumi & Kubernetes: Modify program
pulumi up --yes
pulumi stack output ip

# Pulumi & Kubernetes: Deploy changes
SVC_NAME=$(kubectl get svc -l app=nginx -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward service/${SVC_NAME} 8080:80
curl http://localhost:8080

# Pulumi & Kubernetes: Destroy stack
pulumi destroy -y
pulumi stack ls

# To delete the stack itself
pulumi stack rm dev -y
```
