import { useEffect, useState } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { parseVideo, analyzeVideo, getDimensions, type VideoParseResponse, type Dimension } from '../api/video'



const defaultDimensions: Dimension[] = [
  { name: '性能表现', description: '产品的核心性能、功能实现程度' },
  { name: '质量做工', description: '产品质量、材料、做工精细度' },
  { name: '性价比', description: '价格与性能的匹配程度，是否物有所值' },
  { name: '使用体验', description: '日常使用感受、操作便捷性' },
  { name: '外观设计', description: '产品外观、颜值、设计感' },
  { name: '售后服务', description: '售后服务质量、客服响应、保修政策' },
]

const commentOptions = [
  { value: 100, label: '100条' },
  { value: 500, label: '500条' },
  { value: 1000, label: '1000条（推荐）' },
  { value: 2000, label: '2000条' },
  { value: 5000, label: '5000条' },
  { value: 10000, label: '10000条' },
  { value: 0, label: '全部评论' },
]

function formatNumber(num: number): string {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toLocaleString()
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

const VideoConfirm = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const videoUrl = searchParams.get('video_url')

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [videoInfo, setVideoInfo] = useState<VideoParseResponse | null>(null)
  const [dimensions, setDimensions] = useState<Dimension[]>([])
  const [dimensionsLoading, setDimensionsLoading] = useState(true)
  const [maxComments, setMaxComments] = useState(1000)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (!videoUrl) {
      setLoading(false)
      setError('缺少视频链接参数')
      return
    }

    const fetchVideoInfo = async () => {
      try {
        setLoading(true)
        setError('')
        const result = await parseVideo(videoUrl)
        setVideoInfo(result)
        setLoading(false)

        setDimensionsLoading(true)
        try {
          const dimensionsResult = await getDimensions(result.bvid)
          setDimensions(dimensionsResult.dimensions)
        } catch (err) {
          console.error('Failed to fetch dimensions:', err)
          setDimensions(defaultDimensions)
        } finally {
          setDimensionsLoading(false)
        }
      } catch (err) {
        console.error('Failed to parse video:', err)
        setError('获取视频信息失败，请检查链接是否有效')
        setLoading(false)
      }
    }

    fetchVideoInfo()
  }, [videoUrl])

  const handleAnalyze = async () => {
    if (!videoUrl || !videoInfo || submitting) return

    setSubmitting(true)
    try {
      const result = await analyzeVideo(videoUrl, maxComments)
      navigate(`/progress/${result.task_id}?title=${encodeURIComponent(videoInfo.title)}`)
    } catch (err) {
      console.error('Failed to start analysis:', err)
      setError('启动分析失败，请稍后重试')
      setSubmitting(false)
    }
  }

  const handleCancel = () => {
    navigate('/')
  }

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh]">
        <div className="w-16 h-16 border-4 border-blue-500/30 border-t-blue-500 rounded-full animate-spin mb-6"></div>
        <h2 className="text-2xl font-semibold text-gray-700">正在获取视频信息...</h2>
        <p className="text-gray-500 mt-2">请稍候，正在解析视频链接</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh]">
        <div className="w-16 h-16 mb-6 flex items-center justify-center">
          <svg className="w-12 h-12 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h2 className="text-2xl font-semibold text-red-600">{error}</h2>
        <p className="text-gray-500 mt-2">请检查视频链接是否正确</p>
        <button
          onClick={handleCancel}
          className="mt-6 px-6 py-2 bg-gray-800 text-white rounded-lg hover:bg-gray-700 transition-colors cursor-pointer"
        >
          返回首页
        </button>
      </div>
    )
  }

  if (!videoInfo) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh]">
        <h2 className="text-2xl font-semibold text-red-600">无法获取视频信息</h2>
        <p className="text-gray-500 mt-2">请返回首页重新提交</p>
        <button
          onClick={handleCancel}
          className="mt-6 px-6 py-2 bg-gray-800 text-white rounded-lg hover:bg-gray-700 transition-colors cursor-pointer"
        >
          返回首页
        </button>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="mb-8 text-center">
        <h1 className="text-3xl font-bold text-gray-800 mb-2">确认视频分析</h1>
        <p className="text-gray-500">请确认视频信息和分析设置</p>
      </div>

      <div className="space-y-6">
        {/* 视频信息卡片 */}
        <div className="glass-card p-6">
          {/* 标题 - 居左 */}
          <h2 className="text-xl font-bold text-gray-800 leading-tight mb-4">
            {videoInfo.title}
          </h2>
          
          {/* 封面图 + 信息 - 上下居中对齐 */}
          <div className="flex flex-col md:flex-row gap-6 items-center">
            {/* 左侧：封面图（放大） */}
            <div className="w-full md:w-64 flex-shrink-0">
              <div className="aspect-video rounded-xl overflow-hidden bg-gray-100 shadow-md">
                {videoInfo.cover ? (
                  <img
                    referrerPolicy="no-referrer"
                    src={videoInfo.cover}
                    alt={videoInfo.title}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center bg-gray-200">
                    <svg className="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                )}
              </div>
            </div>

            {/* 右侧：视频信息（靠右） */}
            <div className="flex-1 w-full flex">
              <div className="w-full ">
                <div className="grid grid-cols-2 gap-3">
                  <div className="bg-blue-50/50 rounded-lg p-3 border border-blue-100">
                    <span className="text-xs font-medium text-blue-600 uppercase tracking-wider">UP主</span>
                    <p className="text-sm font-semibold text-gray-800 mt-1 truncate">{videoInfo.author}</p>
                  </div>
                  <div className="bg-purple-50/50 rounded-lg p-3 border border-purple-100">
                    <span className="text-xs font-medium text-purple-600 uppercase tracking-wider">发布时间</span>
                    <p className="text-sm font-semibold text-gray-800 mt-1">{formatDate(videoInfo.pub_date)}</p>
                  </div>
                  <div className="bg-green-50/50 rounded-lg p-3 border border-green-100">
                    <span className="text-xs font-medium text-green-600 uppercase tracking-wider">播放量</span>
                    <p className="text-sm font-semibold text-gray-800 mt-1">{formatNumber(videoInfo.play_count)}</p>
                  </div>
                  <div className="bg-orange-50/50 rounded-lg p-3 border border-orange-100">
                    <span className="text-xs font-medium text-orange-600 uppercase tracking-wider">评论数</span>
                    <p className="text-sm font-semibold text-gray-800 mt-1">{formatNumber(videoInfo.comment_count)}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* 评论数量选择 */}
        <div className="glass-card p-6">
          <label className="block text-sm font-bold text-gray-700 mb-3 flex items-center gap-2">
            <svg className="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
            </svg>
            分析评论数量
          </label>
          <select
            value={maxComments}
            onChange={(e) => setMaxComments(Number(e.target.value))}
            className="w-full  px-4 py-3 bg-white border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all shadow-sm text-gray-700 font-medium cursor-pointer"
          >
            {commentOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          <p className="text-sm text-gray-500 mt-2">
            选择要分析的评论数量，数量越多分析越准确但耗时越长
          </p>
        </div>

        {/* 评价维度 */}
        <div className="glass-card p-6">
          <h4 className="text-sm font-bold text-gray-700 mb-4 flex items-center gap-2">
            <svg className="w-5 h-5 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            评价维度
            <span className="text-xs font-normal text-gray-400 ml-2">
              ({dimensionsLoading ? '...' : (dimensions.length || defaultDimensions.length)} 个维度)
            </span>
          </h4>
          {dimensionsLoading ? (
            <div className="text-center py-8">
              <div className="w-8 h-8 border-2 border-blue-500/30 border-t-blue-500 rounded-full animate-spin mx-auto mb-3"></div>
              <p className="text-gray-500">正在分析维度...</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
              {(dimensions.length > 0 ? dimensions : defaultDimensions).map((dim) => (
                <div
                  key={dim.name}
                  className="bg-gradient-to-br from-blue-50 to-indigo-50 backdrop-blur-sm rounded-xl p-4 border border-blue-100/50 hover:shadow-md transition-shadow"
                >
                  <h5 className="font-bold text-slate-800 mb-1">{dim.name}</h5>
                  <p className="text-xs text-slate-500 leading-relaxed">{dim.description}</p>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* 操作按钮 */}
        <div className="flex gap-4">
          <button
            onClick={handleCancel}
            disabled={submitting || dimensionsLoading}
            className="flex-1 py-4 bg-gray-100 hover:bg-gray-200 text-gray-700 font-bold rounded-xl transition-colors cursor-pointer disabled:opacity-50"
          >
            取消
          </button>
          <button
            onClick={handleAnalyze}
            disabled={submitting || dimensionsLoading}
            className="flex-[2] py-4 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-bold rounded-xl shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all cursor-pointer flex items-center justify-center gap-2 disabled:opacity-50 disabled:transform-none"
          >
            {submitting ? (
              <>
                <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                正在启动分析...
              </>
            ) : (
              <>
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                开始分析
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  )
}

export default VideoConfirm
