#**ScrapeNSave**
A Web Scraper/Archiver written in Go.

##**Usage**
```Usage ./scrapensave url```
```  -a path```
```    	archive path. Default ./Archive```
```  -e file```
```    	exclude file type from archive (ex. "txt|mov")```
```  -f file```
```    	file of links to archive```
```  -n	do not archive
```  -p milliseconds```
```    	polite crawl delay milliseconds (default 500)```

Spider and save links to file

```./scrapensave -n http://www.example.com > links.txt```

or Archive from list of links

```./scrapensave -f links.txt```

Exclude large files from Spidering or Archiving

```./scrapensave -e "mov|jpeg|mp3"```


ScrapeNSave respects robots.txt

##**Installation**
Precompiled binary for Mac 10.11.6 or greater in binaries folder
