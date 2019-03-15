package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/loikg/gocrx/downloader"
	"gopkg.in/cheggaaa/pb.v1"
)

/**
* Command line flags value
 */
var (
	chromeVersion string
	outPutPath    string
	inputPath     string
	extID         string
)

const downloadURL = "https://clients2.google.com/service/update2/crx?response=redirect&acceptformat=crx2,crx3&prodversion=[VERSION]&x=id%3D[EXTENSION_ID]%26installsource%3Dondemand%26uc"

func init() {
	flag.StringVar(&chromeVersion, "version", "72.0", "Chrome version to use")
	flag.StringVar(&outPutPath, "output", "./", "Output path where downloaded crx file will be sotred")
	flag.StringVar(&inputPath, "file", "extension.txt", "File containing list of extension to download")
	flag.StringVar(&extID, "id", "", "ID of the extension to download.")
	flag.Parse()
}

func main() {
	if isExtIDFlagSet() {
		if err := downloadExtensionByID(extID); err != nil {
			fmt.Printf("Failed to download: %v", err)
			os.Exit(1)
		}
	} else {
		if err := downloadFromFile(inputPath, outPutPath); err != nil {
			fmt.Printf("Failed to download: %v", err)
			os.Exit(1)
		}
	}
}

// Check if extID flag as been set on the command line.
// e.g it's default value is an empty string
func isExtIDFlagSet() bool {
	return extID != ""
}

func downloadExtensionByID(ID string) error {
	fmt.Printf("Downloading %s to %s ...\n", ID, outPutPath)
	URL := buildDownloadURL(ID, chromeVersion)
	f, err := os.Create(outPutPath)
	if err != nil {
		return err
	}
	defer f.Close()
	bar := pb.New64(0)
	if err := downloader.DownloadFileTo(URL, f, bar); err != nil {
		os.Remove(f.Name())
		return err
	}
	return nil
}

// buildDownloadURL return the download URL by replacing the extension id (extID) and
// chrome version (chromeVersion) in the base download URL format.
func buildDownloadURL(extID, chromeVersion string) string {
	URL := strings.Replace(downloadURL, "[VERSION]", chromeVersion, 1)
	URL = strings.Replace(URL, "[EXTENSION_ID]", extID, 1)
	return URL
}

const (
	extensionName = 0
	extensionID   = 1
)

func parseFile(in io.Reader, outputPath string) ([]downloader.Job, error) {
	jobs := make([]downloader.Job, 0)
	nbLine := 0
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		linePart := strings.SplitN(line, ":", 2)
		if len(linePart) != 2 {
			fmt.Printf("Skip line %d: %s", nbLine, line)
			continue
		}
		extID := strings.TrimSpace(linePart[extensionID])
		extName := strings.TrimSpace(linePart[extensionName])
		jobs = append(jobs, downloader.Job{
			Name: extName,
			Path: path.Join(outputPath, extName+".crx"),
			ID:   buildDownloadURL(extID, chromeVersion),
		})
		nbLine++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return jobs, nil
}

func downloadFromFile(filePath, outputPath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	jobs, err := parseFile(f, outPutPath)
	if err != nil {
		return err
	}

	manager := downloader.New(5, jobs)
	manager.Start()
	manager.Wait()
	return nil
}
