package collector

import (
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

type Metric interface {
	Get(*sync.WaitGroup) error
	Average(*sync.WaitGroup, *api.Stats, *[]Snapshot, int)
}
