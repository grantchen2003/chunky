package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func HandleInitiateUploadSession(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s\n", strings.Split(r.RemoteAddr, ":")[0])
}

func HandleByteRangesToUpload(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s\n", strings.Split(r.RemoteAddr, ":")[0])

}

func HandleUploadFileChunk(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s\n", strings.Split(r.RemoteAddr, ":")[0])

}

func StartServer(port string, endpointToHandler map[string]func(w http.ResponseWriter, r *http.Request)) {
	for endpoint, handler := range endpointToHandler {
		http.HandleFunc(endpoint, handler)
	}

	fmt.Printf("Server listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
