package handler

import (
	"analytics/data"
	"analytics/storage"
	"log"
	"net/http"
	"time"
)

type EventsManager struct {
	l *log.Logger
	s storage.Storage
}

// NewEventsManager function creates entity which handles incoming requests
func NewEventsManager(l *log.Logger, s storage.Storage) *EventsManager {
	return &EventsManager{l, s}
}

// ProcessEvent method handles incoming events
func (em EventsManager) ProcessEvent(w http.ResponseWriter, r *http.Request) {
	em.l.Println("Handle incoming event")

	event := &data.Event{
		Timestamp: time.Now().UTC(),
	}

	err := event.FromJSON(r.Body)
	if err != nil {
		em.l.Printf("Unable decoding event: %v", err)
		http.Error(w, "Unable decoding event:", http.StatusBadRequest)
		return
	}

	if err := em.s.SendAsync(event.Timestamp, event.Value); err != nil {
		em.l.Printf("Unable writing event: %v", err)
		http.Error(w, "Unable writing event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	em.l.Printf("Processed event with value: %v, timestamp: %v\n", event.Value, event.Timestamp)
}

// GetAvgForLastMin method handles requests to get the average value of the events received in th last min
func (em EventsManager) GetAvgForLastMin(w http.ResponseWriter, r *http.Request) {
	em.l.Println("Handle incoming GET avg request")

	result, err := em.s.GetValuesForLastMin(time.Now())
	if err != nil {
		em.l.Printf("Unable to get values from DB: %v", err)
		http.Error(w, "Unable to get values from DB", http.StatusInternalServerError)
		return
	}

	resp := &data.AvgResponse{
		Rows: 0,
		Avg:  0,
	}

	if len(result) == 0 {
		if err := resp.ToJSON(w); err != nil {
			http.Error(w, "Unable encoding to JSON", http.StatusInternalServerError)
		}

		return
	}

	var counter int
	for _, v := range result {
		counter += v
	}

	resp.Rows = len(result)
	resp.Avg = counter / len(result)

	if err := resp.ToJSON(w); err != nil {
		http.Error(w, "Unable encoding to JSON", http.StatusInternalServerError)
	}
}
