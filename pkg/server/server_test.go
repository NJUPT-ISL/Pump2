package server

import (
	"context"
	"github.com/Mr-Linus/Pump2/pkg/pump2"
	"testing"
)


func TestP2Server_BuildImages(t *testing.T) {
	p := pump2.BuildInfo{
		Name:"testbuild:test",
		Tf:true,TfVersion:"1.14.0",
		Torch:true,
		TorchVersion:"",
		Gpu:true,
		UseToTest:true}
	s := P2Server{}
	_,err := s.BuildImages(context.Background(),&p)
	if err != nil{
		print(err)
	}
}