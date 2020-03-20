package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Package struct {
	ApkName string `xml:"apkname"`
}

type Application struct {
	Packages []Package `xml:"package"`
}

type Fdroid struct {
	Applications []Application `xml:"application"`
}

func fetchIndexJar(url *url.URL) (buf []byte, err error) {
	var resp *http.Response
	resp, err = http.Get(url.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	buf, err = ioutil.ReadAll(resp.Body)
	return
}

func extractIndexXml(indexJar []byte) (fdroid Fdroid, err error) {
	buf := bytes.NewReader(indexJar)
	r, err := zip.NewReader(buf, buf.Size())
	if err != nil {
		return
	}
	for _, f := range r.File {
		if f.Name == "index.xml" {
			var rc io.ReadCloser
			rc, err = f.Open()
			if err != nil {
				return
			}
			defer rc.Close()
			dec := xml.NewDecoder(rc)
			err = dec.Decode(&fdroid)
			return
		}
	}
	err = fmt.Errorf("missing index.xml")
	return
}

func handle(repo string) (apks []*url.URL, err error) {
	baseUrl, err := url.Parse(repo)
	if err != nil {
		return nil, err
	}
	indexRef, err := url.Parse("index.jar")
	if err != nil {
		panic(err)
	}
	indexUrl := baseUrl.ResolveReference(indexRef)
	indexJar, err := fetchIndexJar(indexUrl)
	if err != nil {
		return nil, fmt.Errorf("error while fetching index.jar: %w", err)
	}
	indexXml, err := extractIndexXml(indexJar)
	if err != nil {
		return nil, fmt.Errorf("error while extracting index.xml: %w", err)
	}
	failedApks := 0
	for _, app := range indexXml.Applications {
		for _, pkg := range app.Packages {
			apknameRef, err := url.Parse(pkg.ApkName)
			if err != nil {
				failedApks++
				log.Printf("encountered apkname badly parsing as url part: %#v", pkg.ApkName)
			} else {
				apks = append(apks, indexUrl.ResolveReference(apknameRef))
			}
		}
	}
	err = nil
	if failedApks > 0 {
		err = fmt.Errorf("encountered %d apknames badly parsing as url parts", failedApks)
	}
	return
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("please provide repositories (for example http://192.168.2.55:8888/fdroid/repo/)")
	}
	for _, repo := range flag.Args() {
		log.Printf("processing repository %#v", repo)
		newApks, err := handle(repo)
		log.Printf("captured %d apps from repository %#v", len(newApks), repo)
		if err != nil {
			log.Print(fmt.Errorf("error while handling repository %#v: %w", repo, err))
		}
		for _, apk := range newApks {
			fmt.Println(apk.String())
		}
	}
}
