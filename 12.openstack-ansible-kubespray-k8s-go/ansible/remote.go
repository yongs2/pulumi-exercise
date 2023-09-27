package ansible

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ClusterVariable struct {
	UserName         string
	PrivateKey       string
	Host             pulumi.StringInput
	HttpProxy        string
	MasterIps        string
	WorkerIps        string
	KubeControlHosts int
	MetallbIpRange   string
	KubesprayVersion string
}

func RunCommand(ctx *pulumi.Context, name string, clusterVariable ClusterVariable, opts ...pulumi.ResourceOption) (pCommand *remote.Command, err error) {
	script := &remote.Command{}
	logPrefix := fmt.Sprintf("%v", name)
	localFileName := "ansible/installation.sh.template"
	outputFileName := "/tmp/installation.sh"

	if err = createFileUsingTemplate(localFileName, outputFileName, clusterVariable); err != nil {
		log.Printf("[%v] createFileUsingTemplate.Err[%v]", logPrefix, err)
		return nil, err
	}

	log.Printf("[%v] Connection.Host=[%v],User=[%v]", logPrefix, clusterVariable.Host, clusterVariable.UserName)
	connection := remote.ConnectionArgs{
		Host:       clusterVariable.Host,
		PrivateKey: pulumi.String(clusterVariable.PrivateKey),
		User:       pulumi.String(clusterVariable.UserName),
	}

	// log.Printf("[%v] remote.NewCopyFile.opts=[%v]", logPrefix, opts)
	copyFile, err := remote.NewCopyFile(ctx, "installation.sh_"+name, &remote.CopyFileArgs{
		Connection: connection,
		LocalPath:  pulumi.String(outputFileName),
		RemotePath: pulumi.String(outputFileName),
	}, opts...)
	if err != nil {
		log.Printf("[%v] remote.copyFile[%v].Err[%v]", logPrefix, copyFile, err)
		return nil, err
	}
	// log.Printf("[%v] remote.copyFile[%v].Done", logPrefix, copyFile)

	script, err = remote.NewCommand(ctx, "script_"+name, &remote.CommandArgs{
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

func createFileUsingTemplate(templateFileName string, outputFileName string, data interface{}) (err error) {
	// create template
	var t *template.Template
	if t, err = template.ParseFiles(templateFileName); err != nil {
		log.Printf("ParseFiles.Err[%v]", err)
		return err
	}

	// create output file
	var f *os.File
	if f, err = os.Create(outputFileName); err != nil {
		log.Printf("Create.Err[%v]", err)
		return err
	}
	defer f.Close()

	// excute
	if err = t.Execute(f, data); err != nil {
		log.Printf("Execute.Err[%v]", err)
		return err
	}

	return nil
}
