import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import ReactECharts from 'echarts-for-react'
import ReactMarkdown from 'react-markdown'
import { useToast } from '../hooks/useToast'

interface BrandRanking {
  brand: string
  overall_score: number
  rank: number
  scores: Record<string, number>
}

interface ModelRanking {
  model: string
  brand: string
  overall_score: number
  rank: number
  scores: Record<string, number>
  comment_count: number
}

interface Dimension {
  name: string
  description: string
}

interface ReportStats {
  total_videos: number
  total_comments: number
  comments_by_brand: Record<string, number>
}

interface TypicalComment {
  content: string
  score: number
}

interface BrandAnalysis {
  strengths: string[]
  weaknesses: string[]
}

interface ReportData {
  category: string
  brands: string[]
  dimensions: Dimension[]
  scores: Record<string, Record<string, number>>
  rankings: BrandRanking[]
  model_rankings?: ModelRanking[]
  recommendation: string
  stats?: ReportStats
  top_comments?: Record<string, TypicalComment[]>
  bad_comments?: Record<string, TypicalComment[]>
  brand_analysis?: Record<string, BrandAnalysis>
}

interface ApiResponse {
  id: number
  history_id: number
  category: string
  data: ReportData
  created_at: string
}

const Report = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [report, setReport] = useState<ApiResponse | null>(null)
  const [exporting, setExporting] = useState(false)
  const [specifiedBrands, setSpecifiedBrands] = useState<string[]>([])
  const { showToast } = useToast()

  const handleExportPDF = async () => {
    if (!id) return
    setExporting(true)
    try {
      const response = await fetch(`http://localhost:8080/api/report/${id}/pdf`)
      if (!response.ok) {
        throw new Error('导出失败')
      }
      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `report_${id}.pdf`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    } catch (err) {
      showToast(err instanceof Error ? err.message : '导出PDF失败', 'error')
    } finally {
      setExporting(false)
    }
  }

  useEffect(() => {
    if (!id) return

    const fetchReport = async () => {
      try {
        setLoading(true)
        const response = await fetch(`http://localhost:8080/api/report/${id}`)
        if (!response.ok) {
          throw new Error('报告不存在')
        }
        const data = await response.json()
        setReport(data)

        if (data.history_id) {
          try {
            const historyRes = await fetch(`http://localhost:8080/api/history/${data.history_id}`)
            if (historyRes.ok) {
              const historyData = await historyRes.json()
              setSpecifiedBrands(historyData.brands || [])
            }
          } catch (e) {
            console.error("Failed to fetch history brands", e)
          }
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : '加载报告失败')
      } finally {
        setLoading(false)
      }
    }

    fetchReport()
  }, [id])

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh]">
        <div className="w-16 h-16 border-4 border-blue-500/30 border-t-blue-500 rounded-full animate-spin mb-6"></div>
        <h2 className="text-2xl font-semibold text-gray-700">正在加载报告...</h2>
      </div>
    )
  }

  if (error || !report) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh]">
        <h2 className="text-2xl font-semibold text-red-600 mb-4">加载失败</h2>
        <p className="text-gray-500 mb-6">{error || '报告数据不存在'}</p>
        <button
          onClick={() => navigate('/')}
          className="px-6 py-2 bg-gray-800 text-white rounded-lg hover:bg-gray-700 transition-colors cursor-pointer"
        >
          返回首页
        </button>
      </div>
    )
  }

  const reportData = report.data
  const topBrand = reportData.rankings[0]

  const strengths = reportData.dimensions
    .filter(dim => topBrand?.scores[dim.name] >= 8)
    .map(dim => dim.name)

  const weaknesses = reportData.dimensions
    .filter(dim => topBrand?.scores[dim.name] && topBrand.scores[dim.name] < 6)
    .map(dim => dim.name)

  const colors = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981']

  // ECharts 雷达图配置
  const radarOption = {
    tooltip: {
      trigger: 'item'
    },
    legend: {
      data: reportData.brands.slice(0, 3),
      bottom: 10,
      textStyle: {
        fontSize: 12
      }
    },
    radar: {
      indicator: reportData.dimensions.map(dim => ({
        name: dim.name,
        max: 100
      })),
      splitNumber: 4,
      name: {
        textStyle: {
          color: '#6b7280',
          fontSize: 11
        }
      },
      splitLine: {
        lineStyle: {
          color: '#e5e7eb'
        }
      },
      splitArea: {
        show: true,
        areaStyle: {
          color: ['rgba(255, 255, 255, 0)', 'rgba(249, 250, 251, 0.5)']
        }
      },
      axisLine: {
        lineStyle: {
          color: '#e5e7eb'
        }
      }
    },
    series: [{
      type: 'radar',
      data: reportData.brands.slice(0, 3).map((brand, index) => ({
        value: reportData.dimensions.map(dim => 
          reportData.scores[brand]?.[dim.name] ? reportData.scores[brand][dim.name] * 10 : 0
        ),
        name: brand,
        lineStyle: {
          color: colors[index],
          width: 2
        },
        areaStyle: {
          color: colors[index],
          opacity: 0.2
        },
        itemStyle: {
          color: colors[index]
        }
      }))
    }]
  }

  // ECharts 柱状图配置
  const barOption = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: 80,
      right: 30,
      top: 20,
      bottom: 20
    },
    xAxis: {
      type: 'value',
      max: 10,
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: '#6b7280'
      },
      splitLine: {
        lineStyle: {
          color: '#f0f0f0'
        }
      }
    },
    yAxis: {
      type: 'category',
      data: [...reportData.rankings].reverse().map(r => r.brand),
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: '#6b7280'
      }
    },
    series: [{
      type: 'bar',
      data: [...reportData.rankings].reverse().map((r, index) => ({
        value: Math.round(r.overall_score * 10) / 10,
        itemStyle: {
          color: colors[index % colors.length],
          borderRadius: [0, 8, 8, 0]
        }
      })),
      barWidth: 30,
      label: {
        show: true,
        position: 'right',
        formatter: '{c}',
        color: '#374151',
        fontWeight: 'bold'
      }
    }]
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8 space-y-6">
      <div className="flex justify-between items-end mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-800">产品分析报告</h1>
          <p className="text-gray-500">{reportData.category} | 报告ID: {id}</p>
        </div>
        <div className="flex gap-2">
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
          <button 
            onClick={() => navigate('/history')}
            className="px-4 py-2 bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 transition-colors font-medium text-sm cursor-pointer"
          >
            查看历史
          </button>
        </div>
      </div>

      {reportData.stats && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="glass-card p-6 text-center">
            <div className="text-3xl font-bold text-blue-600">{reportData.stats.total_videos}</div>
            <div className="text-gray-500 text-sm mt-1">分析视频数</div>
          </div>
          <div className="glass-card p-6 text-center">
            <div className="text-3xl font-bold text-purple-600">{reportData.stats.total_comments}</div>
            <div className="text-gray-500 text-sm mt-1">评论总数</div>
          </div>
          <div className="glass-card p-6 text-center">
            <div className="text-3xl font-bold text-green-600">{reportData.brands.length}</div>
            <div className="text-gray-500 text-sm mt-1">对比品牌数</div>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="glass-card flex flex-col items-center justify-center p-8 bg-gradient-to-br from-blue-500 to-indigo-600 text-white border-none">
          <div className="text-6xl font-bold mb-2">
            {topBrand ? topBrand.overall_score.toFixed(1) : '-'}
          </div>
          <div className="text-xl font-medium opacity-90">
            {topBrand ? `${topBrand.brand} 领先` : '暂无数据'}
          </div>
          <div className="mt-4 px-3 py-1 bg-white/20 rounded-full text-sm backdrop-blur-sm">
            共分析 {reportData.brands.length} 个品牌
          </div>
        </div>

        <div className="glass-card">
          <h3 className="text-lg font-semibold text-green-600 mb-4 flex items-center">
            <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5" />
            </svg>
            {topBrand?.brand} 优势
          </h3>
          <ul className="space-y-3">
            {strengths.length > 0 ? strengths.map((item, i) => (
              <li key={i} className="flex items-start text-gray-700 text-sm">
                <span className="mr-2 text-green-500 mt-1">●</span>
                {item}表现优秀
              </li>
            )) : (
              <li className="text-gray-500 text-sm">暂无突出优势</li>
            )}
          </ul>
        </div>

        <div className="glass-card">
          <h3 className="text-lg font-semibold text-amber-500 mb-4 flex items-center">
            <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            待改进
          </h3>
          <ul className="space-y-3">
            {weaknesses.length > 0 ? weaknesses.map((item, i) => (
              <li key={i} className="flex items-start text-gray-700 text-sm">
                <span className="mr-2 text-amber-400 mt-1">●</span>
                {item}有提升空间
              </li>
            )) : (
              <li className="text-gray-500 text-sm">各维度表现均衡</li>
            )}
          </ul>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="glass-card lg:col-span-1 min-h-[400px]">
          <h3 className="text-lg font-semibold text-gray-800 mb-4 text-center">多维度对比</h3>
          <ReactECharts option={radarOption} style={{ height: '320px' }} />
        </div>

        <div className="glass-card lg:col-span-2 min-h-[400px]">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">品牌综合得分排名</h3>
          <ReactECharts option={barOption} style={{ height: '320px' }} />
        </div>
      </div>

      {reportData.model_rankings && reportData.model_rankings.length > 0 && (
        <div className="glass-card">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">🏆 型号排名 (Model Rankings)</h3>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-center py-3 px-4 font-semibold text-gray-600">排名</th>
                  <th className="text-left py-3 px-4 font-semibold text-gray-600">型号</th>
                  <th className="text-left py-3 px-4 font-semibold text-gray-600">品牌</th>
                  <th className="text-center py-3 px-4 font-semibold text-gray-600">综合得分</th>
                  <th className="text-center py-3 px-4 font-semibold text-gray-600">样本数</th>
                </tr>
              </thead>
              <tbody>
                {reportData.model_rankings.map((ranking, index) => (
                  <tr key={`${ranking.brand}-${ranking.model}`} className="border-b border-gray-100 hover:bg-gray-50">
                    <td className="text-center py-3 px-4">
                      <span className={`inline-flex items-center gap-1 ${index === 0 ? 'text-blue-600 font-bold' : ''}`}>
                        {index === 0 && <span className="text-yellow-500">🥇</span>}
                        {index === 1 && <span className="text-gray-400">🥈</span>}
                        {index === 2 && <span className="text-orange-400">🥉</span>}
                        {ranking.rank}
                      </span>
                    </td>
                    <td className="py-3 px-4 font-medium text-gray-800">{ranking.model}</td>
                    <td className="py-3 px-4 text-gray-600">{ranking.brand}</td>
                    <td className="text-center py-3 px-4">
                      <span className="px-2 py-1 rounded text-xs font-bold bg-indigo-100 text-indigo-700">
                        {ranking.overall_score.toFixed(1)}
                      </span>
                    </td>
                    <td className="text-center py-3 px-4 text-gray-500">{ranking.comment_count}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      <div className="glass-card">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">各维度详细得分</h3>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200">
                <th className="text-left py-3 px-4 font-semibold text-gray-600">品牌</th>
                {reportData.dimensions.map(dim => (
                  <th key={dim.name} className="text-center py-3 px-4 font-semibold text-gray-600">
                    {dim.name}
                  </th>
                ))}
                <th className="text-center py-3 px-4 font-semibold text-gray-600">综合</th>
              </tr>
            </thead>
            <tbody>
              {reportData.rankings.map((ranking, index) => (
                <tr key={ranking.brand} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-3 px-4 font-medium">
                    <span className={`inline-flex items-center gap-2 ${index === 0 ? 'text-blue-600' : ''}`}>
                      {index === 0 && <span className="text-yellow-500">🏆</span>}
                      {ranking.brand}
                      {specifiedBrands.length > 0 && !specifiedBrands.includes(ranking.brand) && (
                        <span 
                          className="text-blue-500 cursor-help" 
                          title="此品牌由AI自动发现"
                        >
                          🔍
                        </span>
                      )}
                    </span>
                  </td>
                  {reportData.dimensions.map(dim => {
                    const score = ranking.scores?.[dim.name]
                    return (
                      <td key={dim.name} className="text-center py-3 px-4">
                        <span className={`
                          px-2 py-1 rounded text-xs font-medium
                          ${score && score >= 8 ? 'bg-green-100 text-green-700' :
                            score && score >= 6 ? 'bg-blue-100 text-blue-700' :
                            score ? 'bg-amber-100 text-amber-700' : 'bg-gray-100 text-gray-500'}
                        `}>
                          {score?.toFixed(1) || '-'}
                        </span>
                      </td>
                    )
                  })}
                  <td className="text-center py-3 px-4">
                    <span className="px-2 py-1 rounded text-xs font-bold bg-indigo-100 text-indigo-700">
                      {ranking.overall_score.toFixed(1)}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {(reportData.top_comments || reportData.bad_comments) && (
        <div className="glass-card">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">典型用户评论</h3>
          <div className="space-y-6">
            {reportData.brands.map(brand => {
              const topList = reportData.top_comments?.[brand] || []
              const badList = reportData.bad_comments?.[brand] || []
              if (topList.length === 0 && badList.length === 0) return null
              return (
                <div key={brand} className="border-b border-gray-100 pb-4 last:border-0">
                  <h4 className="font-medium text-gray-700 mb-3">{brand}</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {topList.length > 0 && (
                      <div>
                        <div className="text-sm font-medium text-green-600 mb-2">好评</div>
                        <div className="space-y-2">
                          {topList.slice(0, 2).map((comment, i) => (
                            <div key={i} className="p-3 bg-green-50 border border-green-200 rounded-lg text-sm text-gray-700">
                              <p className="line-clamp-3">{comment.content}</p>
                              <span className="text-xs text-green-600 mt-1 block">评分: {comment.score.toFixed(1)}</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                    {badList.length > 0 && (
                      <div>
                        <div className="text-sm font-medium text-red-600 mb-2">差评</div>
                        <div className="space-y-2">
                          {badList.slice(0, 2).map((comment, i) => (
                            <div key={i} className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-gray-700">
                              <p className="line-clamp-3">{comment.content}</p>
                              <span className="text-xs text-red-600 mt-1 block">评分: {comment.score.toFixed(1)}</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      )}

      <div className="glass-card bg-gradient-to-r from-gray-50 to-gray-100 border border-gray-200">
        <h3 className="text-lg font-bold text-gray-800 mb-4">💡 购买建议</h3>
        <div className="prose prose-sm max-w-none text-gray-700">
          <ReactMarkdown>{reportData.recommendation}</ReactMarkdown>
        </div>
      </div>
    </div>
  )
}

export default Report
