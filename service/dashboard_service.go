package service

import (
	"go-admin/models"
	"time"

	"gorm.io/gorm"
)

type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

type ChartData struct {
	Date       string  `json:"date"`
	GrandTotal float64 `json:"grand_total"`
}

type BestProduct struct {
	Title string  `json:"title"`
	Total float64 `json:"total"`
}

func (s *DashboardService) GetDashboardData() (map[string]interface{}, error) {
	now := time.Now()
	last7Days := now.AddDate(0, 0, -7)

	// Chart sales for last 7 days
	var chartSales []ChartData
	err := s.db.
		Table("transactions").
		Select("DATE(created_at) as date, SUM(grand_total) as grand_total").
		Where("created_at >= ?", last7Days).
		Group("DATE(created_at)").
		Scan(&chartSales).Error
	if err != nil {
		return nil, err
	}

	var salesDate []string
	var grandTotal []float64
	if len(chartSales) > 0 {
		for _, data := range chartSales {
			salesDate = append(salesDate, data.Date)
			grandTotal = append(grandTotal, data.GrandTotal)
		}
	} else {
		salesDate = []string{""}
		grandTotal = []float64{0}
	}

	// Count sales today
	var countSalesToday int64
	s.db.Model(&models.Transaction{}).Where("created_at::date = CURRENT_DATE").Count(&countSalesToday)

	// Sum sales today
	var sumSalesToday float64
	s.db.Model(&models.Transaction{}).Select("SUM(grand_total)").Where("created_at::date = CURRENT_DATE").Scan(&sumSalesToday)

	// Sum profits today
	var sumProfitsToday float64
	s.db.Model(&models.Profit{}).Select("SUM(total)").Where("created_at::date = CURRENT_DATE").Scan(&sumProfitsToday)

	// Products with low stock
	var productsLimitStock []models.Product
	s.db.Where("stock <= ?", 10).Find(&productsLimitStock)

	// Best selling products
	var bestProducts []BestProduct
	s.db.Table("transaction_details").
		Select("products.title as title, SUM(transaction_details.qty) as total").
		Joins("join products on products.id = transaction_details.product_id").
		Group("transaction_details.product_id, products.title").
		Order("total DESC").
		Limit(5).
		Scan(&bestProducts)

	var productTitles []string
	var totalQty []float64
	if len(bestProducts) > 0 {
		for _, p := range bestProducts {
			productTitles = append(productTitles, p.Title)
			totalQty = append(totalQty, p.Total)
		}
	} else {
		productTitles = []string{""}
		totalQty = []float64{0}
	}

	return map[string]interface{}{
		"sales_date":           salesDate,
		"grand_total":          grandTotal,
		"count_sales_today":    countSalesToday,
		"sum_sales_today":      sumSalesToday,
		"sum_profits_today":    sumProfitsToday,
		"products_limit_stock": productsLimitStock,
		"product":              productTitles,
		"total":                totalQty,
	}, nil
}
