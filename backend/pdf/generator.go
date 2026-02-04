package pdf

import (
	"bilibili-analyzer/backend/report"
	"bytes"
	"fmt"
	"log"
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

	if len(reportData.ModelRankings) > 0 {
		if useEnglish {
			sectionHeader(pdf, fontFamily, "Model Rankings")
		} else {
			sectionHeader(pdf, fontFamily, "型号排名")
		}
		drawModelRankingTable(pdf, fontFamily, reportData, useEnglish)
		pdf.Ln(4)
	}

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

	// 2. 使用项目内置的中文字体（跨平台兼容）
	builtinFontPaths := []string{
		"./fonts/siyuan.ttf",
		"../fonts/siyuan.ttf",
		"fonts/siyuan.ttf",
	}
	for _, fontPath := range builtinFontPaths {
		log.Printf("[PDF] 尝试加载内置字体: %s", fontPath)
		if data, err := os.ReadFile(fontPath); err == nil {
			pdf.AddUTF8FontFromBytes("SourceHanSans", "", data)
			pdf.AddUTF8FontFromBytes("SourceHanSans", "B", data)
			pdf.AddUTF8FontFromBytes("SourceHanSans", "I", data)
			pdf.AddUTF8FontFromBytes("SourceHanSans", "BI", data)
			if pdf.Error() == nil {
				log.Printf("[PDF] 内置字体加载成功: %s", fontPath)
				return "SourceHanSans", false
			}
			log.Printf("[PDF] 内置字体加载失败: %v", pdf.Error())
			pdf.ClearError()
		} else {
			log.Printf("[PDF] 读取内置字体失败: %s, 错误: %v", fontPath, err)
		}
	}

	// 3. 自动检测系统中文字体
	systemFonts := getSystemChineseFonts()
	boldFonts := getSystemBoldFonts()
	log.Printf("[PDF] 找到 %d 个可用系统字体，%d 个粗体字体", len(systemFonts), len(boldFonts))

	// 先尝试加载普通字体
	fontData := []byte{}
	fontLoaded := false
	for _, fontPath := range systemFonts {
		log.Printf("[PDF] 尝试加载字体: %s", fontPath)
		if data, err := os.ReadFile(fontPath); err == nil {
			fontData = data
			pdf.AddUTF8FontFromBytes("SystemChinese", "", data)
			if pdf.Error() == nil {
				log.Printf("[PDF] 普通字体加载成功: %s", fontPath)
				fontLoaded = true
				break
			}
			pdf.ClearError()
		}
	}

	// 尝试加载粗体字体（可能使用不同的字体文件）
	boldLoaded := false
	if len(boldFonts) > 0 && fontLoaded {
		for _, fontPath := range boldFonts {
			log.Printf("[PDF] 尝试加载粗体字体: %s", fontPath)
			if data, err := os.ReadFile(fontPath); err == nil {
				pdf.AddUTF8FontFromBytes("SystemChinese", "B", data)
				pdf.AddUTF8FontFromBytes("SystemChinese", "BI", data)
				if pdf.Error() == nil {
					log.Printf("[PDF] 粗体字体加载成功: %s", fontPath)
					boldLoaded = true
					break
				}
				pdf.ClearError()
			}
		}
	}

	// 如果粗体加载失败，使用普通字体模拟
	if fontLoaded && !boldLoaded && len(fontData) > 0 {
		pdf.AddUTF8FontFromBytes("SystemChinese", "B", fontData)
		pdf.AddUTF8FontFromBytes("SystemChinese", "I", fontData)
		pdf.AddUTF8FontFromBytes("SystemChinese", "BI", fontData)
		log.Printf("[PDF] 使用普通字体模拟粗体/斜体")
	}

	if fontLoaded {
		return "SystemChinese", false
	}

	log.Printf("[PDF] 警告: 所有中文字体加载失败，将使用Arial（可能导致中文乱码）")

	// 3. 回退到英文字体
	return "Arial", true
}

