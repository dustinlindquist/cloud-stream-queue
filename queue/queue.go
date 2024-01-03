package queue

import (
	"errors"
	"sync"
)

// JobStatus is a typed string representing valid job states.
type JobStatus string

const (
	// JobStatusQueued represents job available in the queue to be dequeued.
	JobStatusQueued = "QUEUED"
	// JobStatusInProgress represents a job that has been dequeued as is currently be worked on.
	JobStatusInProgress = "IN_PROGRESS"
	// JobStatusConcluded represents a job that has been completed.
	JobStatusConcluded = "CONCLUDED"
)

// JobType is a typed string representing valid job types.
type JobType string

const (
	// JobTypeTimeCritical represents job that is time critical.
	JobTypeTimeCritical = "TIME_CRITICAL"

	// JobTypeNotTimeCritical represents job that is not time critical.
	JobTypeNotTimeCritical = "NOT_TIME_CRITICAL"
)

// ErrJobNotfound represents when a specific job was asked for but not found.
var ErrJobNotfound = errors.New("job not found")

// ErrNoJobs represents when dequeue is called but there are no jobs in the queue.
var ErrNoJobs = errors.New("no jobs")

// Job represents the metadata of a job.
type Job struct {
	ID     int       `json:"id"`
	Type   JobType   `json:"type"`
	Status JobStatus `json:"status"`
}

// Queue stores jobs available to be dequeued as well as information about completed jobs.
type Queue struct {
	mu        sync.Mutex
	Available []*Job
	JobMap    map[int]*Job
	count     int
}

// New returns an initialized queue.
func New() *Queue {
	return &Queue{
		Available: []*Job{},
		JobMap:    make(map[int]*Job),
	}
}

// EnqueueJob add a new job to the queue storing it in both the map of all jobs and slice of available jobs.
func (q *Queue) EnqueueJob(job Job) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	// increment the count prior to using it as our id.
	q.count++
	job.ID = q.count
	job.Status = JobStatusQueued
	q.Available = append(q.Available, &job)

	q.JobMap[q.count] = &job
	return q.count
}

// DequeueJob returns the oldest item in the queue that's status is still queued, ignoring type.
// Returning an error if no jobs are available.
func (q *Queue) DequeueJob() (Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.Available) < 1 {
		return Job{}, ErrNoJobs
	}

	job := q.Available[0]
	q.JobMap[job.ID].Status = JobStatusInProgress
	q.Available = q.Available[1:]

	return *job, nil
}

// ConcludeJob marks a job concluded.
func (q *Queue) ConcludeJob(id int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	_, ok := q.JobMap[id]
	if !ok {
		return ErrJobNotfound
	}

	q.JobMap[id].Status = JobStatusConcluded
	return nil
}

// GetJob returns a job from our queue regardless of its status. Returning an error if not found.
func (q *Queue) GetJob(id int) (Job, error) {
	job, ok := q.JobMap[id]
	if !ok {
		return Job{}, ErrJobNotfound
	}
	return *job, nil
}

// GetQueue is used to expose the state of the queue for debugging.
func (q *Queue) GetQueue() (map[int]*Job, error) {
	return q.JobMap, nil
}
