package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func main() {

}

// given target dir, update it's modtime to the latest modtime of any item found in this dir.
// all subdirs of the dir will also be affected (only updated by their own children, though)
func updateDirModTimeToLatest(targetDir string,dryRun bool,subDirLevel int) {
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
                subDirLevel+1,
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

    // printing only happens to the target and its immediate children
    var shouldPrint bool=subDirLevel<=1

    // if didn't change, print and quit
    if targetDirInfo.ModTime().Equal(bestModTime) {
        if shouldPrint {
            fmt.Printf("%s: no change\n",
                targetDir,
            )
        }

        return

    // otherwise, print out the intended change
    } else {
        if shouldPrint {
            fmt.Printf("%s: %s -> %s\n",
                targetDir,
                formatDate(targetDirInfo.ModTime()),
                formatDate(bestModTime),
            )
        }
    }

    // if dry run, dont make the actual change
    if dryRun {
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