package common

import (
	"fmt"
	"log"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ClusterVariable struct {
	UserName   string
	PrivateKey string
	Host       pulumi.StringInput
}

func RunCommand(ctx *pulumi.Context, name string, clusterVariable ClusterVariable, opts ...pulumi.ResourceOption) (pCommand *remote.Command, err error) {
	script := &remote.Command{}
	logPrefix := fmt.Sprintf("%v", name)
	localFileName := "k8s/common/install.sh"
	outputFileName := "/tmp/install_common.sh"

	log.Printf("[%v] Connection.Host=[%v],User=[%v]", logPrefix, clusterVariable.Host, clusterVariable.UserName)
	connection := remote.ConnectionArgs{
		Host:       clusterVariable.Host,
		PrivateKey: pulumi.String(clusterVariable.PrivateKey),
		User:       pulumi.String(clusterVariable.UserName),
	}

	// log.Printf("[%v] remote.NewCopyFile.opts=[%v]", logPrefix, opts)
	copyFile, err := remote.NewCopyFile(ctx, "install_common_"+name, &remote.CopyFileArgs{
		Connection: connection,
		LocalPath:  pulumi.String(localFileName),
		RemotePath: pulumi.String(outputFileName),
	}, opts...)
	if err != nil {
		log.Printf("[%v] remote.copyFile[%v].Err[%v]", logPrefix, copyFile, err)
		return nil, err
	}
	// log.Printf("[%v] remote.copyFile[%v].Done", logPrefix, copyFile)

	script, err = remote.NewCommand(ctx, "script_common_"+name, &remote.CommandArgs{
		Connection: connection,
		Create:     pulumi.String(fmt.Sprintf("chmod +x %v && %v", outputFileName, outputFileName)),
	}, pulumi.DependsOn([]pulumi.Resource{copyFile}))
	if err != nil {
		log.Printf("[%v] remote.NewCommand.Err[%v]", logPrefix, err)
		return nil, err
	}
	// log.Printf("[%v] remote.NewCommand[%v].Done", logPrefix, script)

	ctx.Export("output", script.Stdout)
	return script, nil
}
