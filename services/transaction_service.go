package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}


func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}


func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}


// func (s *TransactionService) GetAll() ([]models.Transaction, error) {
// 	return s.repo.GetAll()
// }


// func (s *TransactionService) GetByID(id int) (*models.Transaction, error) {
// 	return s.repo.GetByID(id)
// }


// func (s *TransactionService) Update(transaction *models.Transaction) error {
// 	return s.repo.Update(transaction)
// }


// func (s *TransactionService) Delete(id int) error {
// 	return s.repo.Delete(id)
// }
