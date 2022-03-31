# misctools

# Prerequisites
Install Golang SDK 1.18+ : https://go.dev/dl/

# Get the sources
cd $GOPATH/src
mkdir -p github.com/frdrolland
git clone https://github.com/frdrolland/misctools
cd misctools

# Compilation

Cross-compilation: as it is written in Go language, binary file can be cross-compiled to build binary for another OS.

For example, from a Windows machine, you can build misctools.exe windows binary file with this command :
go build
But you can also build misctools binary for Linux 64 bits using :
env GOOS=windows GOARCH=amd64 go build

