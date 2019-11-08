package server

import (
	"context"
	"fmt"
	"github.com/Mr-Linus/Pump2/rpc"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"testing"
)

func TestP2Server_BuildImages(t *testing.T) {
	p := rpc.BuildInfo{
		Name: "testbuild:test",
		Tf:   true, TfVersion: "1.14.0",
		Torch:        true,
		TorchVersion: "",
		Gpu:          true,
		UseToTest:    true}
	s := P2Server{}
	_, err := s.BuildImages(context.Background(), &p)
	if err != nil {
		print(err)
	}
}

func TestCPUFreq(t *testing.T) {
	stat, err := cpu.Info()
	if err != nil {
		println(err)
	}
	fmt.Println(stat)
}

func TestMem(t *testing.T) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(memInfo)
	fmt.Println(1024 * 1024 * 1024)
}
