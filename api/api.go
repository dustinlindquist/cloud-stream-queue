package api

import (
	"cloud-stream-queue/queue"

	"github.com/gin-gonic/gin"
)

type service interface {
	EnqueueJob(job queue.Job) int
	DequeueJob() (queue.Job, error)
	ConcludeJob(id int, result queue.JobResult) error
	GetJob(id int) (queue.Job, error)
	GetQueue() (map[int]*queue.Job, error)
}

// API holds the router and dependencies for fulfilling API calls.
type API struct {
	service service
	// This needs to be exported so we can "Run" the router from main. We may also want to expose a Start(portNum)
	// method on this API that is exported.
	Router *gin.Engine
}

// New returns a configured API with its needed dependencies and a router with routes configured.
func New(service service, router *gin.Engine) API {
	api := API{
		service: service,
	}

	router.POST("/jobs/enqueue", api.EnqueueJob)
	router.GET("/jobs/dequeue", api.DeuqueJob)
	router.PATCH("/jobs/:id/conclude", api.ConcludeJob) // I've used PATCH because we're modifying 1 field of the job in our storage.
	router.GET("/jobs/:id", api.GetJob)
	// The /debug route is commented out and should be removed prior to prod release
	router.GET("/jobs/debug", api.Debug)

	api.Router = router
	return api
}

type httpErr struct {
	Code string `json:"code"`
}
