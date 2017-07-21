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

func TestGetBytes(t *testing.T) {
	f, err := RegistFile("check_test.go")
	if err != nil {
		fmt.Println(err)
	}
	for n := 0; n < 100; n++ {
		text, err := f.GetBytes()
		if err != nil {
			continue
		}
		fmt.Printf("%s", text[:10])
	}

}
