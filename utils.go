package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func fixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}

func savePage(base, uri string, httpBody io.Reader) error {
	b, err := ioutil.ReadAll(httpBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading body ", err)
		return err
	}
	full, err := url.Parse(uri)
	if err != nil {
		fmt.Fprintln(os.Stderr, "problem with parsing url")
		return err
	}

	// handle root index
	if full.Path == "/" || full.Path == "" {
		full.Path = "/index.html"
	}
	savePath := base + "/" + full.Host + full.Path
	err = os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Path %s\n", err)
		return err
	}
	fmt.Printf("Saving file: %s\n", savePath)
	f, err := os.Create(savePath)
	defer f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %s\n", err)
		return err
	}
	n, err := f.Write(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Writing %s:%s\n", savePath, err)
		return err
	}
	fmt.Fprintf(os.Stderr, "Saving %d bytes to %s\n", n, savePath)

	return nil
}

func fetch(uri string) (io.Reader, int64, error) {
	fmt.Fprintf(os.Stderr, "#")
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: transport}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, 0, err
	}

	// more politeness
	req.Header.Set("User-Agent", AGENT)
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body, resp.ContentLength, nil
}
