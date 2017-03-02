package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/jlaffaye/ftp"
)

func getFtp(base string, ftplink *url.URL) {
	server, err := ftp.Dial(ftplink.Host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not connect to ftp server.[%s]\n", err)
		return
	}

	err = server.Login("anonymous", "anonymous")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot Login with anonymous/anonymous.[%s]\n", err)
	}

	// working on this
	if ftplink.Path == "" {
		ftplink.Path = "/"
	}
	var walk func(string, string, *ftp.ServerConn)
	walk = func(base, dir string, s *ftp.ServerConn) {
		err := s.ChangeDir(dir)
		if err != nil {
			fmt.Println(err)
			return
		}
		entries, err := s.List("")
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, entry := range entries {
			if entry.EntryType == ftp.EntryTypeFolder {
				walk(base, entry.Name, s)
				continue
			}

			savep := base + "/ftp/" + ftplink.Host + ftplink.Path

			if _, err := os.Stat(savep); !os.IsNotExist(err) {
				continue
			}

			resp, err := s.Retr(entry.Name)
			if err != nil {
				fmt.Println(err)
				continue
			}

			b, err := ioutil.ReadAll(resp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading body. [%s]\n", err)
				continue
			}

			file, err := createFile(savep)
			if err != nil {
				continue
			}

			n, err := file.Write(b)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing file. [%s]\n", err)
			}
			fmt.Printf("Saving %d bytes to %s\n", n, savep)

			file.Close()
			resp.Close()
			// list dir
			// if dir walk dir
			// if file save
		}
		// start walk with current dir
		walk(base, ftplink.Path, server)

	}
}
