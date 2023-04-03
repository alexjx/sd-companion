package broswer

import (
	"os"
	"path"
	"path/filepath"
)

type ImageFile struct {
	*os.File

	// Name is the name of the file
	Name string
}

func (i *ImageFile) Ext() string {
	return path.Ext(i.Name)
}

func (b *Broswer) Open(p string) (*ImageFile, error) {
	p = filepath.FromSlash(p)
	fPath := path.Join(b.root, p)
	f, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}

	basename := path.Base(fPath)

	return &ImageFile{
		File: f,
		Name: basename,
	}, nil
}
