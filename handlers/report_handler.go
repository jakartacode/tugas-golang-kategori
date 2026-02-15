package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
)


type ReportHandler struct {
	service *services.ReportService
}


func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleReportToday(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			h.ReportToday(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func (h *ReportHandler) ReportToday(w http.ResponseWriter, r *http.Request) {

	reportToday, err := h.service.ReportToday()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reportToday)
}


func (h *ReportHandler) HandleReportRange(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			h.ReportRange(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func (h *ReportHandler) ReportRange(w http.ResponseWriter, r *http.Request) {

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	file, err := os.OpenFile("main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)

	logger.Println("Handler: ReportRange")
	logger.Println("Start date: " + startDate)
	logger.Println("End date: " + endDate)	

	reportRange, err := h.service.ReportRange(startDate, endDate)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reportRange)
}