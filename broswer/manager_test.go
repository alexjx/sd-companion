package broswer_test

import (
	"fmt"
	"testing"

	"github.com/alexjx/image-browser/broswer"
)

func TestBroswer(t *testing.T) {
	b := broswer.NewBroswer("/workspaces/image-browser/", []string{".go"})
	files, err := b.Files()
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Path)
	}
}
