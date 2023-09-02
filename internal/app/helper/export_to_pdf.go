package helper

import (
	"be-idx-tsg/internal/app/httprest/model"
	"fmt"
	"log"
	"math"
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
	// specify each collumn name and size
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
	pageLeftPadding     float64
	pageRightpadding    float64
	pageTopPadding      float64
	headerHeight        float64
	footerHeight        float64
	currentY            float64
	currentX            float64
	lineHeight          float64
	curColWidth         float64
	colWidthList        []float64
	totalColumn         int
	curColIndex         int
	pageHeight          float64
	currRowsheight      float64
	currPageRowHeight   float64
	curRowsBgHeight     float64
	curRowsBgHeightLeft float64
	currRowsIndex       int
	tableWidth          int
	tableMarginX        float64
	currpage            int
	currRowMaxPage      int
	newPageMargin       float64
	footerMarginTop     float64
	footerSpace         float64
	headerSpace         float64
	headerMarginBottom  float64
}

// debugging helper, do not remove
func (p fpdfPageProperties) logStat() {
	log.Printf("value of currRowsIndex: %v", p.currRowsIndex)
	log.Printf("value of currpage: %v", p.currpage)
	log.Printf("value of currRowMaxPage: %v", p.currRowMaxPage)
	// log.Printf("value of pageLeftPadding: %v", p.pageLeftPadding)
	// log.Printf("value of pageRightpadding: %v", p.pageRightpadding)
	// log.Printf("value of pageTopPadding: %v", p.pageTopPadding)
	// log.Printf("value of headerHeight: %v", p.headerHeight)
	log.Printf("value of currentY: %v", p.currentY)
	// log.Printf("value of currentX: %v", p.currentX)
	// log.Printf("value of lineHeight: %v", p.lineHeight)
	// log.Printf("value of curColWidth: %v", p.curColWidth)
	// log.Printf("value of colWidthList: %v", p.colWidthList)
	log.Printf("value of totalColumn: %v", p.totalColumn)
	// log.Printf("value of curColIndex: %v", p.curColIndex)
	// log.Printf("value of pageHeight: %v", p.pageHeight)
	log.Printf("value of currRowsheight: %v", p.currRowsheight)
	log.Printf("value of curRowsBgHeightLeft: %v", p.curRowsBgHeightLeft)
	log.Printf("value of currPageRowHeight: %v", p.currPageRowHeight)
	log.Printf("value of curRowsBgHeight: %v", p.curRowsBgHeight)
	// log.Printf("value of curRowsBgHeightLeft: %v", p.curRowsBgHeightLeft)
	// log.Printf("value of tableWidth: %v", p.tableWidth)
	// log.Printf("value of tableMarginX: %v", p.tableMarginX)
	log.Printf("value of newPageMargin: %v", p.newPageMargin)
	log.Println("========== border here ===================")
}

func (p fpdfPageProperties) isNeedPageBreak(coord float64) bool {
	return coord > p.pageHeight-p.footerSpace
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
	drawFooter(pdf, &pageProps)

	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	pageWidth, pageHeight := pdf.GetPageSize()
	currentY := pageProps.headerSpace
	lineHeight := float64(6)
	pageProps.lineHeight = lineHeight
	pageProps.currentY = currentY
	pageProps.pageHeight = pageHeight

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

	pageProps.tableWidth = totalWidth
	pageProps.tableMarginX = tableMarginX
	pageProps.colWidthList = columnWidth
	pageProps.currRowMaxPage = pdf.PageNo()
	currentX := tableMarginX
	pageProps.currentX = currentX
	pdf.SetLeftMargin(currentX)

	drawTableHeader(pdf, props.HeaderRows, &pageProps)
	drawTable(pdf, &pageProps, data)

	err := pdf.OutputFileAndClose(filenames)
	if err != nil {
		log.Println("failed create Pdf files :", err)
		return "", err
	}
	return filenames, nil
}

