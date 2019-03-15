package download

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/vbauerster/mpb/v4"
)

// downloadFileTo retrieve the resource at the given URL and write it to dst.
// It update the given progress bar and estimate the size of the resource before downloading it with a HEAD HTTP
// request.
func downloadFileTo(url string, dst io.Writer, bar *mpb.Bar) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 204 {
		return fmt.Errorf("request failed with %d %s", resp.StatusCode, resp.Status)
	}
	totalSize, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}
	bar.SetTotal(totalSize, false)
	reqBody := bar.ProxyReader(resp.Body)
	_, err = io.Copy(dst, reqBody)
	if err != nil {
		return err
	}
	return nil
}

// estimateFileSize estimate the size of resource by making an HEAD HTTP request and parsing the Content-Length header.
//func estimateFileSize(url string) (int64, error) {
//	resp, err := http.Head(url)
//	if err != nil {
//		return 0, err
//	}
//	defer resp.Body.Close()
//	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
//	if err != nil {
//		return 0, err
//	}
//	return size, nil
//}