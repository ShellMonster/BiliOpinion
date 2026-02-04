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
import { BrandNetwork } from '../components/Report/Charts/BrandNetwork'
import { BrandCard } from '../components/Report/BrandCard'
import { DimensionFilter } from '../components/Report/DimensionFilter'
import { BrandDetailModal } from '../components/Report/BrandDetailModal'
import { CompetitorCompare } from '../components/Report/CompetitorCompare'
import { DecisionTree } from '../components/Report/DecisionTree'
import { VideoSourceList } from '../components/Report/VideoSourceList'
import type { SentimentStats, ModelRanking } from '../types/report'

type TabType = 'overview' | 'charts' | 'summary' | 'sources'

const Report = () => {
  const { report, loading, error, id } = useReportData()
  const [activeTab, setActiveTab] = useState<TabType>('overview')
  const [exporting, setExporting] = useState(false)
  const [imageExporting, setImageExporting] = useState(false)
  const [excelExporting, setExcelExporting] = useState(false)
  const [selectedBrand, setSelectedBrand] = useState<string | null>(null)
  const [selectedDims, setSelectedDims] = useState<string[]>([])
  const [hideUnknown, setHideUnknown] = useState(true)
  const [hideZeroScore, setHideZeroScore] = useState(true)
  const { showToast } = useToast()

  // 导出图片功能
  const exportImage = async () => {
    const reportContainer = document.getElementById('report-container')
    const data = report?.data
    if (!reportContainer) {
      showToast('未找到报告内容', 'error')
      return
    }
    if (!data) {
      showToast('报告数据未加载', 'error')
      return
    }

    console.log('[Image] 开始导出，容器大小:', reportContainer.scrollWidth, 'x', reportContainer.scrollHeight)

    // 临时转换oklch颜色为rgb
    const elementsWithOklch: {el: HTMLElement, originalColor: string, originalBg: string}[] = []
    const allElements = reportContainer.querySelectorAll('*')
    
    allElements.forEach(el => {
      const htmlEl = el as HTMLElement
      const computedStyle = window.getComputedStyle(htmlEl)
      const color = computedStyle.color
      const bgColor = computedStyle.backgroundColor
      
      if (color.includes('oklch') || bgColor.includes('oklch')) {
        elementsWithOklch.push({
          el: htmlEl,
          originalColor: htmlEl.style.color,
          originalBg: htmlEl.style.backgroundColor
        })
        // 设置为安全的颜色
        if (color.includes('oklch')) {
          htmlEl.style.color = '#374151'
        }
        if (bgColor.includes('oklch')) {
          htmlEl.style.backgroundColor = '#ffffff'
        }
      }
    })
    
    console.log('[Image] 转换了', elementsWithOklch.length, '个oklch颜色元素')

    setImageExporting(true)
    try {
      const canvas = await html2canvas(reportContainer, {
        scale: 2,
        useCORS: true,
        allowTaint: true,
        logging: true,
        backgroundColor: '#ffffff',
      })
      
      console.log('[Image] Canvas生成成功:', canvas.width, 'x', canvas.height)
      
      const link = document.createElement('a')
      link.download = `报告_${data.category}_${id}.png`
      link.href = canvas.toDataURL('image/png')
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      
      console.log('[Image] 导出成功')
      showToast('图片导出成功', 'success')
    } catch (error) {
      console.error('[Image] 导出失败:', error)
      showToast(`导出图片失败: ${error instanceof Error ? error.message : '未知错误'}`, 'error')
    } finally {
      // 恢复原始颜色
      elementsWithOklch.forEach(({el, originalColor, originalBg}) => {
        el.style.color = originalColor
        el.style.backgroundColor = originalBg
      })
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

    console.log('[Excel] 开始导出，数据量:', {
      rankings: data.rankings?.length,
      modelRankings: data.model_rankings?.length,
      dimensions: data.dimensions?.length
    })

    setExcelExporting(true)
    try {
      const wb = XLSX.utils.book_new()
      
      // 1. 品牌排名表
      console.log('[Excel] 创建品牌排名表')
      const brandHeaders = ['排名', '品牌', '综合得分', ...data.dimensions.map(d => d.name)]
      const brandData = data.rankings.map(r => [
        r.rank,
        r.brand,
        r.overall_score.toFixed(1),
        ...data.dimensions.map(d => (r.scores[d.name] || 0).toFixed(1))
      ])
      const brandWs = XLSX.utils.aoa_to_sheet([brandHeaders, ...brandData])
      const brandCols = [{ wch: 6 }, { wch: 15 }, { wch: 10 }, ...data.dimensions.map(() => ({ wch: 12 }))]
      brandWs['!cols'] = brandCols
      XLSX.utils.book_append_sheet(wb, brandWs, '品牌排名')
      
      // 2. 型号排名表
      if (data.model_rankings && data.model_rankings.length > 0) {
        console.log('[Excel] 创建型号排名表')
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
      console.log('[Excel] 创建维度说明表')
      const dimHeaders = ['维度名称', '维度说明']
      const dimData = data.dimensions.map(d => [d.name, d.description])
      const dimWs = XLSX.utils.aoa_to_sheet([dimHeaders, ...dimData])
      dimWs['!cols'] = [{ wch: 15 }, { wch: 50 }]
      XLSX.utils.book_append_sheet(wb, dimWs, '维度说明')

      // 4. 购买建议
      if (data.recommendation) {
        console.log('[Excel] 创建购买建议表')
        const recWs = XLSX.utils.aoa_to_sheet([['购买建议'], [data.recommendation]])
        recWs['!cols'] = [{ wch: 100 }]
        XLSX.utils.book_append_sheet(wb, recWs, '购买建议')
      }
      
      const dateStr = new Date().toISOString().split('T')[0]
      const filename = `报告_${data.category}_${id}_${dateStr}.xlsx`
      console.log('[Excel] 保存文件:', filename)
      XLSX.writeFile(wb, filename)
      console.log('[Excel] 导出成功')
      showToast('Excel导出成功', 'success')
    } catch (error) {
      console.error('[Excel] 导出失败:', error)
      showToast(`导出Excel失败: ${error instanceof Error ? error.message : '未知错误'}`, 'error')
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
  const tabs = [
    { key: 'overview', label: '总览' },
    { key: 'charts', label: '图表' },
    { key: 'summary', label: '深度总结' },
    { key: 'sources', label: '数据来源' }
  ]
  const currentDims = selectedDims.length ? selectedDims : data.dimensions.map(d => d.name)
  
  // 过滤后的数据
  const filteredRankings = data.rankings?.filter(r => {
    if (hideUnknown && r.brand === '未知') return false
    if (hideZeroScore && r.overall_score === 0) return false
    return true
  })

  const filteredModelRankings = data.model_rankings?.filter(m => {
    if (hideUnknown && (m.brand === '未知' || m.model === '通用')) return false
    if (hideZeroScore && m.overall_score === 0) return false
    return true
  })
  
  // 情感分布数据 - 直接使用整体统计，不需要reduce计算
  const totalSentiment: SentimentStats = data.sentiment_distribution || {
    positive_count: 0,
    neutral_count: 0,
    negative_count: 0,
    positive_pct: 0,
    neutral_pct: 0,
    negative_pct: 0
  };

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
            
            {/* 过滤开关 */}
            <div className="flex items-center gap-4 bg-white rounded-lg px-4 py-2 shadow-sm border border-gray-200">
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="hideUnknown"
                  checked={hideUnknown}
                  onChange={(e) => setHideUnknown(e.target.checked)}
                  className="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
                />
                <label htmlFor="hideUnknown" className="text-sm text-gray-700 cursor-pointer">
                  隐藏"未知"品牌和"通用"型号
                </label>
              </div>
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="hideZeroScore"
                  checked={hideZeroScore}
                  onChange={(e) => setHideZeroScore(e.target.checked)}
                  className="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
                />
                <label htmlFor="hideZeroScore" className="text-sm text-gray-700 cursor-pointer">
                  隐藏零分数据
                </label>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredRankings && filteredRankings.length > 0 ? (
                filteredRankings.map(r => <BrandCard key={r.brand} ranking={r} analysis={data.brand_analysis?.[r.brand]} onClick={() => setSelectedBrand(r.brand)} />)
              ) : (
                <div className="col-span-full text-center py-8 text-gray-500">
                  没有符合条件的品牌数据
                </div>
              )}
            </div>
            
            {/* 型号排名 */}
            {filteredModelRankings && filteredModelRankings.length > 0 && (
              <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                <div className="px-6 py-4 border-b border-gray-200">
                  <h2 className="text-xl font-bold text-gray-800">🏆 型号排名</h2>
                  <p className="text-sm text-gray-500 mt-1">基于 AI 分析的具体型号表现</p>
                </div>
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">排名</th>
                        <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">型号</th>
                        <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">品牌</th>
                        <th className="px-4 py-3 text-center text-sm font-medium text-gray-600">综合得分</th>
                        {data.dimensions.map(dim => (
                          <th key={dim.name} className="px-4 py-3 text-center text-sm font-medium text-gray-600">{dim.name}</th>
                        ))}
                        <th className="px-4 py-3 text-center text-sm font-medium text-gray-600">评论数</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                      {filteredModelRankings.map((model: ModelRanking) => (
                        <tr key={`${model.brand}-${model.model}`} className="hover:bg-gray-50">
                          <td className="px-4 py-3">
                            <span className={`inline-flex items-center justify-center w-8 h-8 rounded-full text-sm font-bold ${
                              model.rank === 1 ? 'bg-yellow-100 text-yellow-700' :
                              model.rank === 2 ? 'bg-gray-100 text-gray-700' :
                              model.rank === 3 ? 'bg-orange-100 text-orange-700' :
                              'bg-blue-50 text-blue-600'
                            }`}>
                              {model.rank}
                            </span>
                          </td>
                          <td className="px-4 py-3 font-medium text-gray-900">{model.model}</td>
                          <td className="px-4 py-3 text-gray-600">{model.brand}</td>
                          <td className="px-4 py-3 text-center">
                            <span className={`inline-flex items-center px-2 py-1 rounded text-sm font-medium ${
                              model.overall_score >= 8 ? 'bg-green-100 text-green-700' :
                              model.overall_score >= 6 ? 'bg-blue-100 text-blue-700' :
                              'bg-amber-100 text-amber-700'
                            }`}>
                              {model.overall_score.toFixed(1)}
                            </span>
                          </td>
                          {data.dimensions.map(dim => {
                            const score = model.scores?.[dim.name]
                            return (
                              <td key={dim.name} className="px-4 py-3 text-center text-sm text-gray-600">
                                {score !== undefined ? score.toFixed(1) : '-'}
                              </td>
                            )
                          })}
                          <td className="px-4 py-3 text-center text-sm text-gray-500">{model.comment_count}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
            
            {filteredModelRankings && filteredModelRankings.length === 0 && data.model_rankings && data.model_rankings.length > 0 && (
              <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
                <div className="text-center text-gray-500">
                  没有符合条件的型号数据
                </div>
              </div>
            )}
          </div>
        )}

        {activeTab === 'charts' && (
          <div className="space-y-6">
             <DimensionFilter dimensions={data.dimensions} selectedDimensions={currentDims} onChange={setSelectedDims} />
             <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
               <BrandRadarChart data={data} />
               <BrandScoreChart data={data} />
               <BrandHeatmap data={data} />
               <SentimentPie data={totalSentiment} title="整体情感分布" />
               <KeywordCloud data={data.keyword_frequency || []} />
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

        {activeTab === 'sources' && (
          <div className="space-y-6">
            {data.video_sources && data.video_sources.length > 0 ? (
              <VideoSourceList videos={data.video_sources} />
            ) : (
              <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
                <div className="text-center text-gray-500">
                  暂无数据来源信息
                </div>
              </div>
            )}
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
