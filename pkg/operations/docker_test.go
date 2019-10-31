package operations

import "testing"

func TestImageBuild(t *testing.T) {
	_, err := ImageBuild("sss", "asdasd")
	if err != nil {
		print(err)
	}
}
