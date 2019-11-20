package builder

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"testing"
)

func TestCPUFreq(t *testing.T) {
	stat, err := cpu.Info()
	if err != nil {
		println(err)
	}
	fmt.Println(stat)
}

func TestLoadAvg(t *testing.T) {
	l, err := load.Avg()
	if err != nil {
		println(err)
	}
	print(float32(l.Load5))
}

func TestMem(t *testing.T) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(memInfo)
	fmt.Println(memInfo.Total / 1073741824)
}
