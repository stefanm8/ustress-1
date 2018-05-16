package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/montanaflynn/stats"
	"golang.org/x/net/context"
)

// A message processes url and returns the result on responseChan.
// ctx is places in a struct, but this is ok to do.

var (
	ctx context.Context
)

// the structure that will be passed to channels
type message struct {
	responseChan chan<- *message
	worker       Worker
	ctx          context.Context
}

// Worker details, needed for returning the output and build the report
type Worker struct {
	Request  int     `json:"request"`
	Status   int     `json:"status"`
	Thread   int     `json:"thread"`
	url      string  // should use net.url
	Duration float64 `json:"duration"`
}

// Report is the report structure, object
// @todo calculate percentile 99, 95, 75, 50
type Report struct {
	// id uuid
	// timestamp
	URL       string    `json:"url"`
	TimeStamp time.Time `json:"timestamp"`
	UUID      uuid.UUID `json:"uuid"`
	Stats     struct {
		Median      float64 `json:"median"`
		PercentileA float64 `json:"50_percentile"`
		PercentileB float64 `json:"75_percentile"`
		PercentileC float64 `json:"95_percentile"`
		PercentileD float64 `json:"99_percentile"`
	} `json:"stats"`

	Duration float64   `json:"durationTotal"`
	Workers  []*Worker `json:"data"`
}

func processMessages(id int, work <-chan *message) {
	for job := range work {
		select {
		// If the context is finished, don't bother processing the
		// message.
		case <-job.ctx.Done():
			continue
		default:
		}

		job.worker.doWork(id)

		select {
		case <-job.ctx.Done():
		case job.responseChan <- job:
		}
	}
}

// doWork method for the worker
func (wrk *Worker) doWork(id int) *Worker {
	start := time.Now()
	wrk.Thread = id
	// define a timeout for the http client
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	//httpResponse, error := client.Head(url)
	httpResponse, error := client.Get(wrk.url)
	if error != nil {
		wrk.Status = 503 // status in case of timeout
	} else {
		wrk.Status = httpResponse.StatusCode
	}

	wrk.Duration = time.Since(start).Seconds()
	log.Printf("Worker Reporting: %+v", *wrk)
	return wrk
}

func newRequest(ctx context.Context, worker Worker, q chan<- *message, report *Report) {
	r := make(chan *message)
	select {
	// If the context finishes before we can send msg onto q,
	// exit early
	case <-ctx.Done():
		fmt.Println("Context ended before q could see message")
		return
	case q <- &message{
		responseChan: r,
		worker:       worker,
		// We are placing a context in a struct.  This is ok since it
		// is only stored as a passed message and we want q to know
		// when it can discard this message
		ctx: ctx,
	}:
	}

	select {
	case out := <-r:
		// fmt.Printf("%v\n", out)
		report.Workers = append(report.Workers, &out.worker)
	// If the context finishes before we could get the result, exit early
	case <-ctx.Done():
		fmt.Println("Context ended before q could process message")
	}
}

func (rep *Report) calcStats() *Report {
	var requestDurations []float64
	for _, value := range rep.Workers {
		requestDurations = append(requestDurations, value.Duration)
	}
	rep.Stats.PercentileA, _ = stats.Percentile(requestDurations, 50)
	rep.Stats.PercentileB, _ = stats.Percentile(requestDurations, 75)
	rep.Stats.PercentileC, _ = stats.Percentile(requestDurations, 95)
	rep.Stats.PercentileD, _ = stats.Percentile(requestDurations, 99)
	rep.Stats.Median, _ = stats.Median(requestDurations)
	return rep
}

// NewURLStressReport probes an endpoint and generates a new report
func NewURLStressReport(url string, requests, threads int) ([]byte, error) {
	start := time.Now()
	report := Report{URL: url, TimeStamp: time.Now(), UUID: uuid.New()}

	q := make(chan *message, threads)
	// start number of threads
	for i := 1; i <= threads; i++ {
		go processMessages(i, q)
	}

	// send requests to q
	for k := 1; k <= requests; k++ {
		ctx := context.Background()
		wrk := Worker{url: url, Request: k}
		newRequest(ctx, wrk, q, &report)
	}
	close(q)
	report.calcStats()
	report.Duration = time.Since(start).Seconds()
	b, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}
	// os.Stdout.Write(b)
	return b, nil
}