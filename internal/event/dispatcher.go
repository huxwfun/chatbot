package event

import (
	"fmt"
	"log"
)

type Dispatcher struct {
	jobs     chan job
	channels map[string][]Listener
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		jobs:     make(chan job, 10),
		channels: make(map[string][]Listener),
	}

	go d.consume()

	return d
}

func (d *Dispatcher) Register(listener Listener, names ...string) error {
	for _, name := range names {
		if _, ok := d.channels[name]; !ok {
			d.channels[name] = make([]Listener, 0)
		}
		d.channels[name] = append(d.channels[name], listener)
		log.Printf("channel %s registered %d listeners", name, len(d.channels[name]))
	}

	return nil
}

func (d *Dispatcher) Dispatch(channel string, event interface{}) error {
	if _, ok := d.channels[channel]; !ok {
		return fmt.Errorf("the '%s' channel is not registered", channel)
	}
	d.jobs <- job{channel: channel, event: event}

	return nil
}

func (d *Dispatcher) consume() {
	for job := range d.jobs {
		for _, listener := range d.channels[job.channel] {
			listener(job.event)
		}
	}
}
