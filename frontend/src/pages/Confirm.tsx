import { useEffect, useState } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { confirmNavigationTarget, parseOutcomeFromResponse } from '../lib/confirmFlow'
import { addListItem, buildConfirmPayload, removeListItem } from '../lib/confirmPlan'
import { apiClient } from '../api/client'
import axios from 'axios'

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
  const [error, setError] = useState<string | null>(null)
  const [newBrand, setNewBrand] = useState('')
  const [newKeyword, setNewKeyword] = useState('')
  const [newDimName, setNewDimName] = useState('')

  useEffect(() => {
    if (!requirement) {
        setLoading(false);
        return;
    }

    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        const result = await apiClient.post<ParseResponse>('/parse', { requirement })
        const outcome = parseOutcomeFromResponse(true, result as unknown as Record<string, unknown>)
        if (outcome.kind === 'error') {
          setError(outcome.message)
          setData(null)
          return
        }
        setData(outcome.data as unknown as ParseResponse)
      } catch (err) {
        const apiError = axios.isAxiosError(err) ? err.response?.data : null
        const message = apiError && typeof apiError === 'object' && 'error' in apiError && typeof apiError.error === 'string'
          ? apiError.error
          : '解析需求失败，请检查设置中的 AI 配置'
        setError(message)
        setData(null)
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
      const result = await apiClient.post<{ task_id?: string; error?: string }>('/confirm', buildConfirmPayload({
          requirement: requirement,
          budget: data.budget,
          scenario: data.scenario,
          special_needs: data.special_needs,
          brands: data.brands,
          dimensions: data.dimensions,
          keywords: data.keywords,
        }, {
          video_date_range_months: videoDateRangeMonths,
          min_video_duration: minVideoDuration,
          max_comments: maxComments,
          min_video_comments: minVideoComments,
          min_comments_per_video: minCommentsPerVideo,
          max_comments_per_video_v2: maxCommentsPerVideo
        }))
      const nav = confirmNavigationTarget(true, result, data.product_type)
      if (nav.kind === 'error') {
        setError(nav.message)
        setSubmitting(false)
        return
      }
      navigate(nav.href)
    } catch (err) {
      const apiError = axios.isAxiosError(err) ? err.response?.data : null
      const message = apiError && typeof apiError === 'object' && 'error' in apiError && typeof apiError.error === 'string'
        ? apiError.error
        : '创建任务失败，请稍后重试'
      setError(message)
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
        <p className="text-gray-500 mt-2">{error || '请返回首页重新提交需求'}</p>
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

            {(data.special_needs || []).length > 0 && (
              <p className="text-sm text-slate-600">特殊需求：{(data.special_needs || []).join('、')}</p>
            )}

            <div>
                <h4 className="text-sm font-bold text-gray-600 mb-4 flex items-center gap-2">
                    <span>🏷️</span> 将分析这些品牌
                </h4>
                <div className="flex flex-wrap gap-3 mb-3">
                {(data.brands || []).map(brand => (
                    <button
                      type="button"
                      key={brand}
                      onClick={() => setData({ ...data, brands: removeListItem(data.brands, brand) })}
                      className="px-4 py-2 bg-white/50 rounded-xl text-sm font-medium text-slate-700 border border-slate-200/60"
                    >
                    {brand} ×
                    </button>
                ))}
                </div>
                <div className="flex gap-2">
                  <input value={newBrand} onChange={(e) => setNewBrand(e.target.value)} placeholder="添加品牌" className="flex-1 px-3 py-2 border border-gray-200 rounded-lg text-sm" />
                  <button type="button" className="px-3 py-2 bg-gray-800 text-white rounded-lg text-sm" onClick={() => { setData({ ...data, brands: addListItem(data.brands, newBrand) }); setNewBrand('') }}>添加</button>
                </div>
            </div>

            <div>
                <h4 className="text-sm font-bold text-gray-600 mb-4 flex items-center gap-2">
                    <span>📊</span> 评价维度
                </h4>
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 mb-3">
                {(data.dimensions || []).map(dim => (
                    <div key={dim.name} className="bg-white/40 rounded-xl p-4 border border-white/40">
                    <div className="flex justify-between gap-2">
                      <h5 className="font-bold text-slate-800 mb-1">{dim.name}</h5>
                      <button type="button" className="text-xs text-red-500" onClick={() => setData({ ...data, dimensions: data.dimensions.filter((d) => d.name !== dim.name) })}>删除</button>
                    </div>
                    <p className="text-xs text-slate-500 leading-relaxed">{dim.description}</p>
                    </div>
                ))}
                </div>
                <div className="flex gap-2">
                  <input value={newDimName} onChange={(e) => setNewDimName(e.target.value)} placeholder="添加维度" className="flex-1 px-3 py-2 border border-gray-200 rounded-lg text-sm" />
                  <button type="button" className="px-3 py-2 bg-gray-800 text-white rounded-lg text-sm" onClick={() => {
                    const name = newDimName.trim()
                    if (!name || data.dimensions.some((d) => d.name === name)) return
                    setData({ ...data, dimensions: [...data.dimensions, { name, description: name }] })
                    setNewDimName('')
                  }}>添加</button>
                </div>
            </div>

            <div>
                <h4 className="text-sm font-bold text-gray-600 mb-3 flex items-center gap-2">
                    <span>🔍</span> 搜索关键词
                </h4>
                <div className="flex flex-wrap gap-2 mb-3">
                  {(data.keywords || []).map((word) => (
                    <button type="button" key={word} className="px-3 py-1 bg-gray-50 rounded-lg text-sm" onClick={() => setData({ ...data, keywords: removeListItem(data.keywords, word) })}>
                      {word} ×
                    </button>
                  ))}
                </div>
                <div className="flex gap-2">
                  <input value={newKeyword} onChange={(e) => setNewKeyword(e.target.value)} placeholder="添加关键词" className="flex-1 px-3 py-2 border border-gray-200 rounded-lg text-sm" />
                  <button type="button" className="px-3 py-2 bg-gray-800 text-white rounded-lg text-sm" onClick={() => { setData({ ...data, keywords: addListItem(data.keywords, newKeyword) }); setNewKeyword('') }}>添加</button>
                </div>
            </div>
        </div>

        {error && (
          <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}

        {/* Confirm Button */}
        <button
          onClick={handleConfirm}
          disabled={submitting || data.brands.length === 0 || data.dimensions.length === 0 || data.keywords.length === 0}
          className="w-full py-4 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-bold rounded-xl shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all cursor-pointer flex items-center justify-center gap-2"
        >
          {submitting ? '⏳ 正在创建任务...' : '✓ 确认开始分析'}
        </button>
      </div>
    </div>
  )
}

export default Confirm
