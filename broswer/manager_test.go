package broswer_test

import (
	"testing"

	"github.com/alexjx/sd-companion/broswer"
)

func TestNewBroswer(t *testing.T) {
	b := broswer.NewBroswer(
		"/workspaces/sd-companion/tmp",
		"/workspaces/sd-companion/tmp_trash",
		[]string{".jpg", ".png"},
		80,
	)

	if b.GetRoot() != "/workspaces/sd-companion/tmp" {
		t.Errorf("root path is not correct")
	}

	folders, err := b.Folders()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(folders) != 1 {
		t.Errorf("folders count is not correct")
	}

    
}
