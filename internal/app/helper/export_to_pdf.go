package helper

import (
	"be-idx-tsg/internal/app/httprest/model"
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
	// setting paper width
	PapperWidth float64
	// setting papper height
	Papperheight float64
	// setting page orientation by default "p" / "l"
	PageOrientation string
}

func (opt *PdfTableOptions) getPageOrientation(c *gin.Context) string {
	pageOrientation := c.Query("orientation")
	if c.Query("orientation") == "" && (opt.PageOrientation == "" || !IsContains([]string{"p", "l"}, opt.PageOrientation)) {
		pageOrientation = "p"
		return pageOrientation
	}

	if c.Query("orientation") == "" && IsContains([]string{"p", "l"}, opt.PageOrientation) {
		pageOrientation = opt.PageOrientation
		return pageOrientation
	}
	return pageOrientation
}

func ExportTableToPDF(c *gin.Context, data [][]string, filename string, opt PdfTableOptions) (string, error) {

	pageOrientation := opt.getPageOrientation(c)

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

	pdfFile := pdf.NewMarotoCustomSize(consts.Orientation(strings.ToUpper(pageOrientation)), consts.PageSize(pageSize), "mm", opt.PapperWidth, opt.Papperheight)

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
	if config.HeaderTitle == "" {
		title = "Bursa Effek Indonesia"
	}

	logoPath := config.FooterLogo
	if config.FooterLogo == "" {
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

			page.Col(columnTotal-uint(2*logoWidth), func() {
				page.Text(title, props.Text{
					Style: consts.Bold,
					Size:  18,
					Align: consts.Middle,
				})
			})
			page.ColSpace(uint(logoWidth))
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

func ExportAnnouncementToPdf(c *gin.Context, data model.Announcement, opt PdfTableOptions, filename string) (string, error) {
	filenames := filename
	if filename == "" {
		filenames = "export-to-pdf.pdf"
	}

	pageOrientation := opt.getPageOrientation(c)

	pageSize := opt.PageSize
	if pageSize == "" {
		pageSize = string(consts.A4)
	}

	pdfFile := pdf.NewMarotoCustomSize(consts.Orientation(strings.ToUpper(pageOrientation)), consts.PageSize(pageSize), "mm", opt.PapperWidth, opt.Papperheight)
	pdfFile.SetMaxGridSum(12)
	createHeader(pdfFile, opt, 12)
	createFooter(pdfFile, opt, 12)

	titleHeight := 8
	creatorHeight := 16
	pdfFile.Row(float64(titleHeight), func() {
		infoType := func() string {
			var infoType string
			if infoType == "" {
				infoType = "SEMUA"
				return infoType
			}
			if strings.EqualFold(data.InformationType, "AB") {
				infoType = "ANGGOTA BURSA"
				return infoType
			}
			return infoType
		}()
		pdfFile.Text("Jenis Informasi :"+infoType, props.Text{
			Size:  16,
			Style: consts.Bold,
		})
	})

	pdfFile.Row(float64(creatorHeight), func() {
		pdfFile.Col(4, func() {
			pdfFile.Text("dibuat oleh : "+data.Creator, props.Text{
				Style: consts.Italic,
				Size:  12,
				Color: color.Color{
					Red:   141,
					Green: 137,
					Blue:  137},
			})
		})
		pdfFile.Col(4, func() {
			pdfFile.Text("dibuat pada :"+data.EffectiveDate.Format("15-06-2006"), props.Text{
				Style: consts.Italic,
				Size:  12,
				Color: color.Color{
					Red:   141,
					Green: 137,
					Blue:  137},
			})
		})
	})

	pdfFile.Row(10, func() {
		pdfFile.Text(data.Regarding, props.Text{
			Size: 12,
		})
	})

	errSave := pdfFile.OutputFileAndClose(filenames)
	if errSave != nil {
		log.Println("failed to save pdf :", errSave)
		return "", errSave
	}
	return filenames, nil
}
