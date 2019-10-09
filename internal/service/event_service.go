package service

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang_lesson17/api"
	"golang_lesson17/internal/domain/dao"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var storage dao.Storage

func init() {
	var err error
	storage, err = dao.DefaultStorage()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func addEventHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		http.Error(responseWriter, "Empty body", http.StatusBadRequest)
		return
	}

	defer request.Body.Close()
	data, err := ioutil.ReadAll(request.Body)

	if err != nil {
		http.Error(responseWriter, "Error processing request", http.StatusInternalServerError)
		return
	}

	event := &api.Event{}
	err = proto.Unmarshal(data, event)

	if err != nil {
		http.Error(responseWriter, "Error processing request", http.StatusInternalServerError)
		return
	}

	eventEntity, err := api.Convert(event)

	if err != nil {
		http.Error(responseWriter, "Error processing request", http.StatusInternalServerError)
		return
	}

	id, err := storage.Store(eventEntity)

	if err != nil {
		http.Error(responseWriter, "Error processing request", http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(fmt.Sprintf("id: %d", id)))
}

func getEventHandler(w http.ResponseWriter, r *http.Request) {
	idStr, ok := r.URL.Query()["id"]

	if !ok {
		http.Error(w, "id not found", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr[0], 10, 64)

	if err != nil {
		http.Error(w, "Error parse id key", http.StatusBadRequest)
		return
	}

	event, err := storage.Get(id)

	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	if event == nil {
		http.NotFound(w, r)
		return
	}

	typeName := api.EventType_name[int32(event.Type)]
	result := fmt.Sprintf("Event {id=%d, type=\"%s\", date=\"%s\", msg=\"%s\"}", event.Id, typeName, event.Date, event.Description)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func StartServer(port uint32) error {
	if port <= 0 {
		return fmt.Errorf("wrong port: expected value > 0, actual: %d", port)
	}

	fmt.Printf("Starting server at port: %d \n", port)
	http.HandleFunc("/add", addEventHandler)
	http.HandleFunc("/get", getEventHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
