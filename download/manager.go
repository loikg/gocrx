package download

import (
	"log"
	"os"
	"sync"

	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

// Manager represent a Download manager with workers
type Manager struct {
	wg        *sync.WaitGroup
	nbWorkers int
	jobs      chan *Job
	pool      *mpb.Progress
}

// Job describe what extension a worker will download
type Job struct {
	Name, URL, Output string
	bar               *mpb.Bar
}

// New create a new Manager to mange concurrent downloads
func New(nbWorkers int) *Manager {
	wg := new(sync.WaitGroup)
	return &Manager{
		wg:        wg,
		nbWorkers: nbWorkers,
		jobs:      make(chan *Job, 100),
		pool:      mpb.New(mpb.WithWaitGroup(wg)),
	}
}

// Start start the download manager. It's worker will wait for job added via Add(job Job).
func (dl *Manager) Start() {
	for i := 0; i < dl.nbWorkers; i++ {
		go dl.worker()
	}
}

func (dl *Manager) PostJob(job *Job) {
	job.bar = dl.buildBar(job)
	dl.jobs <- job
	dl.wg.Add(1)
}

// Wait will stop all workers after they finished processing jobs
func (dl *Manager) WaitAndStop() {
	close(dl.jobs)
	dl.pool.Wait()
}

func (dl *Manager) worker() {
	for job := range dl.jobs {
		log.Printf("Downloading %s at %s to %s\n", job.Name, job.URL, job.Output)
		dl.handleJob(job)
	}
}

func (dl *Manager) handleJob(job *Job) {
	defer dl.wg.Done()
	f, err := os.Create(job.Output)
	if err != nil {
		log.Printf("failed to create file '%s': %v\n", job.Output, err)
		job.bar.Abort(false)
	}
	defer f.Close()
	if err := downloadFileTo(job.URL, f, job.bar); err != nil {
		job.bar.Abort(false)
		if err := os.Remove(f.Name()); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
		log.Printf("failed to download '%s' at '%s': %v\n", job.Name, job.URL, err)
	}
}

func (dl *Manager) buildBar(job *Job) *mpb.Bar {
	return dl.pool.AddBar(
		100,
		mpb.PrependDecorators(
			// simple name decorator
			decor.Name(job.Name, decor.WCSyncWidthR),
			// decor.DSyncWidth bit enables column width synchronization
			//decor.Percentage(decor.WCSyncSpace),
			decor.CountersKibiByte(" % 6.1f / % 6.1f", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			// replace ETA decorator with "done" message, OnComplete event
			decor.OnComplete(
				decor.Percentage(), "done !",
			),
		),
	)
}
