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
	norminateRoot string
	root          string
	// skipLength is the length of the root path
	rootSkipLength int

	// extsFilter is the extensions of the files
	extsFilter map[string]struct{}

	// jpeg quality
	quality int
}

func NewBroswer(root string, exts []string, quality int) *Broswer {
	root = filepath.Clean(root)
	logrus.Infof("root path: %q", root)

	// create extensions filter
	extsFilter := make(map[string]struct{})
	for _, ext := range exts {
		ext = strings.ToLower(ext)
		extsFilter[ext] = struct{}{}
	}

	// eval the root if it is a link
	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		logrus.Errorf("eval root path %s error: %v", root, err)
		return nil
	}

	b := &Broswer{
		norminateRoot:  root,
		root:           realRoot,
		rootSkipLength: len(realRoot),
		extsFilter:     extsFilter,
		quality:        quality,
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
			// REVISIT: this is a bug, we should not have this,
			// if the root is a link to a directory
			if fpath == root {
				return nil
			}

			// skip synology eaDir
			if strings.Contains(fpath, "@eaDir") {
				return nil
			}

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
		return files[i].Path < files[j].Path
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

		if !info.IsDir() {
			return nil
		}
		// trim root prefix
		if fpath == root {
			return nil
		}

		// skip synology eaDir
		fpath = fpath[skipLen+1:]
		if strings.Contains(fpath, "@eaDir") {
			return nil
		}

		folders = append(folders, fpath)
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

func (b *Broswer) Delete(p string) error {
	// add .del suffix to the file
	filepath := path.Join(b.root, p)
	delPath := filepath + ".del"
	return os.Rename(filepath, delPath)
}
