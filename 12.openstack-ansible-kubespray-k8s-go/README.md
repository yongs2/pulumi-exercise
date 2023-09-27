# [openstack](https://www.pulumi.com/registry/packages/openstack/)

## [Example](https://www.pulumi.com/registry/packages/openstack/#example)

## [Package Command](https://www.pulumi.com/registry/packages/command/)

- Create a Project

```sh
cd /workspace
mkdir -p 12.openstack-ansible-kubespray-k8s-go && cd 12.openstack-ansible-kubespray-k8s-go
pulumi new openstack-go --name openstack-ansible-kubespray-k8s-go -y --force
pulumi stack ls
```

- Deploying

```sh
pulumi stack init dev
gofmt -w -s *.go; pulumi up -y
    
ssh ubuntu@$(pulumi stack output AnsibleIP)
ssh ubuntu@$(pulumi stack output k8s-master-01)
```

- Clean up

```sh
pulumi destroy -y

# Check exception case
openstack keypair show basic --format=json
openstack flavor show C1M2D20 --format=json
openstack port show basic --format=json
openstack server show basic --format=json

# Exception case
openstack server delete basic
openstack port delete basic
openstack flavor delete C1M2D20
openstack keypair delete basic
```
