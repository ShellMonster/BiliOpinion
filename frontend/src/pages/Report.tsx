import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
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
  const navigate = useNavigate()
  const { report, loading, error, id } = useReportData()
  const [activeTab, setActiveTab] = useState<TabType>('overview')
  const [exporting, setExporting] = useState(false)
  const [selectedBrand, setSelectedBrand] = useState<string | null>(null)
  const [selectedDims, setSelectedDims] = useState<string[]>([])
  const { showToast } = useToast()

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
        <ReportHeader category={data.category} reportId={id || ''} exporting={exporting} onExport={handleExportPDF} onBack={() => navigate('/history')} reportData={data} />
        
        <div className="flex space-x-2 border-b border-gray-200 pb-1">
          {tabs.map(t => <button key={t.key} onClick={() => setActiveTab(t.key as TabType)} className={`px-4 py-2 rounded-t-lg font-medium transition ${activeTab === t.key ? 'bg-white text-blue-600 border-b-2 border-blue-600' : 'text-gray-500 hover:text-gray-700'}`}>{t.label}</button>)}
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
