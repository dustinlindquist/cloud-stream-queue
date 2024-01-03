## Stream Job Queue
This is a simple service that provides a thread safe job queue; allowing multiple callers (expected to be consumers/producers) to queue up jobs to be completed at a later date and then marking those jobs concluded. The job queue stores all jobs, allowing for fetching the information of concluded jobs, in progress jobs, and yet to be deqeued jobs. 

## Testing
To run all of the tests run with the race detector run this command from the top level.
```
 go test ./... -race
```
To run a specific test here is an example:
```
go test ./... -race -run TestEnqueueJob
```

## Starting the service
Run the service by using this command from the top level:
```
go run main.go
```

## Notes
1. Wouldn't use an interger for the identifier, but the spec called for it
2. I don't hold the connection open on dequeue, waiting for more jobs to be queued, we return a 404 with no jobs available as an error. I've typically worked in IO bound systems so returning quickly is kind of my default, though I see the use for it to hold that connection open depending on the consumers. Returning an error feels heavy handed but a 404 is easy to switch on an allow consumers to do with that what they want. 
3. Does the Get /jobs/{job_id} need to return information about completed jobs? This drastically changes what type of datastructure should be used for the queue. I assumed that one would want to be able to fetch all jobs (including completed jobs) This results in the map of jobs that will grow unbounded in memory. If this constraint isn't needed I would not use a map and slice instead opt for a linked list or ring buffer.
4. Should we deny requests to conclude a job that hasn't been dequeued? Should there be validation that it's the same CONSUMER_ID? or have any security that prevents the modifying of the queue. 

## Things to add
1. Logging (hopefully we would have some library/middleware that logs http requests/responses)
2. Metrics (hopefully we would have some library/middleware that emits http requests/responses and latency)
3. Configurable Timeout
4. Configurable ratelimit / backoff
5. Dockerfile
6. Add a MaxSize to the queue, as it currently stores all completed jobs and will grow unbounded
7. Use UUID instead of int and remove the counter
8. Cleanup error handling to smooth the conversion from queue/errs.go to http errors.
9. Generate mocks for the service interface in the api package and write unit tests with expectations and assertions on the response body and http status codes of our routes
10. .env file to hold varibles for the service, like GIN_RELASE_MODE and PORT
11. OpenAPI Spec (swagger) to document the api spec
12. A graceful shutdown that respects our timeout allowing requests in flight to complete prior to the service shutting down.