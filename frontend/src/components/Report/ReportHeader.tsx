import { useState } from 'react'
import html2canvas from 'html2canvas'
import * as XLSX from 'xlsx'
import type { ReportData } from '../../types/report'

interface ReportHeaderProps {
  category: string
  reportId: string | number
  exporting: boolean
  onExport: () => void
  onBack: () => void
  reportData?: ReportData
}

const ReportHeader = ({
  category,
  reportId,
  exporting,
  onExport,
  onBack,
  reportData
}: ReportHeaderProps) => {
  const [imageExporting, setImageExporting] = useState(false)
  const [excelExporting, setExcelExporting] = useState(false)

  // 导出图片功能
  const exportImage = async () => {
    const reportContainer = document.getElementById('report-container')
    if (!reportContainer) {
      alert('未找到报告内容')
      return
    }

    setImageExporting(true)
    try {
      const canvas = await html2canvas(reportContainer, {
        scale: 2,
        useCORS: true,
        allowTaint: true
      })
      
      const link = document.createElement('a')
      link.download = `报告_${category}_${reportId}.png`
      link.href = canvas.toDataURL('image/png')
      link.click()
    } catch (error) {
      console.error('导出图片失败:', error)
      alert('导出图片失败')
    } finally {
      setImageExporting(false)
    }
  }

  // 导出Excel功能
  const exportExcel = () => {
    if (!reportData) {
      alert('报告数据未加载')
      return
    }

    setExcelExporting(true)
    try {
      const wb = XLSX.utils.book_new()
      
      // 品牌排名工作表
      const brandData = reportData.rankings.map(r => ({
        '品牌': r.brand,
        '排名': r.rank,
        '综合得分': r.overall_score,
        ...r.scores
      }))
      const brandWs = XLSX.utils.json_to_sheet(brandData)
      XLSX.utils.book_append_sheet(wb, brandWs, '品牌排名')
      
      // 型号排名工作表
      if (reportData.model_rankings && reportData.model_rankings.length > 0) {
        const modelData = reportData.model_rankings.map(m => ({
          '型号': m.model,
          '品牌': m.brand,
          '排名': m.rank,
          '综合得分': m.overall_score,
          '评论数': m.comment_count,
          ...m.scores
        }))
        const modelWs = XLSX.utils.json_to_sheet(modelData)
        XLSX.utils.book_append_sheet(wb, modelWs, '型号排名')
      }
      
      XLSX.writeFile(wb, `报告_${category}_${reportId}.xlsx`)
    } catch (error) {
      console.error('导出Excel失败:', error)
      alert('导出Excel失败')
    } finally {
      setExcelExporting(false)
    }
  }

  return (
    <div className="flex justify-between items-end mb-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-800">产品分析报告</h1>
        <p className="text-gray-500">{category} | 报告ID: {reportId}</p>
      </div>
      <div className="flex gap-2">
        <button 
          onClick={exportImage}
          disabled={imageExporting}
          className="px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition-colors font-medium text-sm cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          {imageExporting ? (
            <>
              <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              导出中...
            </>
          ) : (
            <>
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              导出图片
            </>
          )}
        </button>
        <button 
          onClick={exportExcel}
          disabled={excelExporting}
          className="px-4 py-2 bg-orange-500 text-white rounded-lg hover:bg-orange-600 transition-colors font-medium text-sm cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          {excelExporting ? (
            <>
              <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              导出中...
            </>
          ) : (
            <>
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              导出Excel
            </>
          )}
        </button>
        <button 
          onClick={onExport}
          disabled={exporting}
          className="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors font-medium text-sm cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          {exporting ? (
            <>
              <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              导出中...
            </>
          ) : (
            <>
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              导出PDF
            </>
          )}
        </button>
        <button 
          onClick={onBack}
          className="px-4 py-2 bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 transition-colors font-medium text-sm cursor-pointer"
        >
          查看历史
        </button>
      </div>
    </div>
  )
}

export default ReportHeader
