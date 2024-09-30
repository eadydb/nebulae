package util

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestPortSet(t *testing.T) {
	pf := &PortSet{}

	// Try to store a port
	pf.Set(9000)

	// Try to load the port
	if alreadySet := pf.LoadOrSet(9000); !alreadySet {
		t.Fatal("didn't load port 9000 correctly")
	}

	if alreadySet := pf.LoadOrSet(4000); alreadySet {
		t.Fatal("didn't store port 4000 correctly")
	}

	if alreadySet := pf.LoadOrSet(4000); !alreadySet {
		t.Fatal("didn't load port 4000 correctly")
	}
}

func TestGetAvailablePort(t *testing.T) {
	N := 100

	var (
		ports  PortSet
		lock   sync.Mutex
		wg     sync.WaitGroup
		errors = map[int]error{}
	)

	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			port := GetAvailablePort(Loopback, 4503, &ports)

			l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", Loopback, port))
			if err != nil {
				lock.Lock()
				errors[port] = err
				lock.Unlock()
			} else {
				l.Close()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	for port, err := range errors {
		t.Errorf("available port (%d) couldn't be used: %v", port, err)
	}
}
