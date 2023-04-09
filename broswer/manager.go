package broswer

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type Broswer struct {
	// root is the root directory of the Broswer
	root string
	// skipLength is the length of the root path
	rootSkipLength int

	// extsFilter is the extensions of the files
	extsFilter map[string]struct{}

	// trash is the trash directory
	trash           string
	trashSkipLength int

	// jpeg quality
	quality int
}

func NewBroswer(root, trash string, exts []string, quality int) *Broswer {
	root = filepath.Clean(root)
	logrus.Infof("root path: %q", root)

	// ensure trash directory exists
	trash = filepath.Clean(trash)
	if err := os.MkdirAll(trash, 0755); err != nil {
		logrus.Errorf("create trash directory %s error: %v", trash, err)
	}
	logrus.Infof("trash path: %q", trash)

	// create extensions filter
	extsFilter := make(map[string]struct{})
	for _, ext := range exts {
		ext = strings.ToLower(ext)
		extsFilter[ext] = struct{}{}
	}

	b := &Broswer{
		root:            root,
		rootSkipLength:  len(root),
		extsFilter:      extsFilter,
		trash:           trash,
		trashSkipLength: len(trash),
		quality:         quality,
	}

	return b
}

func (b *Broswer) GetRoot() string {
	return b.root
}

func (b *Broswer) files(root, folder string, skipLen int) ([]*File, error) {
	files := []*File{}

	targetDir := root
	if folder != "" {
		targetDir = filepath.Join(root, folder)
		targetDir = filepath.Clean(targetDir)
	}

	err := filepath.Walk(targetDir, func(fpath string, info os.FileInfo, err error) error {
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

		if !info.IsDir() {
			// normalize the path relative to the root
			relativePath := fpath[skipLen+1:]

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

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModifiedAt.Before(files[j].ModifiedAt)
	})

	return files, nil
}

func (b *Broswer) folders(root string, skipLen int) ([]string, error) {
	var folders []string

	err := filepath.Walk(root, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Errorf("notify walk path %s error: %v", fpath, err)
			return err
		}

		if info.IsDir() {
			// trim root prefix
			if fpath == root {
				return nil
			}
			fpath = fpath[skipLen+1:]
			folders = append(folders, fpath)
		}
		return nil
	})
	if err != nil {
		logrus.Errorf("walk path %s error: %v", b.root, err)
		return nil, err
	}

	return folders, nil
}

func (b *Broswer) Files(folder string) ([]*File, error) {
	return b.files(b.root, folder, b.rootSkipLength)
}

func (b *Broswer) Folders() ([]string, error) {
	return b.folders(b.root, b.rootSkipLength)
}

func (b *Broswer) TrashFiles(d string) ([]*File, error) {
	return b.files(b.trash, d, b.trashSkipLength)
}

func (b *Broswer) TrashFolders() ([]string, error) {
	return b.folders(b.trash, b.trashSkipLength)
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

func (b *Broswer) Restore(p string) error {
	// normalize the path relative to the root
	filepath := path.Join(b.root, p)
	trashPath := path.Join(b.trash, p)
	return os.Rename(trashPath, filepath)
}
