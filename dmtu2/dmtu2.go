package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func main() {

}

// given target dir, update it's modtime to the latest modtime of any item found in this dir.
// all subdirs of the dir will also be affected (only updated by their own children, though)
func updateDirModTimeToLatest(targetDir string,dryRun bool) {
	var items []os.DirEntry
	var e error
	items,e=os.ReadDir(targetDir)

    if e!=nil {
        panic(e)
    }

    var targetDirInfo fs.FileInfo
    targetDirInfo,e=os.Stat(targetDir)

    if e!=nil {
        panic(e)
    }

    var bestModTime time.Time=targetDirInfo.ModTime()

    var item os.DirEntry
    for _,item = range items {
        // if item is a dir, call this func on it, so its modtime gets updated
        if item.IsDir() {
            updateDirModTimeToLatest(
                filepath.Join(targetDir,item.Name()),
                dryRun,
            )
        }

        // get the mod time of the item, which for dirs should be the time after
        // being properly updated
        var info fs.FileInfo
        info,e=item.Info()

        if e!=nil {
            panic(e)
        }

        if info.ModTime().After(bestModTime) {
            bestModTime=info.ModTime()
        }
    }

    if dryRun {
        return
    }

    if targetDirInfo.ModTime().Equal(bestModTime) {
        return
    }

    e=os.Chtimes(targetDir,time.Now(),bestModTime)

    if e!=nil {
        panic(e)
    }
}

// format for print
func formatDate(date time.Time) string {
    return date.Format("2006-01-02 15:04")
}