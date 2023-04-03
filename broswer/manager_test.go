package broswer_test

import (
	"fmt"
	"testing"

	"github.com/alexjx/sd-companion/broswer"
)

func TestBroswer(t *testing.T) {
	b := broswer.NewBroswer("/workspaces/image-browser/", []string{".go"}, 50)
	files, err := b.Files()
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Path)
	}
}
