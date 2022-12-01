package main

import (
	"log"
	"net/http"
	"os"
	"time"

	ics "github.com/arran4/golang-ical"
)

func handler(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("CALENDAR_URL")

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	calendar, err := ics.ParseCalendar(resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newCalendar := ics.NewCalendar()
	newCalendar.SetMethod(ics.MethodRequest)

	for _, event := range calendar.Events() {
		start, err := event.GetStartAt()
		if err != nil {
			log.Println("event skipped:", err)
			continue
		}
		end, err := event.GetEndAt()
		if err != nil {
			log.Println("event skipped:", err)
			continue
		}

		duration := end.Sub(start)
		if duration > 24*time.Hour {
			continue
		}

		newCalendar.AddVEvent(event)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/calendar")
	if err := newCalendar.SerializeTo(w); err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
