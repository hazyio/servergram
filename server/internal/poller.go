package internal

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Polls struct {
	mu    sync.Mutex
	polls []*Poller
}

func (p *Polls) AddPoll(pl *Poller) *Poller {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.polls = append(p.polls, pl)
	return pl
}
func (p *Polls) StopAll() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, pl := range p.polls {
		pl.Stop()
	}
}

type Poller struct {
	interval time.Duration
	work     func(context.Context)

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

func NewPoller(interval time.Duration, work func(context.Context)) *Poller {
	return &Poller{interval: interval, work: work}
}

func (p *Poller) Start(ctx context.Context) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cancel != nil {
		return // already running
	}

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.done = make(chan struct{})

	go func() {
		defer close(p.done)
		t := time.NewTicker(p.interval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				p.safeRun(ctx)
			}
		}
	}()
}

func (p *Poller) Stop() {
	p.mu.Lock()
	cancel, done := p.cancel, p.done
	p.cancel, p.done = nil, nil
	p.mu.Unlock()

	if cancel == nil {
		return
	}
	cancel()
	<-done // wait for the loop to exit
}

func (p *Poller) safeRun(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("poller panic", "error", r)
		}
	}()
	p.work(ctx)
}
