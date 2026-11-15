package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// args from cli
type CliArgs struct {
    DryRun bool
    TargetPath string
}

func main() {
    var args CliArgs=getCliArgs()

    updateDirModTimeToLatest(
        args.TargetPath,
        args.DryRun,
        0,
    )
}

// given target dir, update it's modtime to the latest modtime of any item found in this dir.
// all subdirs of the dir will also be affected (only updated by their own children, though)
func updateDirModTimeToLatest(targetDir string,dryRun bool,subDirLevel int) time.Time {
    targetDir=filepath.Clean(targetDir)

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
        var itemModTime time.Time

        // if item is a dir, call this func on it, and set the item's mod time to the returned result
        if item.IsDir() {
            itemModTime=updateDirModTimeToLatest(
                filepath.Join(targetDir,item.Name()),
                dryRun,
                subDirLevel+1,
            )

        // if not a dir, set it to normal fileinfo mod time
        } else {
            // get the mod time of the item, which for dirs should be the time after
            // being properly updated
            var info fs.FileInfo
            info,e=item.Info()

            if e!=nil {
                panic(e)
            }

            itemModTime=info.ModTime()
        }

        if itemModTime.After(bestModTime) {
            bestModTime=itemModTime
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

        return bestModTime

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
        return bestModTime
    }

    e=os.Chtimes(targetDir,time.Now(),bestModTime)

    if e!=nil {
        panic(e)
    }

    return bestModTime
}

// format for print
func formatDate(date time.Time) string {
    return date.Format("2006-01-02 15:04:05")
}

// get cli args
func getCliArgs() CliArgs {
    var dryRun *bool = flag.Bool("d", false, "dry run (no changes)")
    var targetPath *string = flag.String("t", "", "target directory (required)")

    flag.Parse()

    if *targetPath == "" {
        fmt.Fprintln(os.Stderr, "Error: target directory (-t) is required.")
        flag.Usage()
        os.Exit(1)
    }

    return CliArgs{
        DryRun:     *dryRun,
        TargetPath: filepath.Clean(*targetPath),
    }
}