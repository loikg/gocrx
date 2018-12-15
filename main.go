package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var chromeVersion string
var outPutPath string
var inputPath string

const downloadURL = "https://clients2.google.com/service/update2/crx?response=redirect&acceptformat=crx2,crx3&prodversion=[VERSION]&x=id%3D[EXTENSION_ID]%26installsource%3Dondemand%26uc"

func init() {
	flag.StringVar(&chromeVersion, "version", "70.0", "Chrome version to use")
	flag.StringVar(&outPutPath, "output", "./", "Output path where downloaded crx file will be sotred")
	flag.StringVar(&inputPath, "file", "extension.txt", "File containing list of extension to download")
	flag.Parse()
}

func main() {
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	readFile(file)
}

func readFile(r io.Reader) {
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		wg.Add(1)
		go func(extURL string) {
			defer wg.Done()
			downloadExtension(extURL)
		}(scanner.Text())
	}
	wg.Wait()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func downloadExtension(extURL string) error {
	fmt.Printf("Dowloading: %s\n", extURL)
	ext, err := parseExtensionURL(extURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %s: %v", extURL, err)
	}
	err = ext.Download(outPutPath)
	if err != nil {
		return fmt.Errorf("failed to download: %s: %v", ext.Name, err)
	}
	return nil
}

type extension struct {
	Name        string
	ID          string
	DownloadURL string
}

// ParseExtensionURL fill an extension struct with information of the corresponding extension
func parseExtensionURL(extURL string) (*extension, error) {
	_, err := url.Parse(extURL)
	if err != nil {
		return nil, fmt.Errorf("failed to patse URL: %v", err)
	}
	parts := strings.Split(extURL, "/")
	extID := parts[len(parts)-1]
	extName := fmt.Sprintf("%s.crx", parts[len(parts)-2])
	downloadURL := strings.Replace(downloadURL, "[VERSION]", chromeVersion, 1)
	downloadURL = strings.Replace(downloadURL, "[EXTENSION_ID]", extID, 1)
	return &extension{
		DownloadURL: downloadURL,
		Name:        extName,
		ID:          extID,
	}, nil
}

func (ext extension) Download(dest string) error {
	res, err := http.Get(ext.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download .crx: %v", err)
	}
	output := filepath.Join(outPutPath, ext.Name)
	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create %s file: %v", ext.Name, err)
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return fmt.Errorf("failed to write to destination file %s: %v", ext.Name, err)
	}
	return nil
}
