package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/grantchen2003/chunky/internal"
)

func HandleInitiateUploadSession(w http.ResponseWriter, r *http.Request) {
	// log.Printf("Serving %s\n", strings.Split(r.RemoteAddr, ":")[0])

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	type Payload struct {
		FileHash           []byte `json:"fileHash"`
		TotalFileSizeBytes int    `json:"TotalFileSizeBytes"`
	}

	var payload Payload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	sessionId, err := internal.GenerateSessionId(16)
	if err != nil {
		http.Error(w, "Failed to generate session id", http.StatusInternalServerError)
		log.Printf("Error generating session id response: %v", err)
		return
	}

	// need to store sessionId, fileHash, and totalFileSizeBytes somewhere

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{"sessionId": sessionId}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func HandleByteRangesToUpload(w http.ResponseWriter, r *http.Request) {
	// log.Printf("Serving %s\n", strings.Split(r.RemoteAddr, ":")[0])

	responseData := map[string][][2]int{
		"ByteRanges": {{6, 100}, {102, 132}, {103, 104}, {133, 152}, {154, 154}, {155, 155}, {156, 159}, {161, 306}},
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func HandleUploadFileChunk(w http.ResponseWriter, r *http.Request) {
	// log.Printf("Serving %s\n", strings.Split(r.RemoteAddr, ":")[0])

	w.WriteHeader(http.StatusOK)
}

func StartServer(port string, endpointToHandler map[string]func(w http.ResponseWriter, r *http.Request)) {
	for endpoint, handler := range endpointToHandler {
		http.HandleFunc(endpoint, handler)
	}

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(port, nil))
}
