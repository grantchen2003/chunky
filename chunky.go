package chunky

import (
	"github.com/grantchen2003/chunky/internal"
	"github.com/grantchen2003/chunky/pkg/client"
)

type Client = client.Client

var (
	NewClient                     = client.NewClient
	ErrResumedOnOngoingUpload     = internal.ErrResumedOnOngoingUpload
	ErrResumedOnNonExistingUpload = internal.ErrResumedOnNonExistingUpload
)
