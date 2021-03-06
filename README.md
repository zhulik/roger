# Roger

Roger is a simple one-way SFTP synchronizer written in Go. It won't delete local files if they are deleted from the remote server.
It doesn't support resuming and doesn't track changes in the files. But it's very handy if you need simply download new
files from your server if it doesn't support rsync.

## Installation

`go get github.com/zhulik/roger`

## Usage

### Single run mode
`roger -local=/home/user/files -remote=sftp://server.example:22/path/files -workers 16`

### Daemon mode

`roger -local=/home/user/files -remote=sftp://server.example:22/path/files -workers 16 -daemon -interval 120`

## Roadmap

- Resuming
- Proper logging when works as a daemon
- Event hooks(started, in progress, finished)
