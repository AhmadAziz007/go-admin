package service

import (
	"bytes"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"go-admin/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type SalesService struct {
	db *gorm.DB
}

func NewSalesService(db *gorm.DB) *SalesService {
	return &SalesService{db: db}
}

func (s *SalesService) FilterSales(startDate, endDate string) ([]models.Transaction, float64, error) {
	var sales []models.Transaction
	var total float64

	// Parse the input dates to time.Time
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, 0, err
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, 0, err
	}

	// Query transactions with relations
	if err := s.db.Preload("User").Preload("Customer").Preload("TransactionDetails").
		Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", startDate, endDate).
		Find(&sales).Error; err != nil {
		return nil, 0, err
	}

	// Calculate total sales
	if err := s.db.Model(&models.Transaction{}).
		Select("SUM(grand_total)").
		Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", startDate, endDate).
		Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	return sales, total, nil
}

func (s *SalesService) ExportExcel(startDate, endDate string) (*excelize.File, error) {
	sales, _, err := s.FilterSales(startDate, endDate)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Sales Report"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1") //Hapus sheet default

	//Set header
	header := []string{"No", "Date", "Invoice", "Cashier", "Customer", "Discount", "Total"}
	for i, h := range header {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, s.headerStyle(f))
	}

	var grandTotal float64
	for i, sale := range sales {
		row := i + 2
		cashier := "Kasir tidak diketahui"
		if sale.User != nil {
			cashier = sale.User.FirstName + " " + sale.User.LastName
		}

		customer := "Umum"
		if sale.Customer != nil {
			customer = sale.Customer.Name
		}

		data := []interface{}{
			i + 1,
			sale.CreatedAt.Format("2006-01-02 15:04:05"),
			sale.Invoice,
			cashier,
			customer,
			sale.Discount,
			sale.GrandTotal,
		}

		for j, d := range data {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, d)
		}
		grandTotal += sale.GrandTotal
	}

	totalRow := len(sales) + 2
	f.SetCellValue(sheet, "A"+strconv.Itoa(totalRow), "TOTAL SALES")
	f.MergeCell(sheet, "A"+strconv.Itoa(totalRow), "F"+strconv.Itoa(totalRow))
	f.SetCellValue(sheet, "G"+strconv.Itoa(totalRow), grandTotal)
	f.SetCellStyle(sheet, "A"+strconv.Itoa(totalRow), "G"+strconv.Itoa(totalRow), s.totalStyle(f))

	for i := range header {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, col, col, 20)
	}
	return f, nil
}

func (s *SalesService) headerStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E6E6E7"}, Pattern: 1},
	})
	return style
}

func (s *SalesService) totalStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Font:   &excelize.Font{Bold: true},
		Fill:   excelize.Fill{Type: "pattern", Color: []string{"#E6E6E7"}, Pattern: 1},
		NumFmt: 4, // Format uang
	})
	return style
}

func (s *SalesService) ExportPDF(startDate, endDate string) ([]byte, error) {
	sales, total, err := s.FilterSales(startDate, endDate)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "SALES REPORT VUE AND GO MANAGEMENT", "", 1, "C", false, 0, "")
	pdf.Ln(1)

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Alamat: Kampung Jati, Kelurahan Jatinegara Kaum, Kecamatan Pulo Gadung, Jakarta Timur", "", 1, "C", false, 0, "")
	pdf.Ln(1)

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Telp: 0812-8888-9999", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 10, "Period: "+startDate+" to "+endDate, "", 1, "L", false, 0, "")
	pdf.Ln(5)

	headers := []string{"No", "Date", "Invoice", "Cashier", "Customer", "Discount", "Total"}
	colWidths := []float64{10, 40, 50, 50, 50, 30, 40}

	pdf.SetFont("Arial", "B", 12)
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	for i, sale := range sales {
		cashier := "Kasir Tidak Diketahui"
		if sale.User != nil {
			cashier = sale.User.FirstName + " " + sale.User.LastName
		}

		customer := "Umum"
		if sale.Customer != nil {
			customer = sale.Customer.Name
		}

		date := sale.CreatedAt.Format("2006-01-02 15:04:05")

		discountStr := strconv.FormatFloat(sale.Discount, 'f', 0, 64)

		data := []string{
			strconv.Itoa(i + 1),
			date,
			sale.Invoice,
			cashier,
			customer,
			discountStr,
			"Rp. " + strconv.FormatFloat(sale.GrandTotal, 'f', 0, 64),
		}

		for j, d := range data {
			pdf.CellFormat(colWidths[j], 10, d, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	pdf.SetFont("Arial", "B", 12)

	totalColWidth := colWidths[0] + colWidths[1] + colWidths[2] + colWidths[3] + colWidths[4]
	pdf.CellFormat(totalColWidth, 10, "TOTAL SALES", "1", 0, "C", false, 0, "")

	pdf.CellFormat(colWidths[5], 10, "", "1", 0, "C", false, 0, "")

	pdf.CellFormat(colWidths[6], 10, "Rp. "+strconv.FormatFloat(total, 'f', 0, 64), "1", 1, "C", false, 0, "")

	// Simpan ke buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
