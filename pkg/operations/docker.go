package operations

import (
	"archive/tar"
	"context"
	"fmt"
	"github.com/Mr-Linus/Pump2/pkg/pump2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
)


func ImageBuild(imagename string, args string) (types.ImageBuildResponse,error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}
	opt := types.ImageBuildOptions{
		Tags:[]string{imagename},
		BuildArgs: map[string]*string{"PACKAGE":&args},
	}
	btx, err := os.Open(os.Getenv("HOME")+"/Archive.tar")
	if err != nil {
		fmt.Println(err)
		return types.ImageBuildResponse{}, err
	}
	out, err := cli.ImageBuild(ctx, btx, opt)
	if err != nil {
		fmt.Println(err)
		return types.ImageBuildResponse{}, err
	}
	return out, nil
}

func ConfigDockerfile(in *pump2.BuildInfo) (dockerfile string, err error) {
	dockerfile = "https://raw.githubusercontent.com/NJUPT-ISL/Pump/master/template/"
	if in.Gpu {
		dockerfile += "gpu.Dockerfile"
	} else {
		dockerfile += "cpu.Dockerfile"
	}
	dockerfile, err = GetOperation(dockerfile)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if in.UseToTest{
		dockerfile += `
EXPOSE 22
CMD    ["/usr/sbin/sshd", "-D"]
`
	}
	f, err := os.Create(os.Getenv("HOME")+"/Dockerfile")
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer f.Close()
	_, err = f.WriteString(dockerfile)
	if err != nil {
		log.Println(err)
		return "", err
	}
	err = TarDockerfile()
	if err != nil {
		log.Println(err)
		return "", err
	}
	return os.Getenv("HOME"),nil
}

func TarDockerfile() error{
	f,err := os.Create(os.Getenv("HOME")+"/Archive.tar")
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()
	tw := tar.NewWriter(f)
	defer tw.Close()
	fInfo, err := os.Stat(os.Getenv("HOME")+"/Dockerfile")
	if err != nil {
		log.Println(err)
		return err
	}
	header, err := tar.FileInfoHeader(fInfo, "")
	err = tw.WriteHeader(header)
	fdocker, err := os.Open(os.Getenv("HOME")+"/Dockerfile")
	if err != nil {
		return err
	}
	_, err = io.Copy(tw, fdocker)
	if err != nil {
		return err
	}
	err = os.Remove(os.Getenv("HOME")+"/Dockerfile")
	if err != nil {
		return err
	}
	return nil
}

func ConfigBuildArgs(in *pump2.BuildInfo) (args string){
	args = ""
	if in.GetTf() {
		args += "tensorflow"
		if in.GetGpu() {
			args += "-gpu"
		}
		if in.TfVersion != ""{
			args += "=="+in.TfVersion+" "
		}
	}
	if in.GetTorch() {
		args += "tensorflow"
		if in.GetTfVersion() != ""{
			args += "=="+in.TfVersion+" "
		}
	}
	args += in.GetDependence()
	return args
}