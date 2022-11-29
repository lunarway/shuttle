package copy

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/exp/slices"
)

// File copies a single file from src to dst
func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	srcfd, err = os.Open(src)
	if err != nil {
		return err
	}
	defer srcfd.Close()

	dstfd, err = os.Create(dst)
	if err != nil {
		return err
	}
	defer dstfd.Close()

	_, err = io.Copy(dstfd, srcfd)
	if err != nil {
		return err
	}

	srcinfo, err = os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcinfo.Mode())
}

// Dir copies a whole directory recursively
func Dir(src string, dst string, ignorelist []string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	srcinfo, err = os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcinfo.Mode())
	if err != nil {
		return err
	}

	fds, err = ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		shouldIgnore := slices.Contains(ignorelist, fd.Name())
		if shouldIgnore {
			continue
		}

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp, ignorelist); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