func drawTable(pdf *fpdf.Fpdf, pageProps *fpdfPageProperties, data [][]string) {
	for r, rows := range data {
		columnWidth := pageProps.colWidthList
		lineHeight := pageProps.lineHeight
		currentY := pageProps.currentY
		pageHeight := pageProps.pageHeight
		footerSpace := pageProps.footerSpace
		currentX := pageProps.tableMarginX
		maxColHeight := getHighestCol(pdf, columnWidth, rows)
		currRowsheight := float64(maxColHeight) * lineHeight
		currPageRowHeight := currRowsheight
		curRowsBgHeightLeft := currRowsheight
		curRowsBgHeight := currPageRowHeight

		if pageProps.isNeedPageBreak(currentY+currRowsheight) && math.Floor((pageHeight-footerSpace-currentY)/lineHeight) > 0 {
			currPageRowHeight = math.Abs(math.Floor((pageHeight-footerSpace-currentY)/lineHeight) * lineHeight)
			curRowsBgHeight = currPageRowHeight
			curRowsBgHeightLeft = currRowsheight - curRowsBgHeight
		}

		pageProps.currRowsheight = currRowsheight
		pageProps.curRowsBgHeightLeft = curRowsBgHeightLeft
		pageProps.currPageRowHeight = currPageRowHeight
		pageProps.curRowsBgHeight = curRowsBgHeight
		pageProps.currRowsIndex = r
		pageProps.currentX = currentX
		pageProps.currpage = pdf.PageNo()

		if pageProps.isNeedPageBreak(currentY + lineHeight) {
			pdf.AddPage()
			pdf.SetPage(pageProps.currpage + 1)
			pageProps.currpage = pdf.PageNo()
			currentY = pageProps.headerSpace
			pageProps.currentY = currentY
		}

		pdf.SetX(currentX)
		pdf.SetY(currentY)
		if pageProps.newPageMargin > 0 {
			pdf.SetY(pageProps.newPageMargin)
			pageProps.currentY = pageProps.newPageMargin
			pageProps.newPageMargin = 0
		}
		pdf.SetPage(pageProps.currRowMaxPage)
		drawRows(pdf, pageProps, rows)

		// do calibration on the props
		currentY = pageProps.currentY
		curRowsBgHeight = pageProps.curRowsBgHeight
		currRowsheight = pageProps.currRowsheight

		//  add height for every line added
		pageProps.currentY = func() float64 {
			if pdf.PageNo() > 1 && currentY == pageProps.headerSpace && currRowsheight > curRowsBgHeight {
				return currentY + currRowsheight - curRowsBgHeight
			}
			return currentY + (float64(maxColHeight) * lineHeight)
		}()
	}
}

func drawRows(pdf *fpdf.Fpdf, pageProps *fpdfPageProperties, rows []string) {
	currentX := pageProps.tableMarginX
	currentY := pageProps.currentY
	columnWidth := pageProps.colWidthList
	lineHeight := pageProps.lineHeight
	pageHeight := pageProps.pageHeight
	// currPageRowHeight := pageProps.currPageRowHeight
	footerSpace := pageProps.footerSpace
	curRowsBgHeightLeft := pageProps.curRowsBgHeightLeft

	// curRowsBgHeight := pageProps.curRowsBgHeight
	curRowsBgHeight := func() float64 {
		if pageProps.isNeedPageBreak(currentY+curRowsBgHeightLeft) && math.Floor((pageHeight-footerSpace-currentY)/lineHeight) != 0 {
			currPageBg := math.Abs(math.Floor((pageHeight-footerSpace-currentY)/lineHeight) * lineHeight)
			return currPageBg
		}
		return math.Abs(curRowsBgHeightLeft)
	}()

	var rowPageOrigin int
	var rowYOrigin float64

	pageProps.currentX = currentX

	pdf.SetFontStyle("")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(240, 240, 240)

	if pageProps.currRowsIndex%2 != 0 {
		pdf.SetAlpha(0, "Normal")
	}

	// do not remove , i forgot why it was here
	// if pageProps.newPageMargin == 0 && pageProps.curRowsBgHeight < pageProps.currRowsheight {
	// 	curRowsBgHeight = pageProps.currRowsheight
	// }

	if curRowsBgHeight > pageProps.currRowsheight {
		curRowsBgHeight = pageProps.currRowsheight
	}
	pageProps.curRowsBgHeight = curRowsBgHeight
	pdf.Rect(currentX, currentY, float64(pageProps.tableWidth), curRowsBgHeight, "F")
	pdf.SetAlpha(1, "Normal")

	pageProps.totalColumn = len(rows) - 1
	for colNumber, col := range rows {
		if colNumber == 0 {
			rowPageOrigin = pageProps.currpage
			rowYOrigin = pageProps.currentY
		}

		pageProps.curColWidth = func() float64 {
			if colNumber > len(columnWidth)-1 {
				return 20
			}

			return columnWidth[colNumber]
		}()

		pageProps.curColIndex = colNumber

		drawCell(pdf, pageProps, col)

		//  reset coordinate when page breaks happen
		if pageProps.currpage != rowPageOrigin {
			pdf.SetPage(rowPageOrigin)
			pageProps.currpage = rowPageOrigin
		}
		if rowYOrigin != pageProps.currentY {
			pdf.SetY(rowYOrigin)
			pageProps.currentY = rowYOrigin
		}

	}
}

