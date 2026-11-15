Program that can update a folder's mod time to be the latest mod time of any item within that folder (recursively). Also affects all sub dirs of the target folder.

# Usage
```bash
dmtu2 -t <target-dir> [-d]
```

- `-t`: Target directory (required).
- `-d`: Dry run (no changes).
- `-h`: Display help.

# Development
In `dmtu2`, `bash build.sh` or `bash build-run.sh`