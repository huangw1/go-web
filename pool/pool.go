package pool

import "sync"

type Pool struct {
	ch chan struct{}
	wg *sync.WaitGroup
}

func New(poolSize, concurrent int) *Pool {
	pool := &Pool{
		ch: make(chan struct{}, poolSize),
		wg: new(sync.WaitGroup),
	}
	if concurrent > 0 {
		pool.wg.Add(concurrent)
	}
	return pool
}

func (p *Pool)Submit(task func()) {
	p.ch <- struct{}{}
	go func() {
		defer func() {
			p.wg.Done()
			<- p.ch
		}()
		task()
	}()
}

func (p *Pool) Wait()  {
	p.wg.Wait()
}