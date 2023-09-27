package osinsta

import (
	"log"

	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/compute"
	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/networking"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type OpenStackInstance struct {
	pulumi.ResourceState

	InstanceIP pulumi.StringOutput
	Instance   *compute.Instance
}

type OpenStackInstanceArgs struct {
	Name        pulumi.StringInput
	KeyPairName pulumi.StringInput
	PrivateKey  string
	FlavorId    pulumi.StringInput
	Image       pulumi.StringInput
	NetworkName string
	FixedIpV4   pulumi.StringInput
}

func NewOpenStackInstance(ctx *pulumi.Context, name string, args *OpenStackInstanceArgs, opts ...pulumi.ResourceOption) (pOpenStackInstance *OpenStackInstance, err error) {
	osinsta := &OpenStackInstance{}

	if err = ctx.RegisterComponentResource("k8s-go:component:OpenStackInstance", name, osinsta, opts...); err != nil {
		return nil, err
	}

	// [Get Network](https://www.pulumi.com/registry/packages/openstack/api-docs/networking/getnetwork/)
	// openstack network show management --format=json
	network, err := networking.LookupNetwork(ctx, &networking.LookupNetworkArgs{
		Name: &args.NetworkName,
	}, nil)
	if err != nil {
		return nil, err
	}
	ctx.Export("network.Id", pulumi.String(network.Id))
	ctx.Export("network.Name", pulumi.String(*network.Name))

	// [Get Subnet](https://www.pulumi.com/registry/packages/openstack/api-docs/networking/getsubnet/)
	// openstack subnet show management --format=json
	subnet, err := networking.LookupSubnet(ctx, &networking.LookupSubnetArgs{
		Name: &args.NetworkName,
	}, nil)
	if err != nil {
		return nil, err
	}
	ctx.Export("subnet.Id", pulumi.String(subnet.Id))

	// Create Port
	// openstack port create --network management --fixed-ip subnet=management,ip-address=192.168.5.48 basic
	newPort, err := networking.NewPort(ctx, name, &networking.PortArgs{
		Name:         args.Name.ToStringOutput(),
		NetworkId:    pulumi.String(network.Id),
		AdminStateUp: pulumi.Bool(true),
		FixedIps: networking.PortFixedIpArray{
			networking.PortFixedIpArgs{
				SubnetId:  pulumi.String(subnet.Id),
				IpAddress: args.FixedIpV4.ToStringOutput(),
			},
		},
	}, opts...)
	if err != nil {
		return nil, err
	}
	_ = newPort.Name.ApplyT(func(v interface{}) string {
		log.Printf("newPort.Name=[%v]", v.(string))
		return v.(string)
	}).(pulumi.StringOutput)
	ctx.Export("newPort.Name", newPort.Name)

	// Create an OpenStack resource (Compute Instance)
	// openstack server create --flavor "C1M2D20" --image ubuntu-20.04.5 --key-name packstack --port basic basic
	osinsta.Instance, err = compute.NewInstance(ctx, name, &compute.InstanceArgs{
		Name:      args.Name.ToStringOutput(),
		ImageName: args.Image.ToStringOutput(),
		FlavorId:  args.FlavorId.ToStringOutput(),
		KeyPair:   args.KeyPairName.ToStringOutput(),
		Networks: compute.InstanceNetworkArray{
			compute.InstanceNetworkArgs{
				Port: newPort.ID(), // ID() in CustomResourceState (pulimi/sdk/go/pulumi/resource.go)
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{newPort}))
	if err != nil {
		return nil, err
	}

	// Export the IP of the instance
	ctx.Export("InstanceIP", osinsta.Instance.AccessIpV4)
	osinsta.InstanceIP = osinsta.Instance.AccessIpV4.ApplyT(func(v interface{}) string {
		log.Printf("osinsta.Instance.AccessIpV4=[%v]", v.(string))
		return v.(string)
	}).(pulumi.StringOutput)
	return osinsta, nil
}
