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

type TabType = 'overview' | 'charts' | 'summary'

const Report = () => {
  const navigate = useNavigate()
  const { report, loading, error, id } = useReportData()
  const [activeTab, setActiveTab] = useState<TabType>('overview')
  const [exporting, setExporting] = useState(false)
  const { showToast } = useToast()

  const handleExportPDF = async () => {
    if (!id) return
    setExporting(true)
    try {
      const response = await fetch(`http://localhost:8080/api/report/${id}/pdf`)
      if (!response.ok) throw new Error('导出失败')
      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `报告_${report?.data.category || '产品'}_${id}.pdf`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
      showToast('PDF导出成功', 'success')
    } catch (err) {
      showToast('导出失败，请重试', 'error')
    } finally {
      setExporting(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-600">加载中...</p>
        </div>
      </div>
    )
  }

  if (error || !report) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-500 mb-4">{error || '报告不存在'}</p>
          <button 
            onClick={() => navigate('/history')}
            className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
          >
            返回历史记录
          </button>
        </div>
      </div>
    )
  }

  const data = report.data

  const tabs: { key: TabType; label: string }[] = [
    { key: 'overview', label: '总览' },
    { key: 'charts', label: '图表' },
    { key: 'summary', label: '总结' },
  ]

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div id="report-container" className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <ReportHeader 
          category={data.category}
          reportId={id || ''}
          exporting={exporting}
          onExport={handleExportPDF}
          onBack={() => navigate('/history')}
          reportData={data}
        />

        <div className="flex flex-wrap gap-2 mb-6 border-b border-gray-200 pb-4">
          {tabs.map(tab => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={`px-4 py-2 rounded-lg font-medium text-sm transition-colors ${
                activeTab === tab.key
                  ? 'bg-blue-500 text-white'
                  : 'bg-white text-gray-600 hover:bg-gray-100'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        <div className="space-y-6">
          {activeTab === 'overview' && (
            <>
              <KeyStatsCards stats={data.stats || { total_videos: 0, total_comments: 0, comments_by_brand: {}}} brandCount={data.brands.length} />
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <BrandScoreChart data={data} />
                <BrandRadarChart data={data} />
              </div>
            </>
          )}

          {activeTab === 'charts' && <RadarBrandSelector data={data} />}

          {activeTab === 'summary' && <EnhancedSummary recommendation={data.recommendation} />}
        </div>
      </div>
    </div>
  )
}

export default Report
