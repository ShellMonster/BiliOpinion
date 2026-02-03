import { useState } from 'react'
import html2canvas from 'html2canvas'
import * as XLSX from 'xlsx'
import { useReportData } from '../hooks/useReportData'
import { useToast } from '../hooks/useToast'
import ReportHeader from '../components/Report/ReportHeader'
import { KeyStatsCards } from '../components/Report/Overview/KeyStatsCards'
import { BrandRadarChart } from '../components/Report/Charts/BrandRadarChart'
import { BrandScoreChart } from '../components/Report/Charts/BrandScoreChart'
import { RadarBrandSelector } from '../components/Report/Charts/RadarBrandSelector'
import { EnhancedSummary } from '../components/Report/EnhancedSummary'
import { BrandHeatmap } from '../components/Report/Charts/BrandHeatmap'
import { KeywordCloud } from '../components/Report/Charts/KeywordCloud'
import { SentimentPie } from '../components/Report/Charts/SentimentPie'
import { ScoreHistogram } from '../components/Report/Charts/ScoreHistogram'
import { BrandNetwork } from '../components/Report/Charts/BrandNetwork'
import { BrandCard } from '../components/Report/BrandCard'
import { DimensionFilter } from '../components/Report/DimensionFilter'
import { BrandDetailModal } from '../components/Report/BrandDetailModal'
import { CompetitorCompare } from '../components/Report/CompetitorCompare'
import { ModelAnalysis } from '../components/Report/ModelAnalysis'
import { DecisionTree } from '../components/Report/DecisionTree'
import type { SentimentStats } from '../types/report'

type TabType = 'overview' | 'charts' | 'summary'

