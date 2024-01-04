package api

import (
	"cloud-stream-queue/queue"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EnqueueJob adds a new job to the queue.
func (a *API) EnqueueJob(c *gin.Context) {
	bytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, httpErr{Code: "invalid_body"})
		return
	}

	jobIn := queue.Job{}
	err = json.Unmarshal(bytes, &jobIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpErr{Code: "invalid_body"})
		return
	}

	id := a.service.EnqueueJob(jobIn)

	var enqueueResp struct {
		ID int `json:"id"`
	}
	enqueueResp.ID = id
	c.JSON(http.StatusOK, enqueueResp)
	return
}

// DeuqueJob dequeues a job from the queue to be completed.
func (a *API) DeuqueJob(c *gin.Context) {
	job, err := a.service.DequeueJob()
	if err != nil {
		if err == queue.ErrNoJobs {
			c.JSON(http.StatusNotFound, httpErr{Code: "no_jobs"})
			return
		}
		c.JSON(http.StatusInternalServerError, httpErr{Code: "internal"})
		return
	}

	var dequeueResp struct {
		Job queue.Job `json:"job"`
	}
	dequeueResp.Job = job
	c.JSON(http.StatusOK, dequeueResp)
	return
}

// ConcludeJob marks a job in the queue concluded.
func (a *API) ConcludeJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// May want to do a 404 to obfuscate what jobs are available.
		c.JSON(http.StatusBadRequest, httpErr{Code: "invalid_id"})
		return
	}

	bytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, httpErr{Code: "invalid_json"})
		return
	}

	var concludeBody struct {
		Result queue.JobResult `json:"result"`
	}

	err = json.Unmarshal(bytes, &concludeBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpErr{Code: "invalid_json"})
		return
	}

	err = a.service.ConcludeJob(id, concludeBody.Result)
	if err != nil {
		if err == queue.ErrJobNotfound {
			c.JSON(http.StatusNotFound, httpErr{Code: "job_not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, httpErr{Code: "internal"})
		return
	}

	c.Status(http.StatusOK)
	return
}

// GetJob returns information about a given job id.
func (a *API) GetJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// May want to do a 404 to obfuscate what jobs are available.
		c.JSON(http.StatusBadRequest, httpErr{Code: "invalid_id"})
		return
	}

	job, err := a.service.GetJob(id)
	if err != nil {
		if err == queue.ErrJobNotfound {
			c.JSON(http.StatusNotFound, httpErr{Code: "job_not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, httpErr{Code: "internal"})
		return
	}
	var getJobResp struct {
		Job queue.Job `json:"job"`
	}
	getJobResp.Job = job
	c.JSON(http.StatusOK, getJobResp)
	return
}

// Debug is used for checking the state of the queue programatically while the service is running.
func (a *API) Debug(c *gin.Context) {
	queueMap, err := a.service.GetQueue()
	if err != nil {
		c.String(http.StatusInternalServerError, "internal error")
		return
	}
	var debugResp struct {
		Queued     []queue.Job `json:"queued"`
		InProgress []queue.Job `json:"in_progress"`
		Concluded  []queue.Job `json:"concluded"`
	}
	// initialize these such that the json has empty arrays instead of nulls.
	debugResp.Queued = []queue.Job{}
	debugResp.InProgress = []queue.Job{}
	debugResp.Concluded = []queue.Job{}

	for _, job := range queueMap {
		if job.Status == queue.JobStatusQueued {
			debugResp.Queued = append(debugResp.Queued, *job)
		} else if job.Status == queue.JobStatusInProgress {
			debugResp.InProgress = append(debugResp.InProgress, *job)
		} else if job.Status == queue.JobStatusConcluded {
			debugResp.Concluded = append(debugResp.Concluded, *job)
		}
	}

	bytes, _ := json.Marshal(debugResp)
	fmt.Println("Debug:", string(bytes))
	c.JSON(http.StatusOK, debugResp)
}
