# build run with some test arguments

set -exu
HERE=$(dirname $(realpath $BASH_SOURCE))
cd $HERE

go build -o dmtu2.exe dmtu2.go
./dmtu2.exe