// getSystemBoldFonts 返回系统粗体中文字体路径列表
func getSystemBoldFonts() []string {
	var fonts []string

	// macOS 粗体字体
	macFonts := []string{
		"/System/Library/Fonts/STHeiti Medium.ttc",   // 华文黑体中粗
		"/System/Library/Fonts/Hiragino Sans GB.ttc", // 冬青黑体
	}
	fonts = append(fonts, macFonts...)

	// Windows 粗体字体
	winFonts := []string{
		"C:\\Windows\\Fonts\\msyhbd.ttc",  // 微软雅黑粗体
		"C:\\Windows\\Fonts\\simhei.ttf",  // 黑体（本身就是粗体）
	}
	fonts = append(fonts, winFonts...)

	// Linux 粗体字体
	linuxFonts := []string{
		"/usr/share/fonts/opentype/noto/NotoSansCJK-Bold.ttc",
		"/usr/share/fonts/noto-cjk/NotoSansCJK-Bold.ttc",
		"/usr/share/fonts/google-noto-cjk/NotoSansCJK-Bold.ttc",
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

// drawModelRankingTable 绘制型号排名表格
// 显示型号、品牌、综合得分和样本数
func drawModelRankingTable(pdf *fpdf.Fpdf, family string, reportData *report.ReportData, useEnglish bool) {
	if len(reportData.ModelRankings) == 0 {
		setFont(pdf, family, "", 11)
		pdf.SetTextColor(75, 85, 99)
		if useEnglish {
			pdf.MultiCell(0, 6, "No model ranking data available", "1", "L", false)
		} else {
			pdf.MultiCell(0, 6, "暂无型号排名数据", "1", "L", false)
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

	// 列宽分配：排名、型号、品牌、综合得分、样本数
	colRank := 16.0
	colModel := usableW * 0.30
	colBrand := usableW * 0.20
	colScore := usableW * 0.18
	colCount := usableW - colRank - colModel - colBrand - colScore
	rowH := 8.0

	// 表头
	if useEnglish {
		pdf.CellFormat(colRank, rowH, "Rank", "1", 0, "C", true, 0, "")
		pdf.CellFormat(colModel, rowH, "Model", "1", 0, "L", true, 0, "")
		pdf.CellFormat(colBrand, rowH, "Brand", "1", 0, "L", true, 0, "")
		pdf.CellFormat(colScore, rowH, "Score", "1", 0, "C", true, 0, "")
		pdf.CellFormat(colCount, rowH, "Samples", "1", 1, "C", true, 0, "")
	} else {
		pdf.CellFormat(colRank, rowH, "排名", "1", 0, "C", true, 0, "")
		pdf.CellFormat(colModel, rowH, "型号", "1", 0, "L", true, 0, "")
		pdf.CellFormat(colBrand, rowH, "品牌", "1", 0, "L", true, 0, "")
		pdf.CellFormat(colScore, rowH, "综合得分", "1", 0, "C", true, 0, "")
		pdf.CellFormat(colCount, rowH, "样本数", "1", 1, "C", true, 0, "")
	}

	// 数据行
	setFont(pdf, family, "", 11)
	for _, r := range reportData.ModelRankings {
		ensureSpace(pdf, rowH)

		pdf.SetTextColor(17, 24, 39)
		pdf.SetFillColor(255, 255, 255)
		pdf.CellFormat(colRank, rowH, fmt.Sprintf("%d", r.Rank), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colModel, rowH, safeText(r.Model), "1", 0, "L", false, 0, "")
		pdf.CellFormat(colBrand, rowH, safeText(r.Brand), "1", 0, "L", false, 0, "")

		// 得分列使用颜色编码
		fillR, fillG, fillB, txtR, txtG, txtB := scoreColors(r.OverallScore)
		pdf.SetFillColor(fillR, fillG, fillB)
		pdf.SetTextColor(txtR, txtG, txtB)
		pdf.CellFormat(colScore, rowH, fmt.Sprintf("%.1f", r.OverallScore), "1", 0, "C", true, 0, "")

		// 样本数列
		pdf.SetTextColor(17, 24, 39)
		pdf.SetFillColor(255, 255, 255)
		pdf.CellFormat(colCount, rowH, fmt.Sprintf("%d", r.CommentCount), "1", 1, "C", false, 0, "")
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

	pdf.SetTextColor(17, 24, 39)
	pdf.SetFillColor(249, 250, 251)
	pdf.SetDrawColor(229, 231, 235)

	ensureSpace(pdf, 20)

	// 获取页面边距和宽度
	left, _, right, _ := pdf.GetMargins()
	pageW, _ := pdf.GetPageSize()
	usableW := pageW - left - right

	// 第一遍：解析所有块，计算总高度
	blocks := parseMarkdownBlocks(text)
	totalH := 0.0
	for _, block := range blocks {
		totalH += estimateBlockHeight(pdf, family, block, usableW-4)
	}
	// 加上内边距
	totalH += 4
	if totalH < 20 {
		totalH = 20
	}

	// 绘制背景框
	y := pdf.GetY()
	pdf.Rect(left, y, usableW, totalH, "FD")

	// 第二遍：渲染所有块
	pdf.SetXY(left+2, y+2) // 添加内边距
	for _, block := range blocks {
		drawMarkdownBlock(pdf, family, block, usableW-4)
	}
}

// markdownBlock 表示 Markdown 块级元素
type markdownBlock struct {
	Type  string // "heading", "paragraph", "list", "quote", "empty"
	Text  string
	Level int     // 标题级别
}

// parseMarkdownBlocks 解析 Markdown 文本为块级元素
func parseMarkdownBlocks(text string) []markdownBlock {
	lines := strings.Split(text, "\n")
	var blocks []markdownBlock
	var currentList []string
	var currentQuote []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 空行
		if trimmed == "" {
			if len(currentList) > 0 {
				blocks = append(blocks, markdownBlock{
					Type: "list",
					Text: strings.Join(currentList, "\n"),
				})
				currentList = nil
			}
			if len(currentQuote) > 0 {
				blocks = append(blocks, markdownBlock{
					Type: "quote",
					Text: strings.Join(currentQuote, "\n"),
				})
				currentQuote = nil
			}
			blocks = append(blocks, markdownBlock{Type: "empty"})
			continue
		}

		// 标题 ## 或 ###
		if strings.HasPrefix(trimmed, "##") {
			if len(currentList) > 0 {
				blocks = append(blocks, markdownBlock{
					Type: "list",
					Text: strings.Join(currentList, "\n"),
				})
				currentList = nil
			}
			if len(currentQuote) > 0 {
				blocks = append(blocks, markdownBlock{
					Type: "quote",
					Text: strings.Join(currentQuote, "\n"),
				})
				currentQuote = nil
			}
			level := 0
			for _, c := range trimmed {
				if c == '#' {
					level++
				} else {
					break
				}
			}
			content := strings.TrimSpace(trimmed[level:])
			blocks = append(blocks, markdownBlock{
				Type:  "heading",
				Text:  content,
				Level: level,
			})
			continue
		}

		// 引用 >
		if strings.HasPrefix(trimmed, ">") {
			if len(currentList) > 0 {
				blocks = append(blocks, markdownBlock{
					Type: "list",
					Text: strings.Join(currentList, "\n"),
				})
				currentList = nil
			}
			content := strings.TrimSpace(trimmed[1:])
			currentQuote = append(currentQuote, content)
			continue
		}

		// 列表 -
		if strings.HasPrefix(trimmed, "-") {
			if len(currentQuote) > 0 {
				blocks = append(blocks, markdownBlock{
					Type: "quote",
					Text: strings.Join(currentQuote, "\n"),
				})
				currentQuote = nil
			}
			content := strings.TrimSpace(trimmed[1:])
			currentList = append(currentList, content)
			continue
		}

		// 普通段落
		if len(currentList) > 0 {
			blocks = append(blocks, markdownBlock{
				Type: "list",
				Text: strings.Join(currentList, "\n"),
			})
			currentList = nil
		}
		if len(currentQuote) > 0 {
			blocks = append(blocks, markdownBlock{
				Type: "quote",
				Text: strings.Join(currentQuote, "\n"),
			})
			currentQuote = nil
		}
		blocks = append(blocks, markdownBlock{
			Type: "paragraph",
			Text: trimmed,
		})
	}

	// 处理最后的块
	if len(currentList) > 0 {
		blocks = append(blocks, markdownBlock{
			Type: "list",
			Text: strings.Join(currentList, "\n"),
		})
	}
	if len(currentQuote) > 0 {
		blocks = append(blocks, markdownBlock{
			Type: "quote",
			Text: strings.Join(currentQuote, "\n"),
		})
	}

	return blocks
}

// estimateBlockHeight 估算块的高度
func estimateBlockHeight(pdf *fpdf.Fpdf, family string, block markdownBlock, maxW float64) float64 {
	setFont(pdf, family, "", 11)
	switch block.Type {
	case "heading":
		setFont(pdf, family, "B", 13)
		// 估算标题行数
		textW := pdf.GetStringWidth(block.Text)
		lines := int(textW/maxW) + 1
		return float64(lines) * 7 + 2
	case "paragraph":
		textW := pdf.GetStringWidth(block.Text)
		lines := int(textW/maxW) + 1
		return float64(lines) * 6.5 + 1
	case "list":
		items := strings.Count(block.Text, "\n") + 1
		return float64(items) * 6.5 + 1
	case "quote":
		textW := pdf.GetStringWidth(block.Text)
		lines := int(textW/(maxW-4)) + 1
		return float64(lines) * 6 + 2
	case "empty":
		return 3
	default:
		return 6.5
	}
}

// drawMarkdownBlock 绘制 Markdown 块
func drawMarkdownBlock(pdf *fpdf.Fpdf, family string, block markdownBlock, maxW float64) {
	left, _, _, _ := pdf.GetMargins()

	switch block.Type {
	case "empty":
		pdf.Ln(3)

	case "heading":
		pdf.SetTextColor(30, 41, 59)
		setFont(pdf, family, "B", 13)
		drawStyledText(pdf, family, block.Text, maxW, 7, 0)
		pdf.Ln(1)

	case "paragraph":
		pdf.SetTextColor(17, 24, 39)
		setFont(pdf, family, "", 11)
		drawStyledText(pdf, family, block.Text, maxW, 6.5, 2)
		pdf.Ln(0.5)

	case "list":
		pdf.SetTextColor(17, 24, 39)
		setFont(pdf, family, "", 11)
		items := strings.Split(block.Text, "\n")
		for _, item := range items {
			drawStyledText(pdf, family, "• "+item, maxW, 6.5, 4)
			pdf.Ln(0.3)
		}

	case "quote":
		// 绘制引用背景
		x := pdf.GetX()
		y := pdf.GetY()
		pdf.SetFillColor(243, 244, 246)
		pdf.Rect(left+2, y, maxW, 6, "FD")
		pdf.SetXY(x+4, y+1)

		pdf.SetTextColor(71, 85, 105)
		setFont(pdf, family, "I", 10)
		drawStyledText(pdf, family, block.Text, maxW-8, 5.5, 2)
		pdf.Ln(1)
	}
}

// drawStyledText 解析并渲染带有 Markdown 样式的文本
// 支持 **bold**、*italic* 和混合样式
func drawStyledText(pdf *fpdf.Fpdf, family, text string, maxW, lineH, indent float64) {
	// 获取左边距，用于换行时重置X坐标
	left, _, _, _ := pdf.GetMargins()

	x := pdf.GetX() + indent
	y := pdf.GetY()
	pdf.SetXY(x, y)

	// 解析文本为样式段
	segments := parseMarkdown(text)

	currentX := x
	for _, seg := range segments {
		// 先设置字体样式，再计算宽度（确保粗体宽度准确）
		setFont(pdf, family, seg.style, 11)
		textW := pdf.GetStringWidth(seg.text)

		// 检查是否需要换行
		if currentX+textW > x+maxW && currentX > x {
			pdf.Ln(lineH)
			// 修复：重置到左边距+缩进，而不是初始位置x
			currentX = left + indent
			pdf.SetX(currentX)
		}

		// 绘制文本段
		pdf.CellFormat(textW, lineH, seg.text, "", 0, "L", false, 0, "")
		currentX += textW
	}

	// 修复：确保结束后换到新行，避免下一次调用时位置错乱
	pdf.Ln(lineH)
}

// textSegment 表示一个带样式的文本段
type textSegment struct {
	text  string // 文本内容
	style string // 样式: "" = normal, "B" = bold, "I" = italic, "BI" = bold+italic
}

// parseMarkdown 解析 Markdown 文本为样式段列表
// 支持 **bold**、*italic* 和嵌套样式
func parseMarkdown(text string) []textSegment {
	var segments []textSegment
	var current strings.Builder
	currentStyle := ""

	i := 0
	for i < len(text) {
		// 检查 **bold**
		if i+1 < len(text) && text[i:i+2] == "**" {
			// 保存当前段
			if current.Len() > 0 {
				segments = append(segments, textSegment{text: current.String(), style: currentStyle})
				current.Reset()
			}

			// 查找结束标记
			end := strings.Index(text[i+2:], "**")
			if end != -1 {
				end += i + 2
				boldText := text[i+2 : end]

				// 检查是否包含 *italic*
				if strings.Contains(boldText, "*") {
					// 递归解析嵌套样式
					nestedSegs := parseMarkdown(boldText)
					for _, seg := range nestedSegs {
						if seg.style == "" {
							seg.style = "B"
						} else if seg.style == "I" {
							seg.style = "BI"
						}
						segments = append(segments, seg)
					}
				} else {
					segments = append(segments, textSegment{text: boldText, style: "B"})
				}

				i = end + 2
				continue
			}
		}

		// 检查 *italic*
		if text[i] == '*' && (i == 0 || text[i-1] != '*') {
			// 保存当前段
			if current.Len() > 0 {
				segments = append(segments, textSegment{text: current.String(), style: currentStyle})
				current.Reset()
			}

			// 查找结束标记
			end := i + 1
			for end < len(text) && text[end] != '*' {
				end++
			}
			if end < len(text) {
				italicText := text[i+1 : end]
				segments = append(segments, textSegment{text: italicText, style: "I"})
				i = end + 1
				continue
			}
		}

		// 普通字符
		current.WriteByte(text[i])
		i++
	}

	// 保存最后一段
	if current.Len() > 0 {
		segments = append(segments, textSegment{text: current.String(), style: currentStyle})
	}

	return segments
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
	// 尝试设置字体，如果失败则使用Arial fallback
	pdf.SetFont(family, style, size)
	if pdf.Error() != nil {
		pdf.ClearError()
		// 如果是有样式的字体失败，尝试无样式版本+模拟粗体
		if style != "" {
			pdf.SetFont(family, "", size)
			if pdf.Error() != nil {
				pdf.ClearError()
				pdf.SetFont("Arial", style, size)
				if pdf.Error() != nil {
					pdf.ClearError()
					pdf.SetFont("Arial", "", size)
				}
			} else {
				// 如果中文字体不支持粗体样式，使用模拟粗体（增加字重）
				if strings.Contains(style, "B") {
					// 设置文本渲染模式为填充来实现粗体效果
					// 注意：fpdf 不直接支持模拟粗体，这里通过字体大小微调实现视觉效果
					// 或者可以接受当前无粗体状态，因为中文字体通常字重已经足够
				}
			}
		} else {
			pdf.SetFont("Arial", "", size)
		}
	}
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
