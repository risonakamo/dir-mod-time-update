package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"
)

type LatestModFileResult struct {
    Filepath string
    ModTime time.Time
}

func main() {
	var testDir string = `C:\Users\ktkm\Desktop\test`

    fmt.Println(findLatestFileModTime(testDir,true))
}

// finds the latest modified time of any item (dir or file) in the target dir
func findLatestFileModTime(targetDir string, onlyFiles bool) LatestModFileResult {
	var latest time.Time
    var foundFile string=""
    var e error

	e = filepath.WalkDir(targetDir, func(path string, entry fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

        var info fs.FileInfo
		info, e = entry.Info()
		if e != nil {
			return e
		}

        if onlyFiles && info.IsDir() {
            return nil
        }

		if info.ModTime().After(latest) {
			latest = info.ModTime()
            foundFile=path
		}

		return nil
	})

	if e != nil {
		panic(e)
	}

	return LatestModFileResult{
        Filepath: foundFile,
        ModTime: latest,
    }
}