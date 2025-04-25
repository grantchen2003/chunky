package chunky

import (
	"github.com/grantchen2003/chunky/internal"
	"github.com/grantchen2003/chunky/pkg/client"
	"github.com/grantchen2003/chunky/pkg/server"
)

type Client = client.Client
type UploadEndpoints = internal.UploadEndpoints

var (
	NewClient = client.NewClient
	NewServer = server.NewServer
)
