package executor

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultMinWorkerQuantity is the default minWorkerQuantity
	DefaultMinWorkerQuantity = 10
	// DefaultMaxWorkerQuantity is the default maxWorkerQuantity
	DefaultMaxWorkerQuantity = 20
	// DefaultScaleUpThreshold is the default scaleUpThreshold
	DefaultScaleUpThreshold = 1 * time.Second
	// DefaultScaleDownThreshold is the default scaleDownThreshold
	DefaultScaleDownThreshold = 1 * time.Minute
)

// Autoscaler struct
type Autoscaler struct {
	scaleUpThreshold   time.Duration
	scaleDownThreshold time.Duration
	minWorkerQuantity  int
	maxWorkerQuantity  int

	stopChan        chan struct{}
	jobChan         chan Job
	workerTokenChan chan bool

	wg    *sync.WaitGroup
	jobWg *sync.WaitGroup
}

// NewAutoscaler return new Autoscaler object
func NewAutoscaler(minWorkerQuantity, maxWorkerQuantity int, scaleUpThreshold, scaleDownThreshold time.Duration) *Autoscaler {

	s := Autoscaler{
		minWorkerQuantity:  minWorkerQuantity,
		maxWorkerQuantity:  maxWorkerQuantity,
		scaleUpThreshold:   scaleUpThreshold,
		scaleDownThreshold: scaleDownThreshold,
		stopChan:           make(chan struct{}),
		jobChan:            make(chan Job),
		workerTokenChan:    make(chan bool, maxWorkerQuantity),
		wg:                 &sync.WaitGroup{},
		jobWg:              &sync.WaitGroup{},
	}

	return &s
}

// Run run Autoscaler worker
func (s *Autoscaler) Run() {
	for i := 0; i < s.maxWorkerQuantity; i++ {
		if i < s.minWorkerQuantity {
			s.workerTokenChan <- false // token=false means can not scale down
			s.scaleUp()
		} else {
			s.workerTokenChan <- true // token=true means can scale down
		}
	}
}

// Stop stop all Autoscaler worker
func (s *Autoscaler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}

// WaitForAllJobsDone wait for all job submitted done
func (s *Autoscaler) WaitForAllJobsDone() {
	s.jobWg.Wait()
}

// AddJobBlocked add new job to autoscaler
func (s *Autoscaler) AddJobBlocked(job Job) error {
	err := WriteChanWithTimeout(s.jobChan, job, s.scaleUpThreshold)
	if err != nil {
		success := s.scaleUp()
		if success {
			s.jobChan <- job
			log.Debugf("autoscaler scale up, current count: %d", s.ActiveWorkerCount())
		} else {
			return fmt.Errorf("heavy load, try again later")
		}
	}
	s.jobWg.Add(1)
	return nil
}

// ActiveWorkerCount return the count of active worker
func (s *Autoscaler) ActiveWorkerCount() int {
	return s.maxWorkerQuantity - len(s.workerTokenChan)
}

func (s *Autoscaler) worker(token bool) {
	defer s.wg.Done()

	timer := time.NewTimer(s.scaleDownThreshold)
	defer timer.Stop()
	for {
		timer.Reset(s.scaleDownThreshold)
		select {
		case <-s.stopChan:
			return
		case <-timer.C:
			if token && s.scaleDown(token) {
				return
			}
		case job := <-s.jobChan:
			job.Exec()
			s.jobWg.Done()
		}
	}
}

// scaleUp may return false
func (s *Autoscaler) scaleUp() bool {
	select {
	case token := <-s.workerTokenChan:
		s.wg.Add(1)
		go s.worker(token)
	default:
		return false
	}
	return true
}

// scaleDown always return true
func (s *Autoscaler) scaleDown(token bool) bool {
	select {
	case s.workerTokenChan <- token:
		log.Debugf("autoscaler scale down, current count: %d", s.ActiveWorkerCount())
		return true
	default:
	}

	return false
}
