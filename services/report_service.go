package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"log"
	"os"
)

type ReportService struct {
	repo *repositories.ReportRepository
}


func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) ReportToday() (*models.ReportToday, error) {
	return s.repo.ReportToday()
}


func (s *ReportService) ReportRange(startDate, endDate string) (*models.ReportRange, error) {
	file, err := os.OpenFile("main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)

	logger.Println("Service: ReportRange")
	logger.Println("Start date: " + startDate)
	logger.Println("End date: " + endDate)	

	return s.repo.ReportRange(startDate, endDate)
}