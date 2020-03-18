package executor

// Executor backend interface
type Executor interface {
	Run()
	Stop()
	WaitForAllJobsDone()
	AddJob(job Job)
}
