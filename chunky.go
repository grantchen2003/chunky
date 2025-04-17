package chunky

import (
	"github.com/grantchen2003/chunky/internal"
	"github.com/grantchen2003/chunky/pkg/client"
)

type Client = client.Client

var (
	NewClient                     = client.NewClient
	UploadCompleted               = internal.UploadCompleted
	UploadFailed                  = internal.UploadFailed
	UploadPaused                  = internal.UploadPaused
	ErrResumedOnOngoingUpload     = internal.ErrResumedOnOngoingUpload
	ErrResumedOnNonExistingUpload = internal.ErrResumedOnNonExistingUpload
)
