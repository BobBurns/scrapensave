package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/temoto/robotstxt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

const PATH string = "Archive"
const AGENT string = "scrapensavebot/1.0"

var visited = make(map[string]bool)

func generator(link string) (<-chan string, <-chan int64) {
	c := make(chan string)
	d := make(chan int64)
	go func() {
		body, count, err := fetch(link)
		if err == nil {
			if count < 0 {
				count = 0
			}
			d <- count
			links := collectLinks(body)
			for _, new := range links {
				absolute := fixUrl(new, link)
				c <- absolute
			}
		}
		close(c)
		close(d)
	}()
	return c, d
}

type Archive struct {
	ArchPath string
	FilePath string
	Exclude  *regexp.Regexp
	Narchive bool
	Delay    time.Duration
	Allow    *robotstxt.Group
	Domain   *regexp.Regexp
	Bytes    int64
	Links    []string
}

func (a *Archive) buildArchive() {
	var polite time.Duration
	// get command line options
	fFlag := flag.String("f", "", "`file` of links to archive")
	eFlag := flag.String("e", "", "exclude `file` type from archive (ex. \"txt|mov\")")
	nFlag := flag.Bool("n", false, "do not archive")
	pFlag := flag.Int("p", 500, "polite crawl delay `milliseconds`")
	aFlag := flag.String("a", "", "archive `path`. Default ./Archive")

	flag.Parse()

	exclude := fmt.Sprintf("[A-Za-z0-9]+\\.(%s)$", *eFlag)
	excludeExp, err := regexp.Compile(exclude)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad -e option. Usage: ")
		flag.PrintDefaults()
		os.Exit(1)
	}

	seed := flag.Arg(0)
	if seed == "" && *fFlag == "" {
		fmt.Fprintf(os.Stderr, "Usage %s url\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	seedUri, err := url.Parse(seed)
	if err != nil {
		fmt.Fprintln(os.Stderr, "please use valid url")
		os.Exit(1)
	}

	polite = time.Duration(*pFlag) * time.Millisecond

	domain := regexp.MustCompile(seedUri.Host)

	//check robots.txt
	robots := &robotstxt.RobotsData{}
	rgroup := &robotstxt.Group{}
	resp, err := http.Get(seedUri.Scheme + "://" + seedUri.Host + "/robots.txt")
	if err == nil {
		robots, err = robotstxt.FromResponse(resp)
		if err == nil {
			rgroup = robots.FindGroup(AGENT)
			polite = rgroup.CrawlDelay
			fmt.Fprintln(os.Stderr, "Found robots.txt with time delay ", polite)
		} else {
			fmt.Fprintln(os.Stderr, "bad robots.txt")
		}

	}

	if *aFlag == "" {
		a.ArchPath = PATH
	} else {
		a.ArchPath = *aFlag
	}
	a.FilePath = *fFlag
	a.Narchive = *nFlag
	a.Exclude = excludeExp
	a.Delay = polite
	a.Allow = rgroup
	a.Domain = domain
	a.Links = []string{seed}

}

func main() {
	var linkReader *bufio.Reader

	arch := Archive{}
	arch.buildArchive()

	// if f flag and file, archive only links from file
	// first scan links to arch.Links
	if arch.FilePath != "" {
		linkFile, err := os.Open(arch.FilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot open %s for reading\n", arch.FilePath)
			os.Exit(1)
		}
		linkReader = bufio.NewReader(linkFile)
		scanner := bufio.NewScanner(linkReader)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			arch.Links = append(arch.Links, scanner.Text())
		}
		goto archive
	}

	// spider all links in domain
	fmt.Fprintln(os.Stderr, "Scanning...")

	for i := 0; i < len(arch.Links); i++ {
		c, d := generator(arch.Links[i])
		arch.Bytes += <-d // this will return first
		for link := range c {
			if !arch.Allow.Test(link) || !arch.Domain.MatchString(link) {
				continue
			}
			if arch.Exclude.MatchString(link) {
				continue
			}

			full, _ := url.Parse(link)
			if !visited[full.Path] {
				arch.Links = append(arch.Links, link)
				visited[full.Path] = true

			}
		}
		time.Sleep(arch.Delay)
	}

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Approximate bytes to archive = %d\n", arch.Bytes)
	if !arch.Narchive {
		answer := ""
		fmt.Fprintf(os.Stderr, "Are you sure you still want to save? [Y/n] ")
		fmt.Scanf("%s", &answer)
		if answer == "n" {
			arch.Narchive = true
		}
	}

archive:
	for _, link := range arch.Links {
		// archive if n flag is not set
		if !arch.Narchive {
			body, _, err := fetch(link)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching %s to save.\n[%s]\n", link, err)
				continue
			}
			err = savePage(arch.ArchPath, link, body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error with Save.\n[%s]\n", err)
				continue
			}
			time.Sleep(arch.Delay)
		}

		fmt.Println(link)
	}

}
