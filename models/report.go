package models

type BestSeller struct {
	ProductName string `json:"nama"`
	QtySold     int    `json:"qty_terjual"`
}

type ReportToday struct {
	TotalRevenue   int          `json:"total_revenue"`
	TotalTransaksi int          `json:"total_transaksi"`
	ProdukTerlaris []BestSeller `json:"produk_terlaris"`
}

type ReportRange struct {
	TotalRevenue   int          `json:"total_revenue"`
	TotalTransaksi int          `json:"total_transaksi"`
	ProdukTerlaris []BestSeller `json:"produk_terlaris"`
}
