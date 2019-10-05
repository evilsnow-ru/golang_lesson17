package service

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang_lesson17/api"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

var storage = make(map[uint64]*api.Event)
var lock sync.RWMutex
var index uint64

func addEvent(event *api.Event) uint64 {
	id := atomic.AddUint64(&index, 1)
	lock.Lock()
	storage[id] = event
	lock.Unlock()
	return id
}

func updateEvent(id uint64, event *api.Event) bool {
	lock.Lock()
	_, ok := storage[id]
	if ok {
		storage[id] = event
	}
	lock.Unlock()
	return ok
}

func deleteEvent(id uint64) {
	lock.Lock()
	delete(storage, id)
	lock.Unlock()
}

func getEvent(id uint64) *api.Event {
	lock.RLock()
	event, ok := storage[id]
	lock.RUnlock()

	if ok {
		return event
	}
	return nil
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

	id := addEvent(event)
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

	event := getEvent(id)

	if event == nil {
		http.NotFound(w, r)
		return
	}

	typeName := api.EventType_name[int32(event.Type)]
	result := fmt.Sprintf("Event {id=%d, type=\"%s\", date=\"%s\", msg=\"%s\"}", event.MsgId, typeName, event.Date, event.Description)
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
