package executor

// Job is the interface of the REAL JOB need to impl
type Job interface {
	// Job exec
	Exec()
}
