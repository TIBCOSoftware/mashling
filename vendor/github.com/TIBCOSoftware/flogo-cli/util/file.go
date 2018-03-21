package fgutil

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const FileURIPrefix = "file://"

type PathInfo struct {
	IsLocal  bool
	IsURL    bool
	IsFile   bool
	FileURL  *url.URL
	FilePath string
	FileName string
}

func GetPathInfo(pathStr string) (*PathInfo, error) {

	fileURL, err := url.Parse(pathStr)

	pi := &PathInfo{}

	if err != nil {
		return nil, err
	}

	if len(fileURL.Scheme) > 0 {
		pi.IsURL = true
		pi.FileURL = fileURL

		pi.FilePath, pi.IsLocal = URLToFilePath(fileURL)

		if pi.IsLocal {
			fileInfo, err := os.Stat(pi.FilePath)

			if err != nil {
				return nil, err
			}

			pi.IsFile = !fileInfo.IsDir()
		}
	} else {
		pi.FilePath = pathStr

		fileInfo, err := os.Stat(pi.FilePath)

		if err == nil {
			//path is a local file or directory
			pi.IsLocal = true
			pi.IsFile = !fileInfo.IsDir()
		}
	}

	if pi.IsFile {
		idx := strings.LastIndex(pathStr, "/")
		pi.FileName = pathStr[idx+1:]
	}

	return pi, nil
}

// ToFilePath convert fileURL to file path
func URLToFilePath(fileURL *url.URL) (string, bool) {

	if fileURL.Scheme == "file" {

		filePath := fileURL.Path

		if runtime.GOOS == "windows" {
			if strings.HasPrefix(filePath, "/") {
				filePath = filePath[1:]
			}
			filePath = filepath.FromSlash(filePath)
		}

		filePath = strings.Replace(filePath, "%20", " ", -1)

		return filePath, true
	}

	return "", false
}

func ToAbsOsPath(filePath string) (string, error) {

	if runtime.GOOS == "windows" {
		filePath = filepath.FromSlash(filePath)
	}

	return filepath.Abs(filePath)
}

func PathToFileURL(filePath string) (string, error) {

	fixedPath, err := ToAbsOsPath(filePath)

	if err != nil {
		return "", err
	}

	fixedPath = strings.Replace(fixedPath, `\`, "/", -1)

	if runtime.GOOS == "windows" {
		return "file:///" + fixedPath, nil
	} else {
		return "file:///" + fixedPath, nil
	}
}

// WriteJSONtoFile encodes the data to json and saves it to a file
func WriteJSONtoFile(filePath string, data interface{}) error {

	f, _ := os.Create(filePath)
	defer f.Close()

	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(j)
	if err != nil {
		return err
	}

	return nil
}

// CopyFile copies the file from the source to the destination file
func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			os.Chmod(dest, sourceinfo.Mode())
		}
	}

	return
}

func IsRemote(path string) bool {
	return strings.HasPrefix(path, "http")
}

func LoadRemoteFile(sourceURL string) (string, error) {

	resp, err := http.Get(sourceURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func LoadLocalFile(path string) (string, error) {

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func CopyRemoteFile(sourceURL string, dest string) (err error) {

	resp, err := http.Get(sourceURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	io.Copy(destfile, resp.Body)

	return nil
}

// CopyDir copies the specified directory and its contents to the specified destination
func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func MoveFiles(source string, dest string) (err error) {

	// get properties of source dir
	_, err = os.Stat(source)
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		srcFile := path.Join(source, obj.Name())
		destFile := path.Join(dest, obj.Name())

		if !obj.IsDir() {
			err = os.Rename(srcFile, destFile)
		}

	}
	return
}

func DeleteFilesWithPrefix(dir string, filePrefix string) int {

	deleted := 0

	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), filePrefix) {
			os.Remove(filepath.Join(dir, f.Name()))
			deleted++
		}
	}
	return deleted
}

// ExecutableExists checks if the specified executable
// exists in the users PATH
func ExecutableExists(execName string) bool {
	_, err := exec.LookPath(execName)
	if err != nil {
		return false
	}
	return true
}

// FileExists determines if the named file exists
func FileExists(filePath string) bool {

	f, err := os.Open(filePath)
	f.Close()
	if err != nil {
		return false
	}
	return true
}
