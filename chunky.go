package chunky

import (
	"github.com/grantchen2003/chunky/internal/upload"
	"github.com/grantchen2003/chunky/pkg/client"
	"github.com/grantchen2003/chunky/pkg/server"
)

type Client = client.Client
type UploadEndpoints = upload.Endpoints

var (
	NewClient = client.NewClient
	NewServer = server.NewServer
)
