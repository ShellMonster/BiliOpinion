import { useState } from 'react'
import html2canvas from 'html2canvas'
import * as XLSX from 'xlsx'
import type { ReportData } from '../../types/report'
import { useToast } from '../../hooks/useToast'

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
  const { showToast } = useToast()

  // 导出图片功能 (使用 html2canvas)
  const exportImage = async () => {
    const reportContainer = document.getElementById('report-container')
    if (!reportContainer) {
      showToast('未找到报告内容', 'error')
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
      showToast('图片导出成功', 'success')
    } catch (error) {
      console.error('导出图片失败:', error)
      showToast('导出图片失败', 'error')
    } finally {
      setImageExporting(false)
    }
  }

  // 导出Excel功能
  const exportExcel = () => {
    if (!reportData) {
      showToast('报告数据未加载', 'error')
      return
    }

    setExcelExporting(true)
    try {
      const wb = XLSX.utils.book_new()
      
      // 1. 品牌排名表
      const brandHeaders = ['排名', '品牌', '综合得分', ...reportData.dimensions.map(d => d.name)]
      const brandData = reportData.rankings.map(r => [
        r.rank,
        r.brand,
        r.overall_score.toFixed(1),
        ...reportData.dimensions.map(d => (r.scores[d.name] || 0).toFixed(1))
      ])
      const brandWs = XLSX.utils.aoa_to_sheet([brandHeaders, ...brandData])
      // 设置列宽
      const brandCols = [{ wch: 6 }, { wch: 15 }, { wch: 10 }, ...reportData.dimensions.map(() => ({ wch: 12 }))]
      brandWs['!cols'] = brandCols
      XLSX.utils.book_append_sheet(wb, brandWs, '品牌排名')
      
      // 2. 型号排名表
      if (reportData.model_rankings && reportData.model_rankings.length > 0) {
        const modelHeaders = ['排名', '型号', '品牌', '综合得分', '评论数', ...reportData.dimensions.map(d => d.name)]
        const modelData = reportData.model_rankings.map(r => [
          r.rank,
          r.model,
          r.brand,
          r.overall_score.toFixed(1),
          r.comment_count,
          ...reportData.dimensions.map(d => (r.scores[d.name] || 0).toFixed(1))
        ])
        const modelWs = XLSX.utils.aoa_to_sheet([modelHeaders, ...modelData])
        const modelCols = [{ wch: 6 }, { wch: 20 }, { wch: 15 }, { wch: 10 }, { wch: 8 }, ...reportData.dimensions.map(() => ({ wch: 12 }))]
        modelWs['!cols'] = modelCols
        XLSX.utils.book_append_sheet(wb, modelWs, '型号排名')
      }

      // 3. 维度说明
      const dimHeaders = ['维度名称', '维度说明']
      const dimData = reportData.dimensions.map(d => [d.name, d.description])
      const dimWs = XLSX.utils.aoa_to_sheet([dimHeaders, ...dimData])
      dimWs['!cols'] = [{ wch: 15 }, { wch: 50 }]
      XLSX.utils.book_append_sheet(wb, dimWs, '维度说明')

      // 4. 购买建议
      if (reportData.recommendation) {
        const recWs = XLSX.utils.aoa_to_sheet([['购买建议'], [reportData.recommendation]])
        recWs['!cols'] = [{ wch: 100 }]
        XLSX.utils.book_append_sheet(wb, recWs, '购买建议')
      }
      
      const dateStr = new Date().toISOString().split('T')[0]
      XLSX.writeFile(wb, `报告_${category}_${reportId}_${dateStr}.xlsx`)
      showToast('Excel导出成功', 'success')
    } catch (error) {
      console.error('导出Excel失败:', error)
      showToast('导出Excel失败', 'error')
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
          className="px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors font-medium text-sm cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
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