const Report = () => {
  const { report, loading, error, id } = useReportData()
  const [activeTab, setActiveTab] = useState<TabType>('overview')
  const [exporting, setExporting] = useState(false)
  const [imageExporting, setImageExporting] = useState(false)
  const [excelExporting, setExcelExporting] = useState(false)
  const [selectedBrand, setSelectedBrand] = useState<string | null>(null)
  const [selectedDims, setSelectedDims] = useState<string[]>([])
  const { showToast } = useToast()

  // 导出图片功能
  const exportImage = async () => {
    const reportContainer = document.getElementById('report-container')
    const data = report?.data
    if (!reportContainer) {
      showToast('未找到报告内容', 'error')
      return
    }
    if (!data) return

    setImageExporting(true)
    try {
      const canvas = await html2canvas(reportContainer, {
        scale: 2,
        useCORS: true,
        allowTaint: true
      })
      
      const link = document.createElement('a')
      link.download = `报告_${data.category}_${id}.png`
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
    const data = report?.data
    if (!data) {
      showToast('报告数据未加载', 'error')
      return
    }

    setExcelExporting(true)
    try {
      const wb = XLSX.utils.book_new()
      
      // 1. 品牌排名表
      const brandHeaders = ['排名', '品牌', '综合得分', ...data.dimensions.map(d => d.name)]
      const brandData = data.rankings.map(r => [
        r.rank,
        r.brand,
        r.overall_score.toFixed(1),
        ...data.dimensions.map(d => (r.scores[d.name] || 0).toFixed(1))
      ])
      const brandWs = XLSX.utils.aoa_to_sheet([brandHeaders, ...brandData])
      // 设置列宽
      const brandCols = [{ wch: 6 }, { wch: 15 }, { wch: 10 }, ...data.dimensions.map(() => ({ wch: 12 }))]
      brandWs['!cols'] = brandCols
      XLSX.utils.book_append_sheet(wb, brandWs, '品牌排名')
      
      // 2. 型号排名表
      if (data.model_rankings && data.model_rankings.length > 0) {
        const modelHeaders = ['排名', '型号', '品牌', '综合得分', '评论数', ...data.dimensions.map(d => d.name)]
        const modelData = data.model_rankings.map(r => [
          r.rank,
          r.model,
          r.brand,
          r.overall_score.toFixed(1),
          r.comment_count,
          ...data.dimensions.map(d => (r.scores[d.name] || 0).toFixed(1))
        ])
        const modelWs = XLSX.utils.aoa_to_sheet([modelHeaders, ...modelData])
        const modelCols = [{ wch: 6 }, { wch: 20 }, { wch: 15 }, { wch: 10 }, { wch: 8 }, ...data.dimensions.map(() => ({ wch: 12 }))]
        modelWs['!cols'] = modelCols
        XLSX.utils.book_append_sheet(wb, modelWs, '型号排名')
      }

      // 3. 维度说明
      const dimHeaders = ['维度名称', '维度说明']
      const dimData = data.dimensions.map(d => [d.name, d.description])
      const dimWs = XLSX.utils.aoa_to_sheet([dimHeaders, ...dimData])
      dimWs['!cols'] = [{ wch: 15 }, { wch: 50 }]
      XLSX.utils.book_append_sheet(wb, dimWs, '维度说明')

      // 4. 购买建议
      if (data.recommendation) {
        const recWs = XLSX.utils.aoa_to_sheet([['购买建议'], [data.recommendation]])
        recWs['!cols'] = [{ wch: 100 }]
        XLSX.utils.book_append_sheet(wb, recWs, '购买建议')
      }
      
      const dateStr = new Date().toISOString().split('T')[0]
      XLSX.writeFile(wb, `报告_${data.category}_${id}_${dateStr}.xlsx`)
      showToast('Excel导出成功', 'success')
    } catch (error) {
      console.error('导出Excel失败:', error)
      showToast('导出Excel失败', 'error')
    } finally {
      setExcelExporting(false)
    }
  }

  const handleExportPDF = async () => {
    if (!id) return; setExporting(true)
    try {
      const response = await fetch(`http://localhost:8080/api/report/${id}/pdf`)
      if (!response.ok) throw new Error('导出失败')
      const blob = await response.blob(), url = window.URL.createObjectURL(blob)
      const a = document.createElement('a'); a.href = url; a.download = `报告_${report?.data.category}_${id}.pdf`
      document.body.appendChild(a); a.click(); document.body.removeChild(a); window.URL.revokeObjectURL(url)
      showToast('PDF导出成功', 'success')
    } catch (err) { showToast('导出失败，请重试', 'error') } finally { setExporting(false) }
  }

  if (loading) return <div className="min-h-screen flex items-center justify-center"><div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500" /></div>
  if (error || !report) return <div className="min-h-screen flex items-center justify-center text-red-500">{error || '报告不存在'}</div>

  const data = report.data
  const tabs = [{ key: 'overview', label: '总览' }, { key: 'charts', label: '图表' }, { key: 'summary', label: '深度总结' }]
  const currentDims = selectedDims.length ? selectedDims : data.dimensions.map(d => d.name)
  
  const totalSentiment = Object.values(data.sentiment_distribution || {}).reduce((acc, curr) => ({
    positive_count: acc.positive_count + curr.positive_count,
    neutral_count: acc.neutral_count + curr.neutral_count,
    negative_count: acc.negative_count + curr.negative_count,
    positive_pct: 0, neutral_pct: 0, negative_pct: 0 
  }), { positive_count: 0, neutral_count: 0, negative_count: 0, positive_pct: 0, neutral_pct: 0, negative_pct: 0 } as SentimentStats);

  return (
    <div className="min-h-screen bg-gray-50 py-8" id="report-container">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 space-y-6">
        <ReportHeader category={data.category} reportId={id || ''} />
        
        <div className="flex justify-between items-center border-b border-gray-200 pb-1">
          <div className="flex space-x-2">
            {tabs.map(t => <button key={t.key} onClick={() => setActiveTab(t.key as TabType)} className={`px-4 py-2 rounded-t-lg font-medium transition ${activeTab === t.key ? 'bg-white text-blue-600 border-b-2 border-blue-600' : 'text-gray-500 hover:text-gray-700'}`}>{t.label}</button>)}
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
              onClick={handleExportPDF}
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
          </div>
        </div>

        {activeTab === 'overview' && (
          <div className="space-y-6">
            <KeyStatsCards stats={data.stats || { total_videos: 0, total_comments: 0, comments_by_brand: {}}} brandCount={data.brands.length} />
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {data.rankings?.map(r => <BrandCard key={r.brand} ranking={r} analysis={data.brand_analysis?.[r.brand]} onClick={() => setSelectedBrand(r.brand)} />)}
            </div>
            <ModelAnalysis modelRankings={data.model_rankings || []} dimensions={data.dimensions} />
          </div>
        )}

        {activeTab === 'charts' && (
          <div className="space-y-6">
             <DimensionFilter dimensions={data.dimensions} selectedDimensions={currentDims} onChange={setSelectedDims} />
             <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
               <BrandRadarChart data={data} />
               <BrandScoreChart data={data} />
               <ScoreHistogram data={data.sentiment_distribution || {}} />
               <BrandHeatmap data={data} />
               <SentimentPie data={totalSentiment} title="整体情感分布" />
               <KeywordCloud data={data.keywords || []} />
             </div>
             <BrandNetwork data={data} />
             <RadarBrandSelector data={data} />
          </div>
        )}

        {activeTab === 'summary' && (
          <div className="space-y-6">
            <CompetitorCompare rankings={data.rankings} dimensions={data.dimensions} />
            <DecisionTree dimensions={data.dimensions} rankings={data.rankings} />
            <EnhancedSummary recommendation={data.recommendation} />
          </div>
        )}
      </div>
      {selectedBrand && (
        <BrandDetailModal
          isOpen={!!selectedBrand}
          onClose={() => setSelectedBrand(null)}
          brandName={selectedBrand}
          ranking={data.rankings?.find(r => r.brand === selectedBrand)}
          analysis={data.brand_analysis?.[selectedBrand]}
          topComments={data.top_comments?.[selectedBrand]}
          badComments={data.bad_comments?.[selectedBrand]}
          dimensions={data.dimensions}
        />
      )}
    </div>
  )
}

export default Report
