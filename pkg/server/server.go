package server

import (
	"net/http"

	"github.com/grantchen2003/chunky/internal"
	"github.com/grantchen2003/chunky/internal/database"
	"github.com/grantchen2003/chunky/internal/filestorer"
	"github.com/grantchen2003/chunky/internal/handler"
)

type Server struct {
	port              string
	handlerToEndpoint map[string]string
	uploadService     *internal.UploadService
}

func NewServer(port string) (*Server, error) {
	db, err := database.NewPostgresql()
	if err != nil {
		return nil, err
	}

	localFileStore, err := filestorer.NewLocalFileStore("./chunky-local-file-store")
	if err != nil {
		return nil, err
	}

	uploadService := internal.NewUploadService(db, localFileStore)

	return &Server{
		port: port,
		handlerToEndpoint: map[string]string{
			"initiateUploadSession": "/initiateUploadSession",
			"byteRangesToUpload":    "/byteRangesToUpload",
			"uploadFileChunk":       "/uploadFileChunk",
		},
		uploadService: uploadService,
	}, nil
}

func (s *Server) SetInitiateUploadSessionEndpoint(endpoint string) {
	s.handlerToEndpoint["initiateUploadSession"] = endpoint
}

func (s *Server) SetByteRangesToUploadEndpoint(endpoint string) {
	s.handlerToEndpoint["byteRangesToUpload"] = endpoint
}

func (s *Server) SetUploadFileChunkEndpoint(endpoint string) {
	s.handlerToEndpoint["uploadFileChunk"] = endpoint
}

func (s *Server) Start() error {
	initiateUploadSessionHandler := handler.NewInitiateUploadSessionHandler(s.uploadService)
	byteRangesToUploadHandler := handler.NewByteRangesToUploadHandler()
	uploadFileChunkHandler := handler.NewUploadFileChunkHandler(s.uploadService)

	http.HandleFunc(s.handlerToEndpoint["initiateUploadSession"], initiateUploadSessionHandler.Handle)
	http.HandleFunc(s.handlerToEndpoint["byteRangesToUpload"], byteRangesToUploadHandler.Handle)
	http.HandleFunc(s.handlerToEndpoint["uploadFileChunk"], uploadFileChunkHandler.Handle)

	return http.ListenAndServe(s.port, nil)
}
