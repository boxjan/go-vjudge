package crawler

import "strings"

type RemoteStatus int

const (
	StatusStart   RemoteStatus = -6 + iota
	StatusPending
	StatusSubmitFailed
	StatusQueuing
	StatusCompiling
	StatusRunning
	StatusAccept
	StatusPresentationError
	StatusWrongAnswer
	StatusTimeLimitExceeded
	StatusMemoryLimitExceeded
	StatusOutputLimitExceeded
	StatusRuntimeError
	StatusCompileError
	StatusSubmitError
	StatusSystemError
	StatusFailedOther
)


var RemoteStatusMap = map[string]RemoteStatus {
	"Wait": StatusQueuing,
	"Pend": StatusQueuing,
	"Queuing": StatusQueuing,
	"Queue": StatusQueuing,

	"Compiling": StatusCompiling,

	"Running": StatusRunning,
	"ing": StatusRunning,

	"Accepted": StatusAccept,

	"Presentation Error": StatusPresentationError,
	"Format Error": StatusPresentationError,

	"Wrong Answer":StatusWrongAnswer,

	"Time Limit Exceed": StatusTimeLimitExceeded,

	"Memory Limit Exceed": StatusMemoryLimitExceeded,

	"Output Limit Exceed": StatusOutputLimitExceeded,

	"Runtime Error": StatusRuntimeError,
	"Segmentation Fault": StatusRuntimeError,
	"Floating Point Error": StatusRuntimeError,
	"Crash": StatusRuntimeError,

	"Compile Error": StatusCompileError,
	"Compilation Error": StatusCompileError,

}

func RemoteStatusType(rawStatus string) RemoteStatus {

	for k, v := range RemoteStatusMap {
		if strings.Contains(rawStatus, k) {
			return v
		}
	}

	return StatusFailedOther
}

