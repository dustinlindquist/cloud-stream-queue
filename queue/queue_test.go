package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnqueueJob(t *testing.T) {
	type test struct {
		name       string
		jobsIn     []Job
		wantJobMap map[int]*Job
	}

	tests := []test{
		{
			name: "enqueue_first_job_time_critical",
			jobsIn: []Job{
				{Type: JobTypeTimeCritical},
			},
			wantJobMap: map[int]*Job{
				1: {
					ID:     1,
					Status: JobStatusQueued,
					Type:   JobTypeTimeCritical,
				},
			},
		},
		{
			name: "enqueue_first_job_not_time_critical",
			jobsIn: []Job{
				{Type: JobTypeNotTimeCritical},
			},
			wantJobMap: map[int]*Job{
				1: {
					ID:     1,
					Status: JobStatusQueued,
					Type:   JobTypeNotTimeCritical,
				},
			},
		},
		{
			name: "enqueue_multiple_jobs",
			jobsIn: []Job{
				{Type: JobTypeNotTimeCritical},
				{Type: JobTypeTimeCritical},
				{Type: JobTypeNotTimeCritical},
			},
			wantJobMap: map[int]*Job{
				1: {
					ID:     1,
					Status: JobStatusQueued,
					Type:   JobTypeNotTimeCritical,
				},
				2: {
					ID:     2,
					Status: JobStatusQueued,
					Type:   JobTypeTimeCritical,
				},
				3: {
					ID:     3,
					Status: JobStatusQueued,
					Type:   JobTypeNotTimeCritical,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q := New()
			for _, job := range tc.jobsIn {
				q.EnqueueJob(job)
			}

			assert.Equal(t, tc.wantJobMap, q.JobMap)
		})
	}
}
