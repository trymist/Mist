// context manager for different deployments
// used for canceling an ongoing deployment

package queue

import (
	"context"
	"sync"
)

var activeDeployments = make(map[int64]context.CancelFunc)
var mu sync.Mutex

func Register(id int64, cancel context.CancelFunc) {
	mu.Lock()
	defer mu.Unlock()
	activeDeployments[id] = cancel
}

func Unregister(id int64) {
	mu.Lock()
	defer mu.Unlock()
	delete(activeDeployments, id)
}

func Cancel(id int64) bool {
	mu.Lock()
	defer mu.Unlock()
	if cancel, exists := activeDeployments[id]; exists {
		cancel()
		delete(activeDeployments, id)
		return true
	}
	return false
}