func drawCell(pdf *fpdf.Fpdf, pageProps *fpdfPageProperties, content string) {
	// detect page incosistent break
	var linePagePosLogs []int

	currColWidth := pageProps.curColWidth
	currentX := pageProps.currentX
	currentY := pageProps.currentY
	lineHeight := pageProps.lineHeight
	// curRowsBgHeight := pageProps.curRowsBgHeight
	// curRowsBgHeightLeft := pageProps.curRowsBgHeightLeft
	// colNumber := pageProps.curColIndex
	totalWidth := pageProps.tableWidth
	tableMarginX := pageProps.tableMarginX
	// footerSpace := pageProps.footerSpace
	// headerSpace := pageProps.headerSpace
	// pageHeight := pageProps.pageHeight
	// curRowsBgHeight := func() float64 {
	// 	if pageProps.isNeedPageBreak(curRowsBgHeightLeft + currentY) {
	// 		currPageBg := curRowsBgHeightLeft - (pageHeight - currentY - footerSpace)
	// 		pageProps.curRowsBgHeightLeft = curRowsBgHeightLeft - currPageBg
	// 		return currPageBg
	// 	}
	// 	return curRowsBgHeightLeft
	// }()

	pdf.SetY(currentY)
	pdf.SetX(currentX)

	// // column border
	// pdf.SetLineWidth(.25)
	// pdf.SetDrawColor(0, 0, 0)
	// pdf.SetAlpha(.25, "Normal")
	// pdf.Line(currentX, currentY, currentX, currentY+curRowsBgHeight)
	// if colNumber == pageProps.totalColumn {
	// 	pdf.Line(currentX+currColWidth, currentY, currentX+currColWidth, currentY+curRowsBgHeight)
	// }
	// pdf.SetAlpha(1, "Normal")
	splittedtext := pdf.SplitLines([]byte(content), currColWidth)
	lastRowY := currentY
	for lineIdx, text := range splittedtext {
		// reset properties when add page
		curRowsBgHeight := pageProps.curRowsBgHeightLeft
		if pageProps.isNeedPageBreak(lastRowY+float64(lineHeight)) && lineIdx > 0 {
			pdf.AddPage()
			pdf.SetPage(linePagePosLogs[lineIdx-1] + 1)
			curPagePosition := pdf.PageNo()
			pageProps.currpage = curPagePosition

			pageProps.currRowMaxPage = func() int {
				if curPagePosition <= pageProps.currRowMaxPage {
					return pageProps.currRowMaxPage
				}
				return curPagePosition
			}()
			curRowsBgHeight = pageProps.curRowsBgHeightLeft

			currentY = pageProps.headerSpace
			lastRowY = currentY

			if pageProps.isNeedPageBreak(lastRowY + curRowsBgHeight) {
				croppedBgHeight := curRowsBgHeight - (pageProps.pageHeight - pageProps.headerSpace - pageProps.footerSpace)
				pageProps.curRowsBgHeightLeft = curRowsBgHeight - croppedBgHeight
				curRowsBgHeight = croppedBgHeight
			}
			pageProps.logStat()
			pageProps.curRowsBgHeightLeft = curRowsBgHeight

			// draw bg
			if pageProps.currRowsIndex%2 != 0 {
				pdf.SetAlpha(0, "Normal")
			}

			pdf.Rect(tableMarginX, lastRowY, float64(totalWidth), curRowsBgHeight, "F")
			pdf.SetAlpha(1, "Normal")
			pageProps.currentY = lastRowY
			// mark if there is new page added
			pageProps.newPageMargin = lastRowY + curRowsBgHeight
		}

		linePagePosLogs = append(linePagePosLogs, pdf.PageNo())

		pdf.SetY(lastRowY)
		pdf.SetX(currentX)
		pdf.SetFont("Arial", "", 12)
		pdf.SetAlpha(1, "Normal")
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(currColWidth, lineHeight, string(text), "", 2, "C", false, 0, getLink(content))
		lastRowY += lineHeight
	}

	currentX += currColWidth
	pageProps.currentX = currentX
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
	drawFooter(pdf, &pageProps)

	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	pageWidth, pageHeight := pdf.GetPageSize()
	currentY := pageProps.headerHeight + 10
	lineHeight := float64(8)
	rowsWidth := pageWidth - pageProps.pageLeftPadding - pageProps.pageRightpadding
	pdf.SetLeftMargin(pageProps.pageLeftPadding)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetY(currentY)
	pdf.MultiCell(rowsWidth, lineHeight, fmt.Sprintf("Jenis informasi : %s", data.Information_Type), "", "", false)

	currentY += lineHeight
	pdf.SetY(currentY)
	pdf.SetTextColor(29, 31, 29)
	pdf.SetFont("Arial", "I", 12)
	pdf.CellFormat(pageWidth/4, lineHeight, fmt.Sprintf("Dibuat oleh : %s", data.Creator), "", 0, "", false, 0, "")

	pdf.CellFormat(pageWidth/4, lineHeight, fmt.Sprintf("Dibuat pada : %s", time.Unix(data.Effective_Date, 0).Format("15-06-2006")), "", 0, "", false, 0, "")

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

func drawFooter(pdf *fpdf.Fpdf, pageProps *fpdfPageProperties) {
	pageWidth, pageHeight := pdf.GetPageSize()
	footerHeight := 10
	footerMarginTop := 10
	pageProps.footerHeight = float64(footerHeight)
	pageProps.footerMarginTop = float64(footerMarginTop)
	pageProps.footerSpace = float64(footerHeight) + float64(footerMarginTop)

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
			return time.Now().Format("02/01/2006") + " - Page " + fmt.Sprintf("%v", pdf.PageNo())
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
	headerHeight := 30
	headerMarginBottom := 10
	pageProps.headerHeight = float64(headerHeight)
	pageProps.headerMarginBottom = float64(headerMarginBottom)
	pageProps.headerSpace = float64(headerHeight) + float64(headerMarginBottom)

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
		leftMargin := (pageWidth - headerTitleWidth) / 2
		pdf.SetX(leftMargin)
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
