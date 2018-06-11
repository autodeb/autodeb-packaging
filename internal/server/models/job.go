package models

// JobType is the type of job
type JobType int

// Job Type enum
const (
	JobTypeUnknown JobType = iota
	JobTypeBuild
	JobTypeAutopkgtest
	JobTypeForward
)

func (jt JobType) String() string {
	switch jt {
	case JobTypeBuild:
		return "build"
	case JobTypeAutopkgtest:
		return "autopkgtest"
	case JobTypeForward:
		return "forward"
	default:
		return "unknown"
	}
}

// JobStatus is the status of the job
type JobStatus int

// Job Status enum
const (
	JobStatusUnknown JobStatus = iota
	JobStatusQueued
	JobStatusAssigned
	JobStatusSuccess
	JobStatusFailed
)

func (js JobStatus) String() string {
	switch js {
	case JobStatusQueued:
		return "queued"
	case JobStatusAssigned:
		return "assigned"
	case JobStatusSuccess:
		return "success"
	case JobStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// Job is a builds a test, etc.
type Job struct {
	ID     uint      `json:"id"`
	Type   JobType   `json:"type"`
	Status JobStatus `json:"status"`

	// The upload that has triggered this job.
	// The uploadID is also set to all child jobs.
	UploadID uint `json:"upload_id"`

	// Some job's artifacts serve as input to other jobs.
	// For example: a build job's artifacts is an input to an autopkgtest job
	InputArtifactID uint `json:"input_artifact_id"`
}
