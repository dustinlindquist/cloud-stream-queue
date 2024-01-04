package queue

import (
	"errors"
	"fmt"
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

// JobResult ...
type JobResult string

const (
	// JobResultFailed represents job that failed and should be reprocessed.
	JobResultFailed = "FAILED"
)

// ErrJobNotfound represents when a specific job was asked for but not found.
var ErrJobNotfound = errors.New("job not found")

// ErrNoJobs represents when dequeue is called but there are no jobs in the queue.
var ErrNoJobs = errors.New("no jobs")

// Job represents the metadata of a job.
type Job struct {
	ID       int       `json:"id"`
	Type     JobType   `json:"type"`
	Status   JobStatus `json:"status"`
	Attempts int       `json:"attempts"`
}

// Queue stores jobs available to be dequeued as well as information about completed jobs.
type Queue struct {
	mu                   sync.Mutex
	NotCriticalAvailable []*Job
	CriticalAvailable    []*Job
	JobMap               map[int]*Job
	count                int
}

// New returns an initialized queue.
func New() *Queue {
	return &Queue{
		NotCriticalAvailable: []*Job{},
		CriticalAvailable:    []*Job{},
		JobMap:               make(map[int]*Job),
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

	switch job.Type {
	case JobTypeTimeCritical:
		q.CriticalAvailable = append(q.CriticalAvailable, &job)
	case JobTypeNotTimeCritical:
		q.NotCriticalAvailable = append(q.NotCriticalAvailable, &job)
	default:
		fmt.Println("Invalid Jop type:", job.Type)
	}

	q.JobMap[q.count] = &job
	return q.count
}

// DequeueJob returns the oldest item in the queue that's status is still queued, ignoring type.
// Returning an error if no jobs are available.
func (q *Queue) DequeueJob() (Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.CriticalAvailable) < 1 && len(q.NotCriticalAvailable) < 1 {
		return Job{}, ErrNoJobs
	}

	var job *Job
	if len(q.CriticalAvailable) > 0 {
		job = q.CriticalAvailable[0]
		q.CriticalAvailable = q.CriticalAvailable[1:]
	} else if len(q.NotCriticalAvailable) > 0 {
		job = q.NotCriticalAvailable[0]
		q.NotCriticalAvailable = q.NotCriticalAvailable[1:]
	}

	q.JobMap[job.ID].Status = JobStatusInProgress

	return *job, nil
}

func (q *Queue) requeueFailed(job Job) error {

	// increment the count prior to using it as our id.
	// q.count++
	// job.ID = q.count
	job.Status = JobStatusQueued

	switch job.Type {
	case JobTypeTimeCritical:
		q.CriticalAvailable = append(q.CriticalAvailable, &job)
	case JobTypeNotTimeCritical:
		q.NotCriticalAvailable = append(q.NotCriticalAvailable, &job)
	default:
		fmt.Println("Invalid Jop type:", job.Type)
	}

	job.Status = JobStatusQueued
	q.JobMap[job.ID] = &job

}

// ConcludeJob marks a job concluded.
func (q *Queue) ConcludeJob(id int, result JobResult) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	job, ok := q.JobMap[id]
	if !ok {
		return ErrJobNotfound
	}

	if result == JobResultFailed {

		fmt.Println("hit")
		// q.JobMap[id].Status = JobStatusConcluded
		q.EnqueueJob(*job)
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
