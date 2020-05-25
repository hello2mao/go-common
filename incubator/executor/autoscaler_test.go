package executor
//
//import (
//	"testing"
//	"time"
//)
//
//var (
//	testScaleUpThreshold   = time.Millisecond
//	testScaleDownThreshold = time.Second * 10
//)
//
//type testJob struct {
//	ch chan struct{}
//}
//
//func (job *testJob) Exec() {
//	<-job.ch
//}
//
//func TestAutoscaler_ActiveWorkerCount(t *testing.T) {
//	testCases := []struct {
//		Name           string
//		Param1         int
//		Param2         int
//		CreateJobCount int
//		PrepareJobs    func(int) []*testJob
//		FinishJobs     func([]*testJob)
//		Expected       int
//	}{
//		{
//			Name:           "test min worker",
//			Param1:         5,
//			Param2:         10,
//			CreateJobCount: 3,
//			PrepareJobs: func(count int) []*testJob {
//				jobs := make([]*testJob, count)
//				for i := range jobs {
//					jobs[i] = &testJob{
//						ch: make(chan struct{}),
//					}
//				}
//				return jobs
//			},
//			FinishJobs: func(jobs []*testJob) {
//				for _, job := range jobs {
//					close(job.ch)
//				}
//			},
//			Expected: 5,
//		},
//		{
//			Name:           "test max worker",
//			Param1:         5,
//			Param2:         10,
//			CreateJobCount: 15,
//			PrepareJobs: func(count int) []*testJob {
//				jobs := make([]*testJob, count)
//				for i := range jobs {
//					jobs[i] = &testJob{
//						ch: make(chan struct{}),
//					}
//				}
//				return jobs
//			},
//			FinishJobs: func(jobs []*testJob) {
//				for _, job := range jobs {
//					close(job.ch)
//				}
//			},
//			Expected: 10,
//		},
//	}
//
//	for _, c := range testCases {
//		func() {
//			scaler := NewAutoscaler(c.Param1, c.Param2, testScaleUpThreshold, testScaleDownThreshold)
//			scaler.Run()
//			jobs := c.PrepareJobs(c.CreateJobCount)
//			for _, job := range jobs {
//				go scaler.AddJobBlocked(job)
//			}
//			time.Sleep(time.Second)
//			if workerCount := scaler.ActiveWorkerCount(); workerCount != c.Expected {
//				t.Errorf("name: %s, expected: %d, actual: %d", c.Name, c.Expected, workerCount)
//			}
//			c.FinishJobs(jobs)
//			scaler.Stop()
//		}()
//	}
//}
//
//func TestAutoscaler_WaitForAllJobsDone(t *testing.T) {
//	testCases := []struct {
//		Name            string
//		CreateJobCount  int
//		PrepareJobs     func(int) []*testJob
//		FinishJobs      func([]*testJob)
//		Delay           time.Duration
//		WaitTimeout     time.Duration
//		ExpectedTimeout bool
//	}{
//		{
//			Name:           "test wait success",
//			CreateJobCount: 1,
//			PrepareJobs: func(count int) []*testJob {
//				jobs := make([]*testJob, count)
//				for i := range jobs {
//					jobs[i] = &testJob{
//						ch: make(chan struct{}),
//					}
//				}
//				return jobs
//			},
//			FinishJobs: func(jobs []*testJob) {
//				for _, job := range jobs {
//					close(job.ch)
//				}
//			},
//			Delay:           time.Second,
//			WaitTimeout:     time.Second * 2,
//			ExpectedTimeout: false,
//		},
//		{
//			Name:           "test wait timeout",
//			CreateJobCount: 1,
//			PrepareJobs: func(count int) []*testJob {
//				jobs := make([]*testJob, count)
//				for i := range jobs {
//					jobs[i] = &testJob{
//						ch: make(chan struct{}),
//					}
//				}
//				return jobs
//			},
//			FinishJobs: func(jobs []*testJob) {
//				for _, job := range jobs {
//					close(job.ch)
//				}
//			},
//			Delay:           time.Second * 2,
//			WaitTimeout:     time.Second,
//			ExpectedTimeout: true,
//		},
//		{
//			Name:           "test no job",
//			CreateJobCount: 0,
//			PrepareJobs: func(count int) []*testJob {
//				jobs := make([]*testJob, count)
//				for i := range jobs {
//					jobs[i] = &testJob{
//						ch: make(chan struct{}),
//					}
//				}
//				return jobs
//			},
//			FinishJobs: func(jobs []*testJob) {
//				for _, job := range jobs {
//					close(job.ch)
//				}
//			},
//			Delay:           time.Second,
//			WaitTimeout:     time.Second * 2,
//			ExpectedTimeout: false,
//		},
//	}
//
//	for _, c := range testCases {
//		scaler := NewAutoscaler(5, 5, testScaleUpThreshold, testScaleDownThreshold)
//		scaler.Run()
//		jobs := c.PrepareJobs(c.CreateJobCount)
//		for _, job := range jobs {
//			go scaler.AddJobBlocked(job)
//		}
//
//		time.Sleep(time.Second)
//
//		ch := make(chan struct{})
//		timer := time.NewTimer(c.WaitTimeout)
//		go func() {
//			scaler.WaitForAllJobsDone()
//			close(ch)
//		}()
//
//		go func() {
//			time.Sleep(c.Delay)
//			c.FinishJobs(jobs)
//		}()
//
//		select {
//		case <-ch:
//			if c.ExpectedTimeout {
//				t.Errorf("name: %s, expected timeout, actual finish", c.Name)
//			}
//		case <-timer.C:
//			if !c.ExpectedTimeout {
//				t.Errorf("name: %s, expected finish, actual timeout", c.Name)
//			}
//		}
//
//		scaler.Stop()
//	}
//}
