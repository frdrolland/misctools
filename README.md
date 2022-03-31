# misctools

# Prerequisites
Install Golang SDK 1.18+ : https://go.dev/dl/

# Get the sources
    cd $GOPATH/src<br>
    mkdir -p github.com/frdrolland<br>
    git clone https://github.com/frdrolland/misctools<br>
    cd misctools<br>

# Compilation

Cross-compilation: as it is written in Go language, binary file can be cross-compiled to build binary for another OS.

For example, from a Windows machine, you can build misctools.exe windows binary file with this command :<br>

    go build<br>
But you can also build misctools binary for Linux 64 bits using :<br>

    env GOOS=windows GOARCH=amd64 go build<br>
    
