package bootstrap

import (
	"fmt"
	"testing"
)

func TestGetStrapList(t *testing.T) {
	InitBootStrap("./bootstrap.list")
	fmt.Println(B.bootlist)
	wanted := []string{"47.52.172.254", "106.14.187.240"}
	if !equal(B.bootlist, wanted) {
		t.Error("Not wanted bootlist")
	}
}

func equal(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
