package downloader

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gopkg.in/cheggaaa/pb.v1"
)

// Downloader represent a Download manager with workers
type Downloader struct {
	nbWorkers int
	jobs      chan Job
	wg        sync.WaitGroup
	pbPool    *pb.Pool
	jobList   []Job
}

// Job describe what extension a worker will downaload
type Job struct {
	Name string
	ID   string
	Path string
	bar  *pb.ProgressBar
}

// New create a new Downloder to mange concurrent downloads
func New(nbWorkers int, jobs []Job) *Downloader {
	nbJobs := len(jobs)
	dl := &Downloader{
		nbWorkers: nbWorkers,
		jobs:      make(chan Job, nbJobs),
		pbPool:    initializeProgressBarPool(jobs),
		jobList:   jobs,
	}
	dl.wg.Add(nbJobs)
	return dl
}

// Start start the download manager. It's worker will wait for job added via Add(job Job).
func (dl *Downloader) Start() {
	for i := 0; i < dl.nbWorkers; i++ {
		go dl.worker()
	}
	if err := dl.pbPool.Start(); err != nil {
		log.Printf("error: %v\n", err)
	}
	for _, job := range dl.jobList {
		dl.jobs <- job
	}
	close(dl.jobs)
}

// Wait will stop all workers after they finished processing jobs
func (dl *Downloader) Wait() {
	dl.wg.Wait()
}

func (dl *Downloader) worker() {
	for job := range dl.jobs {
		// log.Printf("Processing: %v\n", job)
		defer dl.wg.Done()
		if err := dl.handleJob(job); err != nil {
			log.Printf("failed to create file: %v\n", err)
		}
	}
}

func (dl *Downloader) handleJob(job Job) error {
	f, err := os.Create(job.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	log.Println("handleJob()", job)
	if err := DownloadFileTo(job.ID, f, job.bar); err != nil {
		os.Remove(f.Name())
		return err
	}
	return nil
}

// initializeProgressBarPool create a progress bar pool with nbWorkers progress bars.
func initializeProgressBarPool(jobs []Job) *pb.Pool {
	progressBars := make([]*pb.ProgressBar, 0, len(jobs))
	for i, job := range jobs {
		// for _, job := range jobs {
		// }
		job.bar = pb.New64(0).
			SetUnits(pb.U_BYTES).
			Prefix(fmt.Sprintf("%s - Waiting", job.Name))
		progressBars = append(progressBars, job.bar)
		jobs[i] = job
		log.Println(job)
	}
	return pb.NewPool(progressBars...)
}
