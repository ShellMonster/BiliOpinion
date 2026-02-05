import { useEffect, useState } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'

interface ParseResponse {
  understanding: string
  product_type: string
  budget?: string
  scenario?: string
  special_needs?: string[]
  brands: string[]
  dimensions: Array<{
    name: string
    description: string
  }>
  keywords: string[]
}

const Confirm = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const requirement = searchParams.get('requirement')
  
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<ParseResponse | null>(null)
  const [videoDateRangeMonths, setVideoDateRangeMonths] = useState(0)
  const [minVideoDuration, setMinVideoDuration] = useState(30)
  const [maxComments, setMaxComments] = useState(500)
  const [minVideoComments, setMinVideoComments] = useState(0)
  const [minCommentsPerVideo, setMinCommentsPerVideo] = useState(20)
  const [maxCommentsPerVideo, setMaxCommentsPerVideo] = useState(200)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (!requirement) {
        setLoading(false);
        return;
    }

    const fetchData = async () => {
      try {
        setLoading(true)
        const response = await fetch('http://localhost:8080/api/parse', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ requirement })
        })
        const result = await response.json()
        setData(result)
      } catch (error) {
        console.error('Failed to parse requirement:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [requirement])

  const handleConfirm = async () => {
    if (!data || !requirement || submitting) return
    
    setSubmitting(true)
    try {
      const response = await fetch('http://localhost:8080/api/confirm', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          requirement: requirement,
          brands: data.brands,
          dimensions: data.dimensions,
          keywords: data.keywords,
          video_date_range_months: videoDateRangeMonths,
          min_video_duration: minVideoDuration,
          max_comments: maxComments,
          min_video_comments: minVideoComments,
          min_comments_per_video: minCommentsPerVideo,
          max_comments_per_video_v2: maxCommentsPerVideo
        })
      })
      const result = await response.json()
      navigate(`/progress/${result.task_id}?title=${encodeURIComponent(data.product_type)}`)
    } catch (error) {
      console.error('Failed to confirm:', error)
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh]">
        <div className="w-16 h-16 border-4 border-blue-500/30 border-t-blue-500 rounded-full animate-spin mb-6"></div>
        <h2 className="text-2xl font-semibold text-gray-700">正在解析您的需求...</h2>
        <p className="text-gray-500 mt-2">AI 正在分析商品类型、评价维度与品牌信息</p>
      </div>
    )
  }

  if (!data || !requirement) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh]">
        <h2 className="text-2xl font-semibold text-red-600">无法获取分析数据</h2>
        <p className="text-gray-500 mt-2">请返回首页重新提交需求</p>
        <button 
          onClick={() => navigate('/')}
          className="mt-6 px-6 py-2 bg-gray-800 text-white rounded-lg hover:bg-gray-700 transition-colors"
        >
          返回首页
        </button>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="mb-8 text-center">
        <h1 className="text-3xl font-bold text-gray-800 mb-2">确认分析方案</h1>
        <p className="text-gray-500">AI 已为您生成个性化分析计划，请确认细节</p>
      </div>

      <div className="space-y-6">
        {/* Understanding Card */}
        <div className="bg-blue-50/80 backdrop-blur-sm rounded-2xl p-6 border border-blue-100 shadow-sm">
          <h3 className="text-lg font-bold text-blue-900 mb-2">💡 我理解您的需求</h3>
          <p className="text-slate-700 leading-relaxed">{data.understanding}</p>
        </div>

        {/* Analysis Plan Card */}
        <div className="glass-card p-8 space-y-8">
            
            {/* Info Row */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 pb-6 border-b border-gray-100">
                <div>
                    <span className="text-xs font-bold text-gray-400 uppercase tracking-wider">商品类型</span>
                    <p className="text-lg font-medium text-gray-800 mt-1">{data.product_type}</p>
                </div>
                {data.budget && (
                <div>
                    <span className="text-xs font-bold text-gray-400 uppercase tracking-wider">预算范围</span>
                    <p className="text-lg font-medium text-gray-800 mt-1">{data.budget}</p>
                </div>
                )}
                {data.scenario && (
                <div>
                    <span className="text-xs font-bold text-gray-400 uppercase tracking-wider">使用场景</span>
                    <p className="text-lg font-medium text-gray-800 mt-1">{data.scenario}</p>
                </div>
                )}
            </div>

            <div className="bg-blue-50/50 rounded-xl p-4 border border-blue-100">
              <label className="block text-sm font-bold text-gray-700 mb-2 flex items-center gap-2">
                <span>📅</span> 分析时间范围
              </label>
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                {/* 视频发布时间选项 */}
                <div>
                  <select
                    value={videoDateRangeMonths}
                    onChange={(e) => setVideoDateRangeMonths(Number(e.target.value))}
                    className="w-full px-4 py-2 bg-white border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all shadow-sm text-gray-700 font-medium"
                  >
                    <option value={6}>最近 6 个月</option>
                    <option value={12}>最近 1 年</option>
                    <option value={24}>最近 2 年</option>
                    <option value={0}>不限时间 (推荐)</option>
                  </select>
                  <p className="text-xs text-gray-500 mt-1">视频发布时间</p>
                </div>
                {/* 视频时长过滤选项 */}
                <div>
                  <select
                    value={minVideoDuration}
                    onChange={(e) => setMinVideoDuration(Number(e.target.value))}
                    className="w-full px-4 py-2 bg-white border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all shadow-sm text-gray-700 font-medium"
                  >
                    <option value={0}>不限制</option>
                    <option value={30}>至少 30 秒 (推荐)</option>
                    <option value={60}>至少 1 分钟</option>
                    <option value={120}>至少 2 分钟</option>
                    <option value={180}>至少 3 分钟</option>
                    <option value={300}>至少 5 分钟</option>
                  </select>
                  <p className="text-xs text-gray-500 mt-1">过滤短视频</p>
                </div>
                {/* 评论数量限制选项 */}
                <div>
                  <select
                    value={maxComments}
                    onChange={(e) => setMaxComments(Number(e.target.value))}
                    className="w-full px-4 py-2 bg-white border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all shadow-sm text-gray-700 font-medium"
                  >
                    <option value={200}>限制 200 条</option>
                    <option value={500}>限制 500 条 (推荐)</option>
                    <option value={1000}>限制 1000 条</option>
                    <option value={2000}>限制 2000 条</option>
                    <option value={5000}>限制 5000 条</option>
                    <option value={10000}>限制 10000 条</option>
                  </select>
                  <p className="text-xs text-gray-500 mt-1">分析评论数量</p>
                </div>
              </div>
            </div>

            <div className="bg-purple-50/50 rounded-xl p-4 border border-purple-100">
              <label className="block text-sm font-bold text-gray-700 mb-2 flex items-center gap-2">
                <span>🎯</span> 评论抓取策略
              </label>
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                {/* 最小视频评论数 */}
                <div>
                  <select
                    value={minVideoComments}
                    onChange={(e) => setMinVideoComments(Number(e.target.value))}
                    className="w-full px-4 py-2 bg-white border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all shadow-sm text-gray-700 font-medium"
                  >
                    <option value={0}>不限制 (推荐)</option>
                    <option value={50}>至少 50 条</option>
                    <option value={100}>至少 100 条</option>
                    <option value={200}>至少 200 条</option>
                    <option value={500}>至少 500 条</option>
                  </select>
                  <p className="text-xs text-gray-500 mt-1">最小视频评论数</p>
                </div>
                {/* 每视频最少抓取 */}
                <div>
                  <select
                    value={minCommentsPerVideo}
                    onChange={(e) => setMinCommentsPerVideo(Number(e.target.value))}
                    className="w-full px-4 py-2 bg-white border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all shadow-sm text-gray-700 font-medium"
                  >
                    <option value={10}>至少 10 条</option>
                    <option value={20}>至少 20 条 (推荐)</option>
                    <option value={50}>至少 50 条</option>
                    <option value={100}>至少 100 条</option>
                  </select>
                  <p className="text-xs text-gray-500 mt-1">每视频最少抓取</p>
                </div>
                {/* 每视频最多抓取 */}
                <div>
                  <select
                    value={maxCommentsPerVideo}
                    onChange={(e) => setMaxCommentsPerVideo(Number(e.target.value))}
                    className="w-full px-4 py-2 bg-white border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all shadow-sm text-gray-700 font-medium"
                  >
                    <option value={100}>最多 100 条</option>
                    <option value={200}>最多 200 条 (推荐)</option>
                    <option value={500}>最多 500 条</option>
                    <option value={1000}>最多 1000 条</option>
                  </select>
                  <p className="text-xs text-gray-500 mt-1">每视频最多抓取</p>
                </div>
              </div>
            </div>

            {/* Brand Tags */}
            <div>
                <h4 className="text-sm font-bold text-gray-600 mb-4 flex items-center gap-2">
                    <span>🏷️</span> 将分析这些品牌
                </h4>
                <div className="flex flex-wrap gap-3">
                {(data.brands || []).map(brand => (
                    <span key={brand} className="px-4 py-2 bg-white/50 backdrop-blur-sm rounded-xl text-sm font-medium text-slate-700 border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow cursor-default">
                    {brand}
                    </span>
                ))}
                </div>
            </div>

            {/* Dimension Cards */}
            <div>
                <h4 className="text-sm font-bold text-gray-600 mb-4 flex items-center gap-2">
                    <span>📊</span> 评价维度
                </h4>
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
                {(data.dimensions || []).map(dim => (
                    <div key={dim.name} className="bg-white/40 backdrop-blur-sm rounded-xl p-4 border border-white/40 hover:bg-white/60 transition-colors">
                    <h5 className="font-bold text-slate-800 mb-1">{dim.name}</h5>
                    <p className="text-xs text-slate-500 leading-relaxed">{dim.description}</p>
                    </div>
                ))}
                </div>
            </div>

            {/* Keywords */}
            <div>
                <h4 className="text-sm font-bold text-gray-600 mb-3 flex items-center gap-2">
                    <span>🔍</span> 搜索关键词
                </h4>
                <div className="bg-gray-50/50 rounded-lg p-3 text-sm text-slate-600 font-mono border border-gray-100">
                    {(data.keywords || []).join(' | ')}
                </div>
            </div>
        </div>

        {/* Confirm Button */}
        <button
          onClick={handleConfirm}
          disabled={submitting}
          className="w-full py-4 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-bold rounded-xl shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all cursor-pointer flex items-center justify-center gap-2"
        >
          {submitting ? '⏳ 正在创建任务...' : '✓ 确认开始分析'}
        </button>
      </div>
    </div>
  )
}

export default Confirm
