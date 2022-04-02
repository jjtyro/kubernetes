package cache_version

import (
	"fmt"
	"strings"
	"testing"
)

func TestCheckpoint(t *testing.T) {
	fname := "/tmp/test.file"

	b, err := PathExists(fname)
	fmt.Println("File exist:", b)
	fmt.Println("")

	fmt.Println("Write File Done")
	content := "#^ 01234567:v0.1.2.3$#"
	WriteFile(fname, content)

	b, err = PathExists(fname)
	fmt.Println("File exist:", b)
	fmt.Println("")

	str, err := ReadFile(fname)
	if err == nil {
		fmt.Println("Read content:", str)
		fmt.Println("")
		str = strings.Replace(str, "#", "", -1)
		fmt.Println("Trans content:", str)
		fmt.Println("")
	} else {
		fmt.Println(err)
	}
}