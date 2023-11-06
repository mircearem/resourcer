package main

import (
	"context"
	"time"

	"github.com/mircearem/resourcer/rh"
)

// Server type
type Server struct {
	rh         *rh.ResourceHandler
	interval   time.Duration
	mQuitCh    chan struct{}
	toggle     []chan struct{}
	quitCh     []chan struct{}
	intervalCh []chan struct{}
}

// Function to create a new server
func NewServer(ctx context.Context) *Server {
	r := rh.NewHandler(ctx)

	return &Server{
		rh:         r,
		interval:   time.Second,
		toggle:     make([]chan struct{}, 0),
		intervalCh: make([]chan struct{}, 0),
		quitCh:     make([]chan struct{}, 0),
		mQuitCh:    make(chan struct{}),
	}
}

// func (s *Server) updateCpuUsage(quit, itv, toggle chan struct{}) {
// 	go func() {
// 		running := true
// 		ticker := time.NewTicker(s.interval)
// 		for {
// 			select {
// 			// Update the resource usage variable
// 			case <-ticker.C:
// 				s.r.getCpuUsage()
// 			// Update the interval
// 			case <-itv:
// 				ticker.Reset(s.interval)
// 			// Toggle the service on or off
// 			case <-toggle:
// 				if running {
// 					running = false
// 					ticker.Stop()
// 				} else {
// 					running = true
// 					ticker.Reset(s.interval)
// 				}
// 			// Shutdown signal
// 			case <-quit:
// 				close(quit)
// 				close(itv)
// 				close(toggle)
// 				return
// 			}
// 		}
// 	}()
// }

// Run the server
// func (s *Server) Run(i chan time.Duration, toggle chan struct{}) {
// 	// Create all the signal channels for each service
// 	for i := 0; i < 3; i++ {
// 		ch := make(chan struct{})
// 		it := make(chan struct{})
// 		tog := make(chan struct{})

// 		s.quitCh = append(s.quitCh, ch)
// 		s.intervalCh = append(s.intervalCh, it)
// 		s.toggle = append(s.toggle, tog)
// 	}

// 	// Start the individual services, maybe also add toggle chan
// 	s.updateCpuUsage(s.quitCh[0], s.intervalCh[0], s.toggle[0])
// 	s.updateRamUsage(s.quitCh[1], s.intervalCh[1], s.toggle[1])
// 	s.updateHddUsage(s.quitCh[2], s.intervalCh[2], s.toggle[2])

// loop:
// 	for {
// 		select {
// 		case itv := <-i:
// 			// Update the interval for each service
// 			s.interval = itv

// 			// Send signal to all child go routines to update the interval
// 			for i := 0; i < 3; i++ {
// 				s.intervalCh[i] <- struct{}{}
// 			}
// 		case <-toggle:
// 			// Toggle each serice on or off
// 			for i := 0; i < 3; i++ {
// 				s.toggle[i] <- struct{}{}
// 			}

// 		case <-s.mQuitCh:
// 			// Send stop signal to all go routines
// 			for i := 0; i < 3; i++ {
// 				// Send shutdown signals to go routines
// 				s.quitCh[i] <- struct{}{}
// 			}

// 			log.Println("Stopping service")
// 			break loop
// 		}
// 	}
// }

// // Shut the server down
// func (s *Server) Shutdown() {
// 	// Send a quit signal on all of the channels in the quit array
// 	s.mQuitCh <- struct{}{}
// 	close(s.mQuitCh)
// }
