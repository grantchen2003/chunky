package server

import (
	"fmt"
	"net/http"

	"github.com/grantchen2003/chunky/internal"
	"github.com/grantchen2003/chunky/internal/database"
	"github.com/grantchen2003/chunky/internal/filestorer"
	"github.com/grantchen2003/chunky/internal/handler"
)

type Server struct {
	port                 string
	handlerToEndpoint    map[string]string
	db                   database.Database
	fileStorer           filestorer.FileStorer
	uploadSessionService *internal.UploadSessionService
}

func NewServer(port string) (*Server, error) {
	db, err := database.NewPostgresql()
	if err != nil {
		return nil, err
	}

	fileStorer := filestorer.NewFileStore()

	uploadSessionService := internal.NewUploadSessionService(db, fileStorer)

	fmt.Println(db, fileStorer, uploadSessionService)

	return &Server{
		port: port,
		handlerToEndpoint: map[string]string{
			"initiateUploadSession": "/initiateUploadSession",
			"byteRangesToUpload":    "/byteRangesToUpload",
			"uploadFileChunk":       "/uploadFileChunk",
		},
		db:                   db,
		fileStorer:           fileStorer,
		uploadSessionService: uploadSessionService,
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
	fmt.Println(s)

	initiateUploadSessionHandler := handler.NewInitiateUploadSessionHandler(s.uploadSessionService)
	byteRangesToUploadHandler := handler.NewByteRangesToUploadHandler()
	uploadFileChunkHandler := handler.NewUploadFileChunkHandler(s.uploadSessionService)

	http.HandleFunc(s.handlerToEndpoint["initiateUploadSession"], initiateUploadSessionHandler.Handle)
	http.HandleFunc(s.handlerToEndpoint["byteRangesToUpload"], byteRangesToUploadHandler.Handle)
	http.HandleFunc(s.handlerToEndpoint["uploadFileChunk"], uploadFileChunkHandler.Handle)

	return http.ListenAndServe(s.port, nil)
}
