package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}


func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}


func (repo *ReportRepository) ReportToday() (*models.ReportToday, error) {
	sqlTotal := "SELECT COALESCE(SUM(total_amount),0) AS total_revenue, COUNT(*) AS total_transaksi FROM transactions WHERE DATE(created_at)=CURRENT_DATE"
	
	var totalRevenue, totalTransaction int

	transactionToday := repo.db.QueryRow(sqlTotal).Scan(&totalRevenue, &totalTransaction)
	
	if transactionToday != nil {
		return &models.ReportToday {
			TotalRevenue: 0,
			TotalTransaksi: 0,
		}, nil
	}

	var productName string
	var qtySold int

	sqlTopSold := "SELECT p.name AS product_name, COALESCE(SUM(d.quantity),0) AS qty_sold FROM transaction_details d INNER JOIN products p ON d.product_id=p.id INNER JOIN transactions t ON t.id=d.transaction_id  WHERE DATE(created_at)=CURRENT_DATE GROUP BY p.name LIMIT 1"

	topSold := repo.db.QueryRow(sqlTopSold).Scan(&productName, &qtySold)

	if topSold != nil {		
		return &models.ReportToday {
			TotalRevenue: totalRevenue,
			TotalTransaksi: totalTransaction,
			ProdukTerlaris: []models.BestSeller {
				{
					ProductName: productName,
					QtySold: qtySold,
				},
			},
		}, nil
	}

	return &models.ReportToday {
		TotalRevenue: totalRevenue,
		TotalTransaksi: totalTransaction,
		ProdukTerlaris: []models.BestSeller {
			{
				ProductName: productName,
				QtySold: qtySold,
			},
		},
	}, nil
}



func (repo *ReportRepository) ReportRange(startDate, endDate string) (*models.ReportRange, error) {
	var totalRevenue, totalTransaction, qtySold int
	var productName string

	rsTotal := repo.db.QueryRow("SELECT COALESCE(SUM(total_amount),0) AS total_revenue, COUNT(*) AS total_transaksi FROM transactions WHERE DATE(created_at) BETWEEN $1 AND $2", startDate, endDate).Scan(&totalRevenue, &totalTransaction)

	if rsTotal != nil {		
		return nil, rsTotal
	}

	topSold := repo.db.QueryRow("SELECT p.name AS product_name, COALESCE(SUM(d.quantity),0) AS qty_sold FROM transaction_details d INNER JOIN products p ON d.product_id=p.id INNER JOIN transactions t ON t.id=d.transaction_id WHERE DATE(t.created_at) BETWEEN $1 AND $2 GROUP BY p.name LIMIT 1", startDate, endDate).Scan(&productName, &qtySold)

	if topSold != nil {		
		return nil, topSold
	}

	return &models.ReportRange {
		TotalRevenue: totalRevenue,
		TotalTransaksi: totalTransaction,
		ProdukTerlaris: []models.BestSeller {
			{
				ProductName: productName,
				QtySold: qtySold,
			},
		},
	}, nil
}