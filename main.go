package main

import (
	"flag"
	"log"
	"net/url"
	"time"

	"github.com/pkg/sftp"
	"github.com/zhulik/roger/syncer"
	"golang.org/x/crypto/ssh"
)

var (
	localDir  = flag.String("local", "", "Local directory path")
	remoteURL = flag.String("remote", "", "Remote directory URL")
	workers   = flag.Int("workers", 16, "Count of download workers")

	// daemon options
	daemon   = flag.Bool("daemon", false, "Run in daemon mode")
	interval = flag.Int64("interval", 120, "Interval between syncronizations (seconds)")
)

func runSync(url *url.URL, config ssh.ClientConfig) {
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

	if *daemon {
		log.Println("Running in daemon mode...")
		for {
			runSync(url, config)
			log.Printf("Waiting %d seconds...", *interval)
			time.Sleep(time.Duration(*interval) * time.Second)
		}
	} else {
		runSync(url, config)
	}
}
