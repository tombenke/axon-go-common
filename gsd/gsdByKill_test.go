package gsd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"syscall"
	"testing"
)

func TestRegisterByKill(t *testing.T) {
	var mu sync.Mutex
	gsdCbCalled := false

	wg := sync.WaitGroup{}

	// Register the callback handler
	sigsCh := Register(&wg, func(s os.Signal) {
		mu.Lock()
		gsdCbCalled = true
		mu.Unlock()
	})

	// Sent TERM signal, then wait for termination
	sigsCh <- syscall.SIGTERM
	wg.Wait()

	// Checks if callback was called
	mu.Lock()
	assert.True(t, gsdCbCalled)
	mu.Unlock()
}
