/*

the worker pool is goroutine pool dedicated for cutting clips from source video, it is not a general job queue, the general queue will dispatch the clip job payload, and this worker will execute it

-> workflow will be something like

1. worker with service and worker pool size (limited, to stay away from goroutine leaks)

2. then we will launch the pool size goroutine

3. will use till the buffer capacity, not more than that

4. at last we will drain (will take each and every stuff out of channel and clean the channel)

-> will use simple for range iteration to implement draining

*/

package clipping

import "sync"

// worker struct consist of service, poolsize, channels (clip and err), and waitgroup

// clipChan will have clip id's to process

// errChan contains one entry per failed clip

type Worker struct {
	service  *Service
	poolsize int
	clipChan chan string
	errChan  chan error
	wg       sync.WaitGroup
}

// new worker function will have worker struct with the buffered channel double of actual pool size

func NewWorker(service *Service, poolsize int) *Worker {
	if poolsize <= 0 {
		poolsize = 2
	}

	return &Worker{
		service:  service,
		poolsize: poolsize,
		clipChan: make(chan string, poolsize*2),
		errChan:  make(chan error, poolsize*2),
	}
}
