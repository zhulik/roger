package main

import (
	"flag"
	"net/url"

	"github.com/zhulik/roger/syncer"

	"log"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	localDir  = flag.String("local", "", "Local directory path")
	remoteURL = flag.String("remote", "", "Remote directory URL")
	workers   = flag.Int("workers", 8, "Count of download workers")
)

func main() {
	flag.Parse()
	if *localDir == "" {
		panic("Please, specify local dir, see --help")
	}

	if *remoteURL == "" {
		panic("Please, specify remote URL, see --help")
	}

	url, err := url.Parse(*remoteURL)
	if err != nil {
		panic(err)
	}

	username := url.User.Username()
	password, _ := url.User.Password()

	config := ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	log.Println("Connecting...")

	conn, err := ssh.Dial("tcp", url.Host, &config)
	if err != nil {
		panic(err)
	}

	sconn, err := sftp.NewClient(conn, sftp.MaxPacket(1<<15))
	if err != nil {
		panic(err)
	}
	defer sconn.Close()

	syncer.Sync(sconn, *localDir, url.Path, *workers)
}
