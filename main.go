package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type entry struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

type store struct {
	mutex sync.RWMutex
	data  map[string]string // "/gh" -> "https://github.com"
}

var theWholeStore = &store{data: map[string]string{}}

func main() {

	// Redirect handler
	http.HandleFunc("/", handleRedirect)

	// Create/update: POST /urls  {"short":"/gh","long":"https://github.com"}
	http.HandleFunc("/urls", handlePostURL)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRedirect(writer http.ResponseWriter, request *http.Request) {
	theWholeStore.mutex.RLock()
	target, ok := theWholeStore.data[request.URL.Path]
	theWholeStore.mutex.RUnlock()
	if !ok {
		http.NotFound(writer, request)
		return
	}
	http.Redirect(writer, request, target, http.StatusFound) // 302
}

func handlePostURL(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "use POST", http.StatusMethodNotAllowed)
		return
	}
	var bodyStruct entry

	err := json.NewDecoder(request.Body).Decode(&bodyStruct)
	if err != nil {
		http.Error(writer, "invalid json", http.StatusBadRequest)
		return
	}

	if bodyStruct.Short == "" || bodyStruct.Long == "" {
		http.Error(writer, "both short and long are required", http.StatusBadRequest)
		return
	}

	// Ensure short starts with '/'
	if bodyStruct.Short[0] != '/' {
		bodyStruct.Short = "/" + bodyStruct.Short
	}

	err = addToWholeStore(bodyStruct.Short, bodyStruct.Long)
	if err != nil {
		http.Error(writer, "could not store the url", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

func addToWholeStore(short, long string) error {
	theWholeStore.mutex.Lock()
	theWholeStore.data[short] = long
	theWholeStore.mutex.Unlock()
	return nil
}
