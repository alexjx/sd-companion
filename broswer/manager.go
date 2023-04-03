package broswer

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type Broswer struct {
	// root is the root directory of the Broswer
	root string
	// skipLength is the length of the root path
	skipLength int

	// extsFilter is the extensions of the files
	extsFilter map[string]struct{}

	// trash is the trash directory
	trash string

	// jpeg quality
	quality int
}

func NewBroswer(root string, exts []string, quality int) *Broswer {
	root = filepath.Clean(root)
	logrus.Infof("root path: %q", root)

	// ensure trash directory exists
	trash := filepath.Join(root, ".trash")
	if err := os.MkdirAll(trash, 0755); err != nil {
		logrus.Errorf("create trash directory %s error: %v", trash, err)
	}
	trash = filepath.Clean(trash)
	logrus.Infof("trash path: %q", trash)

	// create extensions filter
	extsFilter := make(map[string]struct{})
	for _, ext := range exts {
		ext = strings.ToLower(ext)
		extsFilter[ext] = struct{}{}
	}

	b := &Broswer{
		root:       root,
		skipLength: len(root),
		extsFilter: extsFilter,
		trash:      trash,
		quality:    quality,
	}

	return b
}

func (b *Broswer) GetRoot() string {
	return b.root
}

func (b *Broswer) Files() ([]*File, error) {
	files := []*File{}

	err := filepath.Walk(b.root, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Errorf("walk path %s error: %v", fpath, err)

			// ignore permission denied error
			if os.IsPermission(err) {
				return nil
			}
			return err
		}

		// ignore trash directory
		fpath = filepath.Clean(fpath)
		// logrus.Infof("filepath: %q, trash: %q", fpath, b.trash)
		if strings.HasPrefix(fpath, b.trash) {
			return nil
		}

		if !info.IsDir() {
			// normalize the path relative to the root
			relativePath := fpath[b.skipLength+1:]

			// filter the file by extension
			ext := path.Ext(relativePath)
			ext = strings.ToLower(ext[1:]) // remove the dot
			if _, ok := b.extsFilter[ext]; !ok {
				return nil
			}

			files = append(files, &File{
				Path:       filepath.ToSlash(relativePath),
				Size:       info.Size(),
				ModifiedAt: info.ModTime(),
			})
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("walk path %s error: %v", b.root, err)
		return nil, err
	}

	return files, nil
}

func (b *Broswer) Delete(p string) error {
	// normalize the path relative to the root
	filepath := path.Join(b.root, p)
	trashPath := path.Join(b.trash, p)
	if err := os.MkdirAll(path.Dir(trashPath), 0755); err != nil {
		logrus.Errorf("create trash directory %s error: %v", path.Dir(trashPath), err)
		return err
	}
	return os.Rename(filepath, trashPath)
}
