package scheduler

import (
	"context"
	"errors"
	pb "github.com/Mr-Linus/Pump2/rpc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)
func CheckTask(c *gin.Context) (pb.BuildInfo, error){
	if c.PostForm("name") == ""{
		err := errors.New("The Build Image Name is Null. ")
		LogErrPrint(err)
		return pb.BuildInfo{},err
	}
	if c.PostForm("gpu") != "true" &&  c.PostForm("gpu") != "false"{
		err := errors.New("The Build Image GPU is not set. ")
		LogErrPrint(err)
		return pb.BuildInfo{},err
	}
	if c.PostForm("tf") != "true" &&  c.PostForm("tf") != "false"{
		err := errors.New("The Build Image TensorFlow Framework is not set. ")
		LogErrPrint(err)
		return pb.BuildInfo{},err
	}
	if c.PostForm("torch") != "true" &&  c.PostForm("torch") != "false"{
		err := errors.New("The Build Image PyTorch Framework is not set. ")
		LogErrPrint(err)
		return pb.BuildInfo{},err
	}
	if c.PostForm("test") != "true" &&  c.PostForm("test") != "false"{
		err := errors.New("The Build Image Test mode is not set. ")
		LogErrPrint(err)
		return pb.BuildInfo{},err
	}
	return pb.BuildInfo{
		Name:                 c.PostForm("name"),
		Gpu:                  ChangeToBool(c.PostForm("gpu")),
		Tf:                   ChangeToBool(c.PostForm("tf")),
		Torch:                ChangeToBool(c.PostForm("torch")),
		TfVersion:            c.PostForm("tfver"),
		TorchVersion:         c.PostForm("torchver"),
		Dependence:           c.PostForm("dep"),
		UseToTest:            ChangeToBool(c.PostForm("test")),
	} ,nil
}

func DoTask(address string, task pb.BuildInfo) (bool,error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		LogErrPrint(err)
		return false,err
	}
	defer conn.Close()
	c := pb.NewPump2Client(conn)
	// Contact the server and print out its response.
	ctx := context.Background()
	r, err := c.BuildImages(ctx,&task)
	if err != nil {
		LogErrPrint(err)
		return r.BuildStats,err
	}
	return r.BuildStats,nil
}

func ChangeToBool(key string) bool{
	if key == "true" || key == "True"{
		return true
	}
	return false
}