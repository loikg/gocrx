package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/loikg/gocrx/download"
	"github.com/spf13/cobra"
)

// Command line flags
var (
	chromeVersion string
	nbWorkers     int
)

// Global constant
const (
	downloadURLFormat = "https://clients2.google.com/service/update2/crx?response=redirect&acceptformat=crx2,crx3&prodversion={prodVersion}&x=id%3D{extensionID}%26uc"
	// Indexes used when splitting extension file
	extensionName     = 0
	extensionID       = 1
)

type app struct {
	dlManager *download.Manager
}

func init() {
	// Define flag to retrieve chrome's version for which extension are downloaded.
	rootCmd.PersistentFlags().StringVarP(
		&chromeVersion,
		"chrome",
		"c",
		"72.0",
		"Chrome version for which extension are downloaded",
	)
	// Define flag to retrieve the number of workers to use.
	rootCmd.PersistentFlags().IntVarP(
		&nbWorkers,
		"worker",
		"w",
		4,
		"Number of parallel workers",
	)
}

var rootCmd = &cobra.Command{
	Use:   "gocrx <file|id> [destination]",
	Short: "Quickly download chrome extension.",
	Long: `A tool to download chrome extension .crx files.
Can read from a file or download by extension id.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app := NewApp(nbWorkers)
		app.Run(args)
	},
}

func NewApp(nbWorkers int) *app {
	return &app{
		dlManager: download.New(nbWorkers),
	}
}

func (app *app) Run(args []string) {
	app.dlManager.Start()
	if fileExist(args[0]) {
		f, err := os.Open(args[0])
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		defer f.Close()
		if err := app.downloadExtensionFromFile(f, getOutput(args, "./")); err != nil {
			log.Printf("error: %v\n", err)
		}
	} else {
		app.downloadExtensionByID(args[0], getOutput(args, "download.crx"))
	}
	app.dlManager.WaitAndStop()
}

func getOutput(args []string, fallback string) string {
	var output string
	if len(args) > 1 {
		output = args[1]
	} else {
		output = fallback
	}
	return output
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (app *app) downloadExtensionByID(ID, output string) {
	app.dlManager.PostJob(&download.Job{
		Name:   "#1",
		URL:    buildDownloadURL(ID, chromeVersion),
		Output: output,
	})
}

func (app *app) downloadExtensionFromFile(src io.Reader, outPutPath string) error {
	scanner := bufio.NewScanner(src)
	nbLine := 0
	for scanner.Scan() {
		line := scanner.Text()
		linePart := strings.SplitN(line, ":", 2)
		if len(linePart) != 2 {
			log.Printf("Skip line %d: %s", nbLine, line)
			nbLine++
			continue
		}
		extName := strings.TrimSpace(linePart[extensionName])
		app.dlManager.PostJob(&download.Job{
			Name:    extName,
			URL:    buildDownloadURL(strings.TrimSpace(linePart[extensionID]), chromeVersion),
			Output: path.Join(outPutPath, extName+".crx"),
		})
		nbLine++
	}
	return nil
}

// buildDownloadURL return the download URL by replacing the extension id (extID) and
// chrome version (chromeVersion) in the base download URL format.
func buildDownloadURL(extID, chromeVersion string) string {
	URL := strings.Replace(downloadURLFormat, "{prodVersion}", chromeVersion, 1)
	URL = strings.Replace(URL, "{extensionID}", extID, 1)
	return URL
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}