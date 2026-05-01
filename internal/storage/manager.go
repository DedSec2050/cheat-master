package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"cheat-master/internal/models"
)

const (
	CredentialsFile = "credentials.json"
	JobsFile        = "jobs.json"
)

// CredentialsManager handles user credential storage and retrieval
type CredentialsManager struct {
	mu    sync.RWMutex
	store *models.CredentialsStore
	file  string
}

// NewCredentialsManager creates a new credentials manager
func NewCredentialsManager(filePath string) (*CredentialsManager, error) {
	if filePath == "" {
		filePath = CredentialsFile
	}

	cm := &CredentialsManager{
		store: &models.CredentialsStore{
			Users: make(map[string]*models.User),
		},
		file: filePath,
	}

	// Load existing credentials
	if err := cm.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}

	return cm, nil
}

// Load reads credentials from file
func (cm *CredentialsManager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := ioutil.ReadFile(cm.file)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, cm.store); err != nil {
		return fmt.Errorf("failed to parse credentials: %w", err)
	}

	fmt.Printf("✅ Loaded %d users from credentials file\n", len(cm.store.Users))
	return nil
}

// saveLocked writes credentials to file (assumes lock is already held)
func (cm *CredentialsManager) saveLocked() error {
	data, err := json.MarshalIndent(cm.store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	if err := ioutil.WriteFile(cm.file, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

// Save writes credentials to file
func (cm *CredentialsManager) Save() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	return cm.saveLocked()
}

// RegisterUser adds a new user
func (cm *CredentialsManager) RegisterUser(id, email, password string) (*models.User, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Check if user already exists
	if _, exists := cm.store.Users[id]; exists {
		return nil, fmt.Errorf("user %s already exists", id)
	}

	user := &models.User{
		ID:        id,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cm.store.Users[id] = user

	// Save to file (using saveLocked since we already hold the write lock)
	if err := cm.saveLocked(); err != nil {
		return nil, err
	}

	fmt.Printf("✅ User registered: %s (%s)\n", id, email)
	return user, nil
}

// GetUser retrieves a user by ID
func (cm *CredentialsManager) GetUser(userID string) (*models.User, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	user, exists := cm.store.Users[userID]
	if !exists {
		return nil, fmt.Errorf("user %s not found", userID)
	}

	return user, nil
}

// ListUsers returns all registered users
func (cm *CredentialsManager) ListUsers() []*models.User {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	users := make([]*models.User, 0, len(cm.store.Users))
	for _, user := range cm.store.Users {
		users = append(users, user)
	}

	return users
}

// JobManager handles job scheduling and state
type JobManager struct {
	mu    sync.RWMutex
	queue *models.JobQueue
	file  string
}

// NewJobManager creates a new job manager
func NewJobManager(filePath string) (*JobManager, error) {
	if filePath == "" {
		filePath = JobsFile
	}

	jm := &JobManager{
		queue: &models.JobQueue{
			Jobs: make(map[string]*models.Job),
		},
		file: filePath,
	}

	// Load existing jobs
	if err := jm.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load jobs: %w", err)
	}

	return jm, nil
}

// Load reads jobs from file
func (jm *JobManager) Load() error {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	data, err := ioutil.ReadFile(jm.file)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, jm.queue); err != nil {
		return fmt.Errorf("failed to parse jobs: %w", err)
	}

	fmt.Printf("✅ Loaded %d jobs from queue file\n", len(jm.queue.Jobs))
	return nil
}

// saveLocked writes jobs to file (assumes lock is already held)
func (jm *JobManager) saveLocked() error {
	data, err := json.MarshalIndent(jm.queue, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal jobs: %w", err)
	}

	if err := ioutil.WriteFile(jm.file, data, 0600); err != nil {
		return fmt.Errorf("failed to write jobs file: %w", err)
	}

	return nil
}

// Save writes jobs to file
func (jm *JobManager) Save() error {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	return jm.saveLocked()
}

// CreateJob creates a new job
func (jm *JobManager) CreateJob(userID, courseSlug string) (*models.Job, error) {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	jobID := fmt.Sprintf("job_%s_%s_%d", userID, courseSlug, time.Now().Unix())

	job := &models.Job{
		ID:        jobID,
		UserID:    userID,
		CourseSlug: courseSlug,
		Status:    "pending",
		CreatedAt: time.Now(),
		Progress:  0,
	}

	jm.queue.Jobs[jobID] = job

	// Save to file (using saveLocked since we already hold the write lock)
	if err := jm.saveLocked(); err != nil {
		return nil, err
	}

	fmt.Printf("✅ Job created: %s for user %s\n", jobID, userID)
	return job, nil
}

// GetJob retrieves a job by ID
func (jm *JobManager) GetJob(jobID string) (*models.Job, error) {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	job, exists := jm.queue.Jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job %s not found", jobID)
	}

	return job, nil
}

// UpdateJobStatus updates job status
func (jm *JobManager) UpdateJobStatus(jobID, status string) error {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	job, exists := jm.queue.Jobs[jobID]
	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	job.Status = status
	job.UpdatedAt = time.Now()

	if status == "running" && job.StartedAt.IsZero() {
		job.StartedAt = time.Now()
	}
	if status == "completed" || status == "failed" {
		job.CompletedAt = time.Now()
	}

	// Save to file (using saveLocked since we already hold the write lock)
	return jm.saveLocked()
}

// ListPendingJobs returns all pending jobs
func (jm *JobManager) ListPendingJobs() []*models.Job {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	var pending []*models.Job
	for _, job := range jm.queue.Jobs {
		if job.Status == "pending" {
			pending = append(pending, job)
		}
	}

	return pending
}

// ListUserJobs returns all jobs for a user
func (jm *JobManager) ListUserJobs(userID string) []*models.Job {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	var userJobs []*models.Job
	for _, job := range jm.queue.Jobs {
		if job.UserID == userID {
			userJobs = append(userJobs, job)
		}
	}

	return userJobs
}
