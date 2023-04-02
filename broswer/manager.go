package broswer

import (
	"os"
	"path"
	"path/filepath"
)

type Broswer struct {
	// root is the root directory of the Broswer
	root string
	// skipLength is the length of the root path
	skipLength int

	// extsFilter is the extensions of the files
	extsFilter map[string]struct{}
}

func NewBroswer(root string, exts []string) *Broswer {
	root = path.Clean(root)

	// create extensions filter
	extsFilter := make(map[string]struct{})
	for _, ext := range exts {
		extsFilter[ext] = struct{}{}
	}

	b := &Broswer{
		root:       root,
		skipLength: len(root),
		extsFilter: extsFilter,
	}

	return b
}

func (b *Broswer) GetRoot() string {
	return b.root
}

func (b *Broswer) Files() ([]*File, error) {
	var files []*File
	err := filepath.Walk(b.root, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			// ignore permission denied error
			if os.IsPermission(err) {
				return nil
			}
			return err
		}

		if !info.IsDir() {
			// normalize the path relative to the root
			relativePath := filepath[b.skipLength:]

			// filter the file by extension
			ext := path.Ext(relativePath)
			if _, ok := b.extsFilter[ext]; !ok {
				return nil
			}

			files = append(files, &File{
				Path:       relativePath,
				Size:       info.Size(),
				ModifiedAt: info.ModTime(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (b *Broswer) Delete(p string) error {
	// normalize the path relative to the root
	filepath := path.Join(b.root, p)
	return os.Remove(filepath)
}

func (b *Broswer) Content(p string) ([]byte, error) {
	return nil, nil
}

func (b *Broswer) Metadata(p string) (map[string]string, error) {
	return nil, nil
}
