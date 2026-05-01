package worker

import (
	"fmt"
	"sync"
	"time"

	"cheat-master/internal/client"
	"cheat-master/internal/models"
	"cheat-master/internal/storage"
)

// WorkerPool manages concurrent job execution
type WorkerPool struct {
	maxWorkers    int
	jobChannel    chan *models.Job
	wg            sync.WaitGroup
	ctx            *WorkerContext
	credManager   *storage.CredentialsManager
	jobManager    *storage.JobManager
	activeJobs    sync.Map // map[string]bool
	resourceLimit int      // Max concurrent jobs
}

// WorkerContext holds context for job execution
type WorkerContext struct {
	MaxRetries     int
	MaxPolls       int
	TimeoutSeconds int
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(maxWorkers int, credManager *storage.CredentialsManager, jobManager *storage.JobManager) *WorkerPool {
	return &WorkerPool{
		maxWorkers:   maxWorkers,
		jobChannel:   make(chan *models.Job, maxWorkers*2),
		credManager:  credManager,
		jobManager:   jobManager,
		resourceLimit: maxWorkers,
		ctx: &WorkerContext{
			MaxRetries:     3,
			MaxPolls:       6,
			TimeoutSeconds: 300,
		},
	}
}

// Start begins processing jobs
func (wp *WorkerPool) Start() {
	fmt.Printf("🚀 Starting worker pool with %d workers\n", wp.maxWorkers)

	for i := 0; i < wp.maxWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i + 1)
	}

	fmt.Println("✅ Worker pool started")
}

// Stop gracefully stops the worker pool
func (wp *WorkerPool) Stop() {
	fmt.Println("⏹ Stopping worker pool...")
	close(wp.jobChannel)
	wp.wg.Wait()
	fmt.Println("✅ Worker pool stopped")
}

// SubmitJob adds a job to the queue
func (wp *WorkerPool) SubmitJob(job *models.Job) error {
	// Check resource availability
	activeCount := wp.getActiveJobCount()
	if activeCount >= wp.resourceLimit {
		return fmt.Errorf("resource limit reached: %d/%d", activeCount, wp.resourceLimit)
	}

	select {
	case wp.jobChannel <- job:
		fmt.Printf("📥 Job submitted: %s\n", job.ID)
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout submitting job")
	}
}

// worker processes jobs from the channel
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	fmt.Printf("👷 Worker %d started\n", id)

	for job := range wp.jobChannel {
		wp.processJob(id, job)
	}

	fmt.Printf("👷 Worker %d stopped\n", id)
}

// processJob executes a job
func (wp *WorkerPool) processJob(workerID int, job *models.Job) {
	wp.activeJobs.Store(job.ID, true)
	defer wp.activeJobs.Delete(job.ID)

	fmt.Printf("⚙️  Worker %d: Processing job %s\n", workerID, job.ID)

	// Update job status to running
	wp.jobManager.UpdateJobStatus(job.ID, "running")

	// Get user credentials
	user, err := wp.credManager.GetUser(job.UserID)
	if err != nil {
		wp.jobManager.UpdateJobStatus(job.ID, "failed")
		fmt.Printf("❌ Worker %d: User not found: %v\n", workerID, err)
		return
	}

	// Create client and login
	lectureClient := client.NewClient()
	if err := lectureClient.Login(user.Email, user.Password); err != nil {
		wp.jobManager.UpdateJobStatus(job.ID, "failed")
		fmt.Printf("❌ Worker %d: Login failed: %v\n", workerID, err)
		return
	}

	// Process the course
	// TODO: Implement course processing logic
	fmt.Printf("✅ Worker %d: Job %s completed\n", workerID, job.ID)
	wp.jobManager.UpdateJobStatus(job.ID, "completed")
}

// GetStats returns worker pool statistics
func (wp *WorkerPool) GetStats() map[string]interface{} {
	activeCount := wp.getActiveJobCount()
	pendingJobs := wp.jobManager.ListPendingJobs()

	return map[string]interface{}{
		"max_workers":    wp.maxWorkers,
		"active_jobs":    activeCount,
		"pending_jobs":   len(pendingJobs),
		"resource_limit": wp.resourceLimit,
		"queue_size":     len(wp.jobChannel),
	}
}

// getActiveJobCount returns number of active jobs
func (wp *WorkerPool) getActiveJobCount() int {
	count := 0
	wp.activeJobs.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}
