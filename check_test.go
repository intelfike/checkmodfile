package checkmodfile

import (
	"fmt"
	"testing"
)

func TestCheck(t *testing.T) {
	f, err := RegistFile("check_test.go")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f.IsLatest())
}
