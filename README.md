ScrapeNSave
===========

A Web Scraper/Archiver written in Go.

Usage
-----

```
Usage ./scrapensave url
  -a path
    	archive path. Default ./Archive
  -e file
    	exclude file type from archive (ex. "txt|mov")
  -f file
    	file of links to archive
  -n	do not archive
  -p milliseconds
    	polite crawl delay milliseconds (default 500)
```

Archive all web pages in domain www.example.com to ./Archive

`./scrapensave https://www.example.com`

Spider and save links to file

`./scrapensave -n http://www.example.com > links.txt`

or Archive from list of links

`./scrapensave -f links.txt`

Exclude large files from Spidering or Archiving

`./scrapensave -e "mov|jpeg|mp3"`


ScrapeNSave respects robots.txt

Installation
------------

Precompiled binaries for Mac (compiled on El Capitan 10.11.6) and Ubuntu (Compiled on Mint 18) in binaries folder

Or make executable scrapensave with `go build`

You will need package robotstxt

`go get github.com/temoto/robotstxt`


**Caution**

Running this program can fill up lots of disk space.  I would recommend mounting an external drive and giving scrapensave the full path of the drive with the -a option.

If errors occur when archiving, link urls will be appended to scrapensave.log

**Caveats**

When archiving a large number of files, you may want to change ```ulimit -n``` to meet a large number of file descriptors being open. I am using ```ulimit -n 15000``` on Ubuntu.


Inspired from Jack Danger's [gocrawler](https://github.com/JackDanger/gocrawler)
Checkout [6brand.com](https://jdanger.com/)

