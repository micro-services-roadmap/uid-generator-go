package worker

import (
	"github.com/micro-services-roadmap/cloudflare/kv/worker"
)

type CloudflareAssigner struct {
}

func (c CloudflareAssigner) AssignWorkerId() int64 {
	for i := 0; i < 10; i++ {
		if id, err := worker.NextWorkerID(); err != nil {
			continue
		} else {
			return id
		}
	}
	panic("Could not assign worker id")
}
