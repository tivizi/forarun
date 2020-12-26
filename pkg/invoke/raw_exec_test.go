package invoke

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestExec(t *testing.T) {
	exec := RawExec{Stderr: os.Stderr}
	bs, err := exec.ExecPlugin(context.Background(), "/home/ujued/main", []byte("{}"), []string{"TE=1TE2"})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bs))
}
