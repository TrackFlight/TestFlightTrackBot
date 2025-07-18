package utils

type Pool struct {
	jobs  int
	queue chan struct{}
}

func NewPool(size int) *Pool {
	if size < 0 {
		size = 0
	}
	return &Pool{
		jobs:  size,
		queue: make(chan struct{}, size),
	}
}

func (p *Pool) Enqueue(job func()) {
	if p.jobs == 0 {
		go job()
		return
	}
	p.queue <- struct{}{}
	go func() {
		defer func() {
			<-p.queue
		}()
		job()
	}()
}
