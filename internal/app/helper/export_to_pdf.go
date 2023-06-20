package helper

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type PdfTableOptions struct {
	// default is "Bursa Effek Indonesia"
	HeaderTitle string
	// specify each collumn name
	HeaderRows []string
	// // L (landscape) or P (potrait), default is P
	// Orientation string
	// "A3", "A4", "Legal", "Letter", "A5" default is "A4"
	PageSize string
	// path to logo default is globe idx
	HeaderLogo string
	// path to logo
	FooterLogo string
	// line header and footer color default is maroon
	LineColor *color.Color
	// even indexed rows bg color in table, default is gray
	TableBgCol *color.Color
}

func ExportTableToPDF(c *gin.Context, data [][]string, filename string, opt PdfTableOptions) (string, error) {

	pageOrientation := c.DefaultQuery("orientation", "p")

	headers := opt.HeaderRows
	if len(opt.HeaderRows) <= 0 {
		headers = data[0]
	}

	exportedData := data[1:]

	filenames := filename
	if filenames == "" {
		filenames = "export-to-pdf.pdf"
	}

	pageSize := opt.PageSize
	if pageSize == "" {
		pageSize = string(consts.A4)
	}

	pdfFile := pdf.NewMaroto(consts.Orientation(strings.ToUpper(pageOrientation)), consts.PageSize(pageSize))

	collumnGridSize := len(headers)
	pdfFile.SetMaxGridSum(float64(collumnGridSize))

	collumnwidth := uint(1)
	evenRowsBg := opt.TableBgCol
	if evenRowsBg == nil {
		evenRowsBg = &color.Color{
			Red:   245,
			Green: 245,
			Blue:  245,
		}
	}

	createHeader(pdfFile, opt, uint(collumnGridSize))
	createFooter(pdfFile, opt, uint(collumnGridSize))

	var tableColWidth []uint
	// share same width to each collumn
	for i := 0; i < len(headers); i++ {
		tableColWidth = append(tableColWidth, collumnwidth)
	}

	pdfFile.TableList(headers, exportedData, props.TableList{
		Align:                  consts.Middle,
		VerticalContentPadding: 2,
		HeaderProp: props.TableListContent{
			GridSizes: tableColWidth,
			Style:     consts.Bold,
		},
		ContentProp: props.TableListContent{
			GridSizes: tableColWidth,
		},
		HeaderContentSpace:   2,
		AlternatedBackground: evenRowsBg,
	})

	errSave := pdfFile.OutputFileAndClose(filenames)
	if errSave != nil {
		log.Println("failed to save pdf :", errSave)
		return "", errSave
	}
	return filenames, nil
}

func createHeader(page pdf.Maroto, config PdfTableOptions, columnTotal uint) {
	title := config.HeaderTitle
	if title == "" {
		title = "Bursa Effek Indonesia"
	}

	logoPath := config.FooterLogo
	if logoPath == "" {
		logoPath = "internal/app/helper/icon-globe-idx.png"
	}

	page.RegisterHeader(func() {
		logoWidth := 1
		if columnTotal >= 4 {
			logoWidth = int(columnTotal) / 4
		}

		page.Row(24, func() {
			page.Col(uint(logoWidth), func() {
				_ = page.FileImage(logoPath, props.Rect{
					// left space to make gap with line below
					Percent: 80,
				})
			})

			page.Col(columnTotal-uint(logoWidth), func() {
				page.Text(title, props.Text{
					Style: consts.Bold,
					Size:  18,
					Align: consts.Middle,
				})
			})
		})
		drawLine(page, config)
	})
}

func createFooter(page pdf.Maroto, config PdfTableOptions, columnTotal uint) {
	footerHeight := float64(8)
	footerPaddingTop := footerHeight / 4

	logoPath := config.FooterLogo
	if logoPath == "" {
		logoPath = "internal/app/helper/idx-logo-2.png"
	}

	dateAndNameWidth := 1
	if columnTotal > 8 {
		dateAndNameWidth = 2
	}

	page.RegisterFooter(func() {
		page.Row(footerHeight, func() {
			drawLine(page, config)
		})
		page.Row(footerHeight, func() {
			page.Col(uint(dateAndNameWidth), func() {
				page.Text(time.Now().Format("02 January 2006"), props.Text{
					Align: consts.Left,
					Size:  8,
					Top:   footerPaddingTop,
				})
			})

			page.Col(1, func() {
				errImg := page.FileImage(logoPath, props.Rect{
					Top: footerPaddingTop,
				})
				if errImg != nil {
					log.Println(errImg)
				}
			})
			page.Col(uint(dateAndNameWidth), func() {
				page.Text("Sistem Portal Anggota Bursa", props.Text{
					Align: consts.Left,
					Size:  8,
					Top:   footerPaddingTop,
				})
			})
		})
	})
}

func drawLine(page pdf.Maroto, config PdfTableOptions) {
	pageWidth, _ := page.GetPageSize()
	lineColor := config.LineColor
	if lineColor == nil {
		lineColor = &color.Color{
			Red: 159, Green: 14, Blue: 15,
		}
	}

	page.Line(.5, props.Line{
		Width: pageWidth,
		Style: consts.Solid,
		// red color
		Color: *lineColor,
	})
}
