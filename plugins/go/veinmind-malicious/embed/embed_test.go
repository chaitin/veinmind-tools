package embed

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(ReadFile("scripts/.env"))
}
