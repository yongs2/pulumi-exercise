# [openstack](https://www.pulumi.com/registry/packages/openstack/)

## [Example](https://www.pulumi.com/registry/packages/openstack/#example)

- Create a Project

```sh
cd /workspace
mkdir -p 11.openstack-compute-instance-go && cd 11.openstack-compute-instance-go
pulumi new openstack-go --name openstack-compute-instance-go -y --force
pulumi stack ls
```

- Deploying

```sh
pulumi up -y

ssh ubuntu@192.168.5.48
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
