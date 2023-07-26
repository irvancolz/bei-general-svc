package helper

import (
	"be-idx-tsg/internal/app/httprest/model"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pdf/fpdf"
	"github.com/johnfercher/maroto/pkg/color"
)

type TableHeader struct {
	Title    string
	Width    float64
	Children []TableHeader
}

func (t TableHeader) GetWidth(props []TableHeader) float64 {
	var result float64
	if len(t.Children) <= 0 {
		return t.Width
	}

	for _, header := range props {
		result += header.GetWidth(header.Children)
	}

	return result
}

type PdfTableOptions struct {
	// default is "Bursa Effek Indonesia"
	HeaderTitle string
	// specify each collumn name
	HeaderRows []TableHeader
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

type fpdfPageProperties struct {
	pageLeftPadding  float64
	pageRightpadding float64
	pageTopPadding   float64
	headerHeight     float64
	currentY         float64
	currentX         float64
	lineHeight       float64
}

func ExportTableToPDF(c *gin.Context, data [][]string, filename string, props *PdfTableOptions) (string, error) {

	filenames := filename
	if filenames == "" {
		filenames = "export-to-pdf.pdf"
	}

	pdf := fpdf.NewCustom(createPageConfig(c, props))
	pageProps := fpdfPageProperties{
		pageLeftPadding:  15,
		pageRightpadding: 15,
		pageTopPadding:   5,
	}

	pdf.SetAutoPageBreak(false, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.SetMargins(pageProps.pageLeftPadding, pageProps.pageTopPadding, pageProps.pageRightpadding)

	drawHeader(pdf, props.getHeaderTitle(), &pageProps)
	drawFooter(pdf)

	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	pageWidth, pageHeight := pdf.GetPageSize()
	currentY := pageProps.headerHeight + 10
	lineHeight := float64(6)
	pageProps.lineHeight = lineHeight
	pageProps.currentY = currentY

	var columnWidth []float64
	for _, header := range props.HeaderRows {
		columnWidth = append(columnWidth, header.GetWidth(header.Children))
	}

	var totalWidth int
	for _, col := range columnWidth {
		totalWidth += int(col)
	}

	// centering the table
	tableMarginX := func() float64 {
		if totalWidth < int(pageWidth)-int(pageProps.pageLeftPadding)-int(pageProps.pageRightpadding) {
			return (pageWidth - float64(totalWidth)) / 2
		}

		return pageProps.pageLeftPadding
	}()

	currentX := tableMarginX
	pageProps.currentX = currentX
	pdf.SetLeftMargin(currentX)
	// drawTableHeader

	drawTableHeader(pdf, props.HeaderRows, &pageProps)
	currentY = pageProps.currentY

	for r, rows := range data {
		currentX = tableMarginX
		maxColHeight := getHighestCol(pdf, columnWidth, rows)

		// reset properties when add page
		if currentY+float64(maxColHeight) > pageHeight-30 {
			pdf.AddPage()
			pdf.SetPage(pdf.PageNo() + 1)
			currentY = pageProps.headerHeight + 10
		}

		pdf.SetFontStyle("")
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFillColor(240, 240, 240)

		if r%2 != 0 {
			pdf.SetAlpha(0, "Normal")
		}

		currRowsheight := float64(maxColHeight) * lineHeight

		// draw bg
		pdf.Rect(currentX, currentY, float64(totalWidth), currRowsheight, "F")
		pdf.SetAlpha(1, "Normal")

		pdf.SetX(currentX)
		pdf.SetY(currentY)

		for colNumber, col := range rows {

			// prevent out of index panic
			currColWidth := func() float64 {
				if colNumber > len(columnWidth)-1 {
					return 20
				}

				return columnWidth[colNumber]
			}()

			pdf.SetY(currentY)
			pdf.SetX(currentX)

			pdf.SetAlpha(.25, "Normal")
			pdf.Line(currentX, currentY, currentX, currentY+currRowsheight)
			if colNumber == len(rows)-1 {
				pdf.Line(pageWidth-pageProps.pageRightpadding, currentY, pageWidth-pageProps.pageRightpadding, currentY+currRowsheight)

			}
			pdf.SetAlpha(1, "Normal")

			splittedtext := pdf.SplitLines([]byte(col), currColWidth)
			for _, text := range splittedtext {
				pdf.CellFormat(currColWidth, lineHeight, string(text), "", 2, "C", false, 0, getLink(col))
			}
			currentX += currColWidth
		}

		currentY += float64(maxColHeight) * lineHeight
		pageProps.currentY = currentY

	}

	err := pdf.OutputFileAndClose(filenames)
	if err != nil {
		log.Println("failed create Pdf files :", err)
		return "", err
	}
	return filenames, nil
}

func ExportAnnouncementToPdf(c *gin.Context, data model.Announcement, opt PdfTableOptions, filename string) (string, error) {

	filenames := filename
	if filenames == "" {
		filenames = "export-to-pdf.pdf"
	}

	pdf := fpdf.NewCustom(createPageConfig(c, &opt))
	pageProps := fpdfPageProperties{
		pageLeftPadding:  15,
		pageRightpadding: 15,
		pageTopPadding:   5,
	}

	pdf.SetAutoPageBreak(false, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.SetMargins(pageProps.pageLeftPadding, pageProps.pageTopPadding, pageProps.pageRightpadding)

	drawHeader(pdf, opt.getHeaderTitle(), &pageProps)
	drawFooter(pdf)

	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	pageWidth, pageHeight := pdf.GetPageSize()
	currentY := pageProps.headerHeight + 10
	lineHeight := float64(8)
	rowsWidth := pageWidth - pageProps.pageLeftPadding - pageProps.pageRightpadding
	pdf.SetLeftMargin(pageProps.pageLeftPadding)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetY(currentY)
	pdf.MultiCell(rowsWidth, lineHeight, fmt.Sprintf("Jenis informasi : %s", data.InformationType), "", "", false)

	currentY += lineHeight
	pdf.SetY(currentY)
	pdf.SetTextColor(29, 31, 29)
	pdf.SetFont("Arial", "I", 12)
	pdf.CellFormat(pageWidth/4, lineHeight, fmt.Sprintf("Dibuat oleh : %s", data.Creator), "", 0, "", false, 0, "")

	pdf.CellFormat(pageWidth/4, lineHeight, fmt.Sprintf("Dibuat pada : %s", data.EffectiveDate.Format("15-06-2006")), "", 0, "", false, 0, "")

	currentY += lineHeight + 10
	pdf.SetFont("Arial", "", 12)
	splittedText := pdf.SplitLines([]byte(data.Regarding), rowsWidth)

	for _, line := range splittedText {

		if currentY+lineHeight >= pageHeight-20 {
			pdf.AddPage()
			pdf.SetPage(pdf.PageNo() + 1)
			currentY = pageProps.headerHeight + 10
			pdf.SetLeftMargin(pageProps.pageLeftPadding)
		}

		pdf.SetY(currentY)
		pdf.SetTextColor(29, 31, 29)
		pdf.CellFormat(pageWidth/4, lineHeight, string(line), "", 2, "", false, 0, "")
		currentY += lineHeight
	}

	err := pdf.OutputFileAndClose(filenames)
	if err != nil {
		log.Println("failed create Pdf files :", err)
		return "", err
	}
	return filenames, nil
}

func drawFooter(pdf *fpdf.Fpdf) {
	pageWidth, pageHeight := pdf.GetPageSize()
	footerHeight := 10

	pdf.SetFooterFunc(func() {
		footerImgHeight := 7
		pdf.SetTextColor(0, 0, 0)
		// bottom Line
		pdf.SetFillColor(240, 240, 240)
		pdf.Rect(0, pageHeight-float64(footerHeight), pageWidth, float64(footerHeight), "F")

		// footer background
		pdf.SetDrawColor(159, 14, 15)
		pdf.SetLineWidth(.8)
		pdf.Line(0, pageHeight-float64(footerHeight), pageWidth, pageHeight-float64(footerHeight))

		// current Time
		pdf.SetFont("Times", "", 10)
		pdf.SetLeftMargin(0)
		pdf.SetY(pageHeight - float64(footerHeight))
		footerDate := func() string {
			return time.Now().Format("02/01/2006") + " - Page " + fmt.Sprintf("%v", pdf.PageNo()) + " Of " + fmt.Sprintf("%v", pdf.PageCount())
		}()
		pdf.MultiCell(50, 8, footerDate, "", "C", false)

		// app name
		pdf.SetY(pageHeight - float64(footerHeight))
		appNameWidth := 50
		pdf.SetX(pageWidth - float64(appNameWidth))
		pdf.MultiCell(float64(appNameWidth), 8, "Sistem IDX Portal", "", "C", false)

		// idx footer logo
		pdf.ImageOptions("internal/app/helper/idx-logo-2.png", (pageWidth-float64(footerImgHeight))/2, pageHeight-float64(footerHeight)+2, float64(footerImgHeight), float64(footerImgHeight), false, fpdf.ImageOptions{}, 0, "")
	})
}

func drawHeader(pdf *fpdf.Fpdf, title string, pageProps *fpdfPageProperties) {
	pageProps.headerHeight = 30
	pdf.SetHeaderFunc(func() {
		pdf.SetTextColor(0, 0, 0)
		pageWidth, pageHeight := pdf.GetPageSize()
		headerTitleWidth := pageWidth / 2
		headerTitle := title
		headerImgHeight := pageProps.headerHeight - 5
		a4PageWidth := 210
		watermarkWidth := float64(a4PageWidth) * 0.85

		// watermark
		pdf.SetAlpha(0.1, "Normal")
		pdf.ImageOptions("internal/app/helper/icon-globe-idx.png", (pageWidth-watermarkWidth)/2, (pageHeight-watermarkWidth)/2, watermarkWidth, watermarkWidth, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
		pdf.SetAlpha(1, "Normal")

		// header bg
		pdf.SetFillColor(240, 240, 240)
		pdf.Rect(0, 0, pageWidth, float64(pageProps.headerHeight+pageProps.pageTopPadding), "F")

		// header logo
		pdf.ImageOptions("internal/app/helper/icon-globe-idx.png", pageProps.pageLeftPadding, pageProps.pageTopPadding, headerImgHeight, headerImgHeight, false, fpdf.ImageOptions{}, 0, "")

		// header title
		pdf.SetLeftMargin((pageWidth - headerTitleWidth) / 2)
		pdf.SetFontSize(18)
		pdf.SetFontStyle("B")
		pdf.MultiCell(headerTitleWidth, 10, headerTitle, "", "C", false)
		pdf.SetLeftMargin(0)

		// header Bot Line
		pdf.SetDrawColor(159, 14, 15)
		pdf.SetLineWidth(.8)
		pdf.Line(0, pageProps.headerHeight+pageProps.pageTopPadding, pageWidth, pageProps.headerHeight+pageProps.pageTopPadding)

	})
}

func getHighestCol(pdf *fpdf.Fpdf, colWidth []float64, data []string) int {
	result := 0

	for i, text := range data {
		currTextHeight := pdf.SplitLines([]byte(text), colWidth[i])
		if len(currTextHeight) > result {
			result = len(currTextHeight)
		}
	}
	return result
}

func getLink(str string) string {
	pattern := `^(https?|ftp|file):\/\/[-\w+&@#/%?=~_|!:,.;]*[-\w+&@#/%=~_|]$`
	reg := regexp.MustCompile(pattern)

	if reg.MatchString(str) {
		return str
	}

	return ""
}

// return A4 size by default
func (opt *PdfTableOptions) getPageSize() fpdf.SizeType {
	results := fpdf.SizeType{}

	results.Ht = func() float64 {
		if opt.Papperheight <= 0 {
			return 297.0
		}
		return opt.Papperheight
	}()

	results.Wd = func() float64 {
		if opt.PapperWidth <= 0 {
			return 210.0
		}
		return opt.PapperWidth
	}()

	return results
}
func createPageConfig(c *gin.Context, props *PdfTableOptions) *fpdf.InitType {
	results := fpdf.InitType{
		OrientationStr: props.getPageOrientation(c),
		UnitStr:        "mm",
		FontDirStr:     "",
		Size:           props.getPageSize(),
	}

	return &results
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

func (opt *PdfTableOptions) getHeaderTitle() string {
	if opt.HeaderTitle == "" {
		return "Bursa Efek Indonesia"
	}
	return opt.HeaderTitle
}

func drawTableHeader(pdf *fpdf.Fpdf, headers []TableHeader, pageProps *fpdfPageProperties) {
	if len(headers) <= 0 {
		return
	}

	currentX := pageProps.currentX
	lineHeight := pageProps.lineHeight
	currentY := pageProps.currentY

	pdf.SetFontStyle("B")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFillColor(50, 117, 168)

	// get how high this row are
	var headerTexts []string
	var headerWidths []float64
	for _, header := range headers {
		headerTexts = append(headerTexts, header.Title)
		headerWidths = append(headerWidths, header.GetWidth(header.Children))
	}

	maxColHeight := getHighestCol(pdf, headerWidths, headerTexts)
	currRowsheight := float64(maxColHeight) * lineHeight

	var totalWidth int
	for _, col := range headerWidths {
		totalWidth += int(col)
	}
	// draw bg
	pdf.Rect(currentX, currentY, float64(totalWidth), currRowsheight, "F")

	for _, header := range headers {
		pdf.SetLeftMargin(currentX)
		pdf.SetY(currentY)
		curColWidth := header.GetWidth(header.Children)

		splittedtext := pdf.SplitLines([]byte(header.Title), curColWidth)
		for _, text := range splittedtext {
			pdf.CellFormat(curColWidth, lineHeight, string(text), "", 2, "C", false, 0, getLink(header.Title))
		}

		currentX += curColWidth
		if len(header.Children) > 0 {
			currentY += lineHeight
		}

		pageProps.currentY = currentY
		drawTableHeader(pdf, header.Children, pageProps)
		pageProps.currentX = currentX
	}
	pageProps.currentY += float64(maxColHeight * int(lineHeight))

}

func GenerateTableHeaders(titles []string, widths []float64) []TableHeader {
	var result []TableHeader

	for i, title := range titles {
		item := TableHeader{
			Title: title,
			Width: func() float64 {
				if i >= len(widths) {
					return 20
				}
				return widths[i]
			}(),
		}
		result = append(result, item)
	}

	return result
}
