package pdf

import (
	"bilibili-analyzer/backend/report"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

func GeneratePDF(reportData *report.ReportData, reportID uint) ([]byte, error) {
	if reportData == nil {
		return nil, fmt.Errorf("reportData is nil")
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	fontFamily, useEnglish := loadFont(pdf)

	now := time.Now()

	pdf.SetTextColor(17, 24, 39)
	setFont(pdf, fontFamily, "B", 20)
	if useEnglish {
		pdf.CellFormat(0, 12, "Product Analysis Report", "", 1, "C", false, 0, "")
	} else {
		pdf.CellFormat(0, 12, "产品分析报告", "", 1, "C", false, 0, "")
	}
	pdf.Ln(2)

	setFont(pdf, fontFamily, "", 11)
	var meta string
	if useEnglish {
		meta = fmt.Sprintf("Category: %s | Report ID: %d | Generated: %s", safeText(reportData.Category), reportID, now.Format("2006-01-02 15:04"))
	} else {
		meta = fmt.Sprintf("类目: %s | 报告ID: %d | 生成时间: %s", safeText(reportData.Category), reportID, now.Format("2006-01-02 15:04"))
	}
	pdf.SetTextColor(55, 65, 81)
	pdf.CellFormat(0, 6, meta, "", 1, "L", false, 0, "")
	pdf.Ln(4)

	if useEnglish {
		sectionHeader(pdf, fontFamily, "Top Brand Overview")
	} else {
		sectionHeader(pdf, fontFamily, "Top 品牌概览")
	}
	drawTopBrandCard(pdf, fontFamily, reportData, useEnglish)
	pdf.Ln(4)

	if useEnglish {
		sectionHeader(pdf, fontFamily, "Brand Rankings")
	} else {
		sectionHeader(pdf, fontFamily, "品牌排名")
	}
	drawRankingTable(pdf, fontFamily, reportData, useEnglish)
	pdf.Ln(4)

	// 得分对比柱状图
	if useEnglish {
		sectionHeader(pdf, fontFamily, "Score Comparison")
	} else {
		sectionHeader(pdf, fontFamily, "得分对比")
	}
	drawBarChart(pdf, fontFamily, reportData, useEnglish)
	pdf.Ln(4)

	if useEnglish {
		sectionHeader(pdf, fontFamily, "Dimension Scores")
	} else {
		sectionHeader(pdf, fontFamily, "维度得分")
	}
	drawDimensionMatrix(pdf, fontFamily, reportData, useEnglish)
	pdf.Ln(4)

	if useEnglish {
		sectionHeader(pdf, fontFamily, "Purchase Recommendation")
	} else {
		sectionHeader(pdf, fontFamily, "购买建议")
	}
	drawRecommendation(pdf, fontFamily, reportData, useEnglish)

	if pdf.Error() != nil {
		return nil, fmt.Errorf("PDF generation failed: %w", pdf.Error())
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("PDF output failed: %w", err)
	}
	return buf.Bytes(), nil
}

func loadFont(pdf *fpdf.Fpdf) (string, bool) {
	// 1. 优先使用环境变量指定的字体
	if fontPath := strings.TrimSpace(os.Getenv("BILIBILI_PDF_FONT_PATH")); fontPath != "" {
		if data, err := os.ReadFile(fontPath); err == nil {
			pdf.AddUTF8FontFromBytes("CustomFont", "", data)
			pdf.AddUTF8FontFromBytes("CustomFont", "B", data)
			if pdf.Error() == nil {
				return "CustomFont", false
			}
			pdf.ClearError()
		}
	}

	// 2. 自动检测系统中文字体
	systemFonts := getSystemChineseFonts()
	for _, fontPath := range systemFonts {
		if data, err := os.ReadFile(fontPath); err == nil {
			pdf.AddUTF8FontFromBytes("SystemChinese", "", data)
			pdf.AddUTF8FontFromBytes("SystemChinese", "B", data)
			if pdf.Error() == nil {
				return "SystemChinese", false
			}
			pdf.ClearError()
		}
	}

	// 3. 回退到英文字体
	return "Arial", true
}

// getSystemChineseFonts 返回系统中文字体路径列表（按优先级排序）
func getSystemChineseFonts() []string {
	var fonts []string

	// macOS 字体路径
	macFonts := []string{
		"/System/Library/Fonts/PingFang.ttc",                   // 苹方（macOS 10.11+）
		"/System/Library/Fonts/STHeiti Light.ttc",              // 华文黑体
		"/System/Library/Fonts/STHeiti Medium.ttc",             // 华文黑体
		"/Library/Fonts/Arial Unicode.ttf",                     // Arial Unicode
		"/System/Library/Fonts/Supplemental/Songti.ttc",        // 宋体
		"/System/Library/Fonts/Supplemental/PingFang.ttc",      // 苹方补充
		"/System/Library/Fonts/Hiragino Sans GB.ttc",           // 冬青黑体
		"/System/Library/Fonts/Supplemental/Arial Unicode.ttf", // Arial Unicode 补充
	}
	fonts = append(fonts, macFonts...)

	// Windows 字体路径
	winFonts := []string{
		"C:\\Windows\\Fonts\\msyh.ttc",    // 微软雅黑
		"C:\\Windows\\Fonts\\msyhbd.ttc",  // 微软雅黑粗体
		"C:\\Windows\\Fonts\\simhei.ttf",  // 黑体
		"C:\\Windows\\Fonts\\simsun.ttc",  // 宋体
		"C:\\Windows\\Fonts\\simkai.ttf",  // 楷体
		"C:\\Windows\\Fonts\\STKAITI.TTF", // 华文楷体
		"C:\\Windows\\Fonts\\STSONG.TTF",  // 华文宋体
	}
	fonts = append(fonts, winFonts...)

	// Linux 字体路径
	linuxFonts := []string{
		"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",    // Noto Sans CJK
		"/usr/share/fonts/noto-cjk/NotoSansCJK-Regular.ttc",         // Noto Sans CJK (另一路径)
		"/usr/share/fonts/google-noto-cjk/NotoSansCJK-Regular.ttc",  // Noto Sans CJK (Fedora)
		"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",    // Noto Sans CJK (Ubuntu)
		"/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf", // Droid Sans Fallback
		"/usr/share/fonts/wenquanyi/wqy-microhei/wqy-microhei.ttc",  // 文泉驿微米黑
		"/usr/share/fonts/wenquanyi/wqy-zenhei/wqy-zenhei.ttc",      // 文泉驿正黑
		"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",            // 文泉驿微米黑 (另一路径)
		"/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc",              // 文泉驿正黑 (另一路径)
	}
	fonts = append(fonts, linuxFonts...)

	// 过滤存在的字体文件
	var existingFonts []string
	for _, f := range fonts {
		if _, err := os.Stat(f); err == nil {
			existingFonts = append(existingFonts, f)
		}
	}

	return existingFonts
}

func sectionHeader(pdf *fpdf.Fpdf, family, title string) {
	pdf.SetFillColor(238, 242, 255)
	pdf.SetDrawColor(203, 213, 225)
	pdf.SetTextColor(30, 41, 59)
	setFont(pdf, family, "B", 13)
	pdf.CellFormat(0, 8, title, "1", 1, "L", true, 0, "")
	pdf.Ln(2)
}

func drawTopBrandCard(pdf *fpdf.Fpdf, family string, reportData *report.ReportData, useEnglish bool) {
	left, _, right, _ := pdf.GetMargins()
	pageW, _ := pdf.GetPageSize()
	cardW := pageW - left - right
	cardH := 26.0

	var topBrand string
	var topScore float64
	if len(reportData.Rankings) > 0 {
		topBrand = reportData.Rankings[0].Brand
		topScore = reportData.Rankings[0].OverallScore
	}
	if topBrand == "" {
		if useEnglish {
			topBrand = "No Data"
		} else {
			topBrand = "暂无数据"
		}
	}

	x := left
	y := pdf.GetY()

	pdf.SetFillColor(240, 253, 244)
	pdf.SetDrawColor(187, 247, 208)
	pdf.RoundedRect(x, y, cardW, cardH, 2.5, "1234", "DF")

	pdf.SetXY(x+6, y+5)
	pdf.SetTextColor(6, 95, 70)
	setFont(pdf, family, "B", 12)
	pdf.CellFormat(0, 6, fmt.Sprintf("Top 1: %s", safeText(topBrand)), "", 1, "L", false, 0, "")

	pdf.SetX(x + 6)
	pdf.SetTextColor(17, 24, 39)
	setFont(pdf, family, "", 10.5)

	if len(reportData.Rankings) > 0 {
		if useEnglish {
			pdf.CellFormat(0, 6, fmt.Sprintf("Score: %.1f/10 | Brands: %d | Dimensions: %d", topScore, len(reportData.Brands), len(reportData.Dimensions)), "", 1, "L", false, 0, "")
		} else {
			pdf.CellFormat(0, 6, fmt.Sprintf("综合得分: %.1f/10 | 品牌数: %d | 维度数: %d", topScore, len(reportData.Brands), len(reportData.Dimensions)), "", 1, "L", false, 0, "")
		}
	} else {
		if useEnglish {
			pdf.CellFormat(0, 6, fmt.Sprintf("Brands: %d | Dimensions: %d", len(reportData.Brands), len(reportData.Dimensions)), "", 1, "L", false, 0, "")
		} else {
			pdf.CellFormat(0, 6, fmt.Sprintf("品牌数: %d | 维度数: %d", len(reportData.Brands), len(reportData.Dimensions)), "", 1, "L", false, 0, "")
		}
	}

	pdf.SetY(y + cardH)
}

func drawRankingTable(pdf *fpdf.Fpdf, family string, reportData *report.ReportData, useEnglish bool) {
	if len(reportData.Rankings) == 0 {
		setFont(pdf, family, "", 11)
		pdf.SetTextColor(75, 85, 99)
		if useEnglish {
			pdf.MultiCell(0, 6, "No ranking data available", "1", "L", false)
		} else {
			pdf.MultiCell(0, 6, "暂无排名数据", "1", "L", false)
		}
		return
	}

	left, _, right, _ := pdf.GetMargins()
	pageW, _ := pdf.GetPageSize()
	usableW := pageW - left - right

	setFont(pdf, family, "B", 11)
	pdf.SetFillColor(241, 245, 249)
	pdf.SetDrawColor(203, 213, 225)
	pdf.SetTextColor(15, 23, 42)

	colRank := 16.0
	colBrand := usableW*0.55 - colRank
	if colBrand < 40 {
		colBrand = 40
	}
	colScore := usableW - colRank - colBrand
	rowH := 8.0

	if useEnglish {
		pdf.CellFormat(colRank, rowH, "Rank", "1", 0, "C", true, 0, "")
		pdf.CellFormat(colBrand, rowH, "Brand", "1", 0, "L", true, 0, "")
		pdf.CellFormat(colScore, rowH, "Score", "1", 1, "C", true, 0, "")
	} else {
		pdf.CellFormat(colRank, rowH, "排名", "1", 0, "C", true, 0, "")
		pdf.CellFormat(colBrand, rowH, "品牌", "1", 0, "L", true, 0, "")
		pdf.CellFormat(colScore, rowH, "综合得分", "1", 1, "C", true, 0, "")
	}

	setFont(pdf, family, "", 11)
	for _, r := range reportData.Rankings {
		ensureSpace(pdf, rowH)

		pdf.SetTextColor(17, 24, 39)
		pdf.SetFillColor(255, 255, 255)
		pdf.CellFormat(colRank, rowH, fmt.Sprintf("%d", r.Rank), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colBrand, rowH, safeText(r.Brand), "1", 0, "L", false, 0, "")

		fillR, fillG, fillB, txtR, txtG, txtB := scoreColors(r.OverallScore)
		pdf.SetFillColor(fillR, fillG, fillB)
		pdf.SetTextColor(txtR, txtG, txtB)
		pdf.CellFormat(colScore, rowH, fmt.Sprintf("%.1f", r.OverallScore), "1", 1, "C", true, 0, "")
	}

	pdf.SetTextColor(17, 24, 39)
}

func drawDimensionMatrix(pdf *fpdf.Fpdf, family string, reportData *report.ReportData, useEnglish bool) {
	if len(reportData.Dimensions) == 0 {
		setFont(pdf, family, "", 11)
		pdf.SetTextColor(75, 85, 99)
		if useEnglish {
			pdf.MultiCell(0, 6, "No dimension data", "1", "L", false)
		} else {
			pdf.MultiCell(0, 6, "暂无维度数据", "1", "L", false)
		}
		return
	}
	if len(reportData.Brands) == 0 {
		setFont(pdf, family, "", 11)
		pdf.SetTextColor(75, 85, 99)
		if useEnglish {
			pdf.MultiCell(0, 6, "No brand data", "1", "L", false)
		} else {
			pdf.MultiCell(0, 6, "暂无品牌数据", "1", "L", false)
		}
		return
	}

	left, _, right, _ := pdf.GetMargins()
	pageW, _ := pdf.GetPageSize()
	usableW := pageW - left - right

	dimColW := 34.0
	minBrandColW := 22.0
	availForBrands := usableW - dimColW
	if availForBrands < minBrandColW {
		availForBrands = minBrandColW
	}

	maxCols := int(availForBrands / minBrandColW)
	if maxCols < 1 {
		maxCols = 1
	}

	for start := 0; start < len(reportData.Brands); start += maxCols {
		end := start + maxCols
		if end > len(reportData.Brands) {
			end = len(reportData.Brands)
		}
		brands := reportData.Brands[start:end]

		if len(reportData.Brands) > maxCols {
			setFont(pdf, family, "", 10)
			pdf.SetTextColor(75, 85, 99)
			if useEnglish {
				pdf.CellFormat(0, 5, fmt.Sprintf("Brands %d-%d of %d", start+1, end, len(reportData.Brands)), "", 1, "L", false, 0, "")
			} else {
				pdf.CellFormat(0, 5, fmt.Sprintf("品牌 %d-%d / %d", start+1, end, len(reportData.Brands)), "", 1, "L", false, 0, "")
			}
			pdf.Ln(1)
		}

		brandColW := (usableW - dimColW) / float64(len(brands))
		rowH := 8.0
		headH := 8.5

		ensureSpace(pdf, headH+rowH)
		setFont(pdf, family, "B", 10.5)
		pdf.SetFillColor(241, 245, 249)
		pdf.SetDrawColor(203, 213, 225)
		pdf.SetTextColor(15, 23, 42)

		if useEnglish {
			pdf.CellFormat(dimColW, headH, "Dimension", "1", 0, "C", true, 0, "")
		} else {
			pdf.CellFormat(dimColW, headH, "维度", "1", 0, "C", true, 0, "")
		}
		for _, b := range brands {
			pdf.CellFormat(brandColW, headH, safeText(b), "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)

		setFont(pdf, family, "", 10.5)
		for _, d := range reportData.Dimensions {
			ensureSpace(pdf, rowH)
			pdf.SetTextColor(17, 24, 39)
			pdf.SetFillColor(255, 255, 255)
			pdf.CellFormat(dimColW, rowH, safeText(d.Name), "1", 0, "L", false, 0, "")

			for _, b := range brands {
				score, ok := lookupScore(reportData, b, d.Name)
				if !ok {
					pdf.SetTextColor(107, 114, 128)
					pdf.SetFillColor(255, 255, 255)
					pdf.CellFormat(brandColW, rowH, "-", "1", 0, "C", false, 0, "")
					continue
				}

				fillR, fillG, fillB, txtR, txtG, txtB := scoreColors(score)
				pdf.SetFillColor(fillR, fillG, fillB)
				pdf.SetTextColor(txtR, txtG, txtB)
				pdf.CellFormat(brandColW, rowH, fmt.Sprintf("%.1f", score), "1", 0, "C", true, 0, "")
			}
			pdf.Ln(-1)
		}

		pdf.Ln(3)
	}

	pdf.SetTextColor(17, 24, 39)
}

func drawRecommendation(pdf *fpdf.Fpdf, family string, reportData *report.ReportData, useEnglish bool) {
	text := strings.TrimSpace(reportData.Recommendation)
	if text == "" {
		if useEnglish {
			text = "No recommendation available"
		} else {
			text = "暂无购买建议"
		}
	}

	setFont(pdf, family, "", 11)
	pdf.SetTextColor(17, 24, 39)
	pdf.SetFillColor(249, 250, 251)
	pdf.SetDrawColor(229, 231, 235)

	ensureSpace(pdf, 20)
	pdf.MultiCell(0, 6.5, text, "1", "L", true)
}

func lookupScore(reportData *report.ReportData, brand, dimName string) (float64, bool) {
	if reportData == nil || reportData.Scores == nil {
		return 0, false
	}
	brandScores, ok := reportData.Scores[brand]
	if !ok || brandScores == nil {
		return 0, false
	}
	s, ok := brandScores[dimName]
	return s, ok
}

func scoreColors(score float64) (fillR, fillG, fillB, txtR, txtG, txtB int) {
	if score >= 8.0 {
		return 220, 252, 231, 6, 95, 70
	}
	if score >= 6.0 {
		return 255, 251, 235, 146, 64, 14
	}
	return 254, 226, 226, 153, 27, 27
}

func ensureSpace(pdf *fpdf.Fpdf, needH float64) {
	_, pageH := pdf.GetPageSize()
	_, _, _, bottom := pdf.GetMargins()
	if pdf.GetY()+needH > pageH-bottom {
		pdf.AddPage()
	}
}

func setFont(pdf *fpdf.Fpdf, family, style string, size float64) {
	pdf.SetFont(family, style, size)
}

func safeText(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "-"
	}
	return s
}

// drawBarChart 绘制品牌综合得分横向柱状图
// 每个品牌一行，柱子长度按得分比例绘制，颜色根据得分等级区分
func drawBarChart(pdf *fpdf.Fpdf, family string, reportData *report.ReportData, useEnglish bool) {
	// 无数据时显示提示
	if len(reportData.Rankings) == 0 {
		setFont(pdf, family, "", 11)
		pdf.SetTextColor(75, 85, 99)
		if useEnglish {
			pdf.MultiCell(0, 6, "No ranking data available", "1", "L", false)
		} else {
			pdf.MultiCell(0, 6, "暂无排名数据", "1", "L", false)
		}
		return
	}

	// 获取页面尺寸和边距
	left, _, right, _ := pdf.GetMargins()
	pageW, _ := pdf.GetPageSize()
	usableW := pageW - left - right

	// 布局参数
	labelW := 45.0                           // 品牌名称列宽
	scoreW := 20.0                           // 得分数值列宽
	barMaxW := usableW - labelW - scoreW - 4 // 柱状图最大宽度（留4mm间距）
	rowH := 8.0                              // 每行高度
	barH := 6.0                              // 柱子高度（比行高略小，留上下间距）

	// 遍历所有品牌绘制柱状图
	for _, r := range reportData.Rankings {
		// 确保有足够空间，否则换页
		ensureSpace(pdf, rowH)

		y := pdf.GetY()
		x := left

		// 1. 绘制品牌名称（左侧）
		setFont(pdf, family, "", 10)
		pdf.SetTextColor(17, 24, 39)
		pdf.SetXY(x, y)
		pdf.CellFormat(labelW, rowH, safeText(r.Brand), "", 0, "L", false, 0, "")

		// 2. 绘制彩色柱状条（中间）
		// 计算柱子宽度：得分/10 * 最大宽度
		barW := (r.OverallScore / 10.0) * barMaxW
		if barW < 1 {
			barW = 1 // 最小宽度1mm，确保可见
		}

		// 获取得分对应的颜色
		fillR, fillG, fillB, _, _, _ := scoreColors(r.OverallScore)
		pdf.SetFillColor(fillR, fillG, fillB)

		// 绘制矩形条（垂直居中）
		barX := x + labelW + 2 // 留2mm间距
		barY := y + (rowH-barH)/2
		pdf.Rect(barX, barY, barW, barH, "F")

		// 3. 绘制得分数值（右侧）
		_, _, _, txtR, txtG, txtB := scoreColors(r.OverallScore)
		pdf.SetTextColor(txtR, txtG, txtB)
		setFont(pdf, family, "B", 10)
		pdf.SetXY(x+labelW+barMaxW+4, y)
		pdf.CellFormat(scoreW, rowH, fmt.Sprintf("%.1f", r.OverallScore), "", 0, "R", false, 0, "")

		// 移动到下一行
		pdf.SetY(y + rowH)
	}

	// 恢复默认文本颜色
	pdf.SetTextColor(17, 24, 39)
}
