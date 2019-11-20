package yaml

import (
	"os"
	"testing"
)

func TestReadNodeYaml(t *testing.T) {
	println(os.Getenv("PWD"))
	conf,err := ReadNodeYaml("../../template/pump2-scheduler.yaml")
	if err != nil{
		println(err)
	}
	println(conf.Nodes[0].IP)
}