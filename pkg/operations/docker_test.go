package operations

import "testing"

func TestImageBuild(t *testing.T) {
	s := "asdasd"
	_, err := ImageBuild("sss",&s)
	if err != nil {
		print(err)
	}
}
