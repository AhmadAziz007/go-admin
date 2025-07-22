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

type ProfitService struct {
	DB *gorm.DB
}

func NewProfitService(db *gorm.DB) *ProfitService {
	return &ProfitService{DB: db}
}

func (s *ProfitService) FilterProfits(startDate, endDate string) ([]models.Profit, float64, error) {
	// Parse tanggal dari string ke time.Time
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, 0, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, 0, err
	}

	var profits []models.Profit
	if err := s.DB.Preload("Transaction").
		Where("DATE(created_at) BETWEEN ? AND ?", start, end).
		Find(&profits).Error; err != nil {
		return nil, 0, err
	}

	var total_profit float64
	if err := s.DB.Model(&models.Profit{}).
		Where("DATE(created_at) BETWEEN ? AND ?", start, end).
		Select("SUM(total)").Scan(&total_profit).Error; err != nil {
		return nil, 0, err
	}

	return profits, total_profit, nil
}

func (s *ProfitService) ExportExcel(startDate, endDate string) (*excelize.File, error) {
	profits, _, err := s.FilterProfits(startDate, endDate)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Sales Report"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1") //Hapus sheet default

	header := []string{"No", "Date", "Invoice", "Total"}
	for i, h := range header {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, s.headerStyle(f))
	}

	var total float64
	for i, profit := range profits {
		row := i + 2
		transaction := "Kode invoice"
		if profit.Transaction != nil {
			transaction = profit.Transaction.Invoice
		}

		data := []interface{}{
			i + 1,
			profit.CreatedAt.Format("2006-01-02 15:04:05"),
			transaction,
			profit.Total,
		}

		for j, d := range data {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, d)
		}
		total += profit.Total
	}
	totalRow := len(profits) + 2
	f.SetCellValue(sheet, "A"+strconv.Itoa(totalRow), "TOTAL")
	f.MergeCell(sheet, "A"+strconv.Itoa(totalRow), "C"+strconv.Itoa(totalRow))
	f.SetCellValue(sheet, "D"+strconv.Itoa(totalRow), total)
	f.SetCellStyle(sheet, "A"+strconv.Itoa(totalRow), "D"+strconv.Itoa(totalRow), s.totalStyle(f))

	for i := range header {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, col, col, 20)
	}
	return f, nil
}

func (s *ProfitService) headerStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E6E6E7"}, Pattern: 1},
	})
	return style
}

func (s *ProfitService) totalStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Font:   &excelize.Font{Bold: true},
		Fill:   excelize.Fill{Type: "pattern", Color: []string{"#E6E6E7"}, Pattern: 1},
		NumFmt: 4, // Format uang
	})
	return style
}

func (s *ProfitService) ExportPDF(startDate, endDate string) ([]byte, error) {
	profits, totalProfit, err := s.FilterProfits(startDate, endDate)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "PROFIT REPORT VUE AND GO MANAGEMENT", "", 1, "C", false, 0, "")
	pdf.Ln(1)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "Alamat: Kampung Jati, Kelurahan Jatinegara Kaum, Kecamatan Pulo Gadung, Jakarta Timur", "", 1, "C", false, 0, "")
	pdf.Ln(1)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "Telp: 0812-8888-9999", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 10, "Period: "+startDate+" to "+endDate, "", 1, "L", false, 0, "")
	pdf.Ln(5)

	headers := []string{"No", "Date", "Invoice", "Total"}

	colWidths := []float64{10, 40, 50, 50, 50, 40}
	for i, header := range headers {
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	for i, profit := range profits {
		transaction := "Kode invoice"
		if profit.Transaction != nil {
			transaction = profit.Transaction.Invoice
		}

		data := []string{
			strconv.Itoa(i + 1),
			profit.CreatedAt.Format("2006-01-02 15:04:05"),
			transaction,
			"Rp. " + strconv.FormatFloat(profit.Total, 'f', 0, 64),
		}

		for j, d := range data {
			pdf.CellFormat(colWidths[j], 10, d, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Total
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(colWidths[0]+colWidths[1]+colWidths[2], 10, "TOTAL", "1", 0, "C", false, 0, "")
	pdf.CellFormat(colWidths[3], 10, "Rp. "+strconv.FormatFloat(totalProfit, 'f', 0, 64), "1", 0, "C", false, 0, "")

	// Simpan ke buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
