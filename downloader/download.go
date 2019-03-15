package downloader

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/cheggaaa/pb.v1"
)

// DownloadFileTo download the extension to the given file and show a progress bar.
func DownloadFileTo(url string, f *os.File, bar *pb.ProgressBar) error {
	// log.Printf("Downloading %s\n", url)
	totalSize, err := getExtFileSize(url)
	if err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bar.Set64(totalSize).SetUnits(pb.U_BYTES).Prefix("downloading").Start()
	bar.Start()
	defer bar.Finish()
	r := bar.NewProxyReader(resp.Body)
	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}
	return nil
}

// getExtFileSize make an HTTP HEAD request to retrieve the Content-Length.
func getExtFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}
