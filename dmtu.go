package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type LatestModFileResult struct {
    Filepath string
    ModTime time.Time
}

func main() {
	// var testDir string = `C:\Users\ktkm\Desktop\test`

    // fmt.Println(findLatestFileModTime(testDir,true))

	updateModTimeWithSubItems(`C:\Users\ktkm\Desktop\test\pappi`)
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

// target a folder and set its mod time to the latest mod time of an item found in the folder
func updateModTimeWithSubItems(targetDir string) {
    var result LatestModFileResult
    result = findLatestFileModTime(targetDir, false)

    if result.Filepath == "" {
        fmt.Printf("No files found in %s\n", targetDir)
        return
    }

    var info fs.FileInfo
    var e error

    info, e = os.Stat(targetDir)
    if e != nil {
        panic(e)
    }

    e = os.Chtimes(targetDir, info.ModTime(), result.ModTime)
    if e != nil {
        panic(e)
    }

	fmt.Printf("%s: %s -> %s\n",
		targetDir,
		formatDate(info.ModTime()),
		formatDate(result.ModTime),
	)
}

// format for print
func formatDate(date time.Time) string {
    return date.Format("2006-01-02 15:04")
}