package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var videoFiles []VideoFile
var filesToMove []VideoFile

// VideoFile is just a holder for files to be moved.
type VideoFile struct {
	fullPath string
	name     string
	size     int64
	modified time.Time
}

func main() {

	dest := flag.String("dest", "/tmp/", "destination")
	src := flag.String("src", ".", "the (recursive) source directory of files to watch")

	flag.Parse()
	// if last char is not /, add it
	var lastChar = (*dest)[len(*dest)-1:]
	if !os.IsPathSeparator(lastChar[0]) {
		*dest = *dest + "/"
	}

	// get dir
	file, _ := filepath.Abs(*src)
	fmt.Println("Searching: " + file)

	// traverse dir
	err := filepath.Walk(file, walkFunc)

	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(5 * time.Second)
	loopOverFiles(videoFiles)

	// move files to new location
	err = moveFiles(filesToMove, *dest, *src)

	if err != nil {
		fmt.Println(err)
	}

}

// moveFiles moves a []VideoFile to dest
func moveFiles(vf []VideoFile, dest string, src string) (err error) {
	for _, element := range vf {

		folders := strings.SplitAfter(element.fullPath, filepath.Clean(src))
		fullDest, _ := filepath.Split(dest + folders[1])

		cleanFile := filepath.Clean(fullDest + element.name)

		line := fmt.Sprintf("Moving: %s to: %s", element.fullPath, cleanFile)
		fmt.Println(line)

		err = os.MkdirAll(filepath.Clean(fullDest), os.ModePerm)

		if runtime.GOOS == "linux" {
			err = os.Rename(element.fullPath, cleanFile)
		} else {
			err = CopyFile(element.fullPath, cleanFile)
		}
	}
	return
}

// loopOverFiles checks the file sizes and
// adds them to a []VideoFile if they are not changing.
func loopOverFiles(vf []VideoFile) {
	for _, element := range vf {
		fileString := element.fullPath

		updatedFile, _ := os.Stat(fileString)

		// We can move the file if the file has not changed in size
		// and if the last modified time is greater than 5 mins ago
		if (updatedFile.Size() == element.size) &&
			(time.Now().Add(-5 * time.Minute).After(element.modified)) {
			filesToMove = append(filesToMove, element)
		}
	}
}

// walkFunc recursively finds files and adds them to an array
func walkFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		fmt.Println("Found file: " + path)
		videoFiles = append(videoFiles, VideoFile{name: info.Name(), fullPath: path, size: info.Size(), modified: info.ModTime()})
	}
	return nil
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	err = copyFileContents(src, dst)
	err = os.Remove(src)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
