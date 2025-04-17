package chunky

import "github.com/grantchen2003/chunky/pkg/client"

type Client = client.Client

var (
	NewClient                     = client.NewClient
	UploadCompleted               = client.UploadCompleted
	UploadFailed                  = client.UploadFailed
	UploadPaused                  = client.UploadPaused
	ErrResumedOnOngoingUpload     = client.ErrResumedOnOngoingUpload
	ErrResumedOnNonExistingUpload = client.ErrResumedOnNonExistingUpload
)
