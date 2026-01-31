import { useEffect, useState } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'

interface VideoInfo {
  title: string
  cover: string
  author: string
  duration: string
}

interface AnalysisTarget {
  productName: string
  brand: string
  category: string
  features: string[]
}

const Confirm = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const url = searchParams.get('url')
  
  const [loading, setLoading] = useState(true)
  const [videoInfo, setVideoInfo] = useState<VideoInfo | null>(null)
  const [target, setTarget] = useState<AnalysisTarget | null>(null)

  useEffect(() => {
    // Simulate API call to fetch video info and initial analysis
    const timer = setTimeout(() => {
      setVideoInfo({
        title: "【何同学】快充伤电池？8000次循环测试真相...",
        cover: "https://i0.hdslb.com/bfs/archive/8d601321727771701387707770138770.jpg", // Mock cover
        author: "老师好我叫何同学",
        duration: "12:34"
      })
      setTarget({
        productName: "120W 快充充电器",
        brand: "小米",
        category: "数码配件/充电器",
        features: ["充电速度", "发热情况", "电池健康影响", "便携性"]
      })
      setLoading(false)
    }, 1500)

    return () => clearTimeout(timer)
  }, [url])

  const handleConfirm = () => {
    // Generate a random task ID and navigate to progress
    const taskId = 'task_' + Math.random().toString(36).substr(2, 9)
    navigate(`/progress/${taskId}`)
  }

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh]">
        <div className="w-16 h-16 border-4 border-blue-500/30 border-t-blue-500 rounded-full animate-spin mb-6"></div>
        <h2 className="text-2xl font-semibold text-gray-700">正在解析视频内容...</h2>
        <p className="text-gray-500 mt-2">AI 正在识别商品信息与评论数据</p>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="mb-8 text-center">
        <h1 className="text-3xl font-bold text-gray-800 mb-2">确认分析目标</h1>
        <p className="text-gray-500">请确认 AI 自动提取的商品信息是否准确</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* Left: Video Info */}
        <div className="glass-card">
          <h3 className="text-lg font-semibold text-gray-700 mb-4 border-b border-gray-200/50 pb-2">
            视频来源
          </h3>
          <div className="aspect-video bg-gray-200 rounded-lg mb-4 flex items-center justify-center overflow-hidden relative group">
             {/* Placeholder for video cover */}
             <div className="absolute inset-0 bg-gray-800/10 flex items-center justify-center text-gray-400">
                Mock Cover
             </div>
          </div>
          <h4 className="font-medium text-gray-900 line-clamp-2 mb-2">
            {videoInfo?.title}
          </h4>
          <div className="flex justify-between text-sm text-gray-500">
            <span>UP主: {videoInfo?.author}</span>
            <span>时长: {videoInfo?.duration}</span>
          </div>
          <div className="mt-4 pt-4 border-t border-gray-200/50 text-xs text-gray-400 break-all">
            {url}
          </div>
        </div>

        {/* Right: Analysis Target */}
        <div className="glass-card flex flex-col">
          <h3 className="text-lg font-semibold text-gray-700 mb-4 border-b border-gray-200/50 pb-2">
            分析维度
          </h3>
          
          <div className="space-y-6 flex-1">
            <div>
              <label className="text-xs font-bold text-gray-400 uppercase tracking-wider">商品名称</label>
              <div className="text-xl font-bold text-gray-800 mt-1">{target?.productName}</div>
            </div>
            
            <div className="grid grid-cols-2 gap-4">
               <div>
                <label className="text-xs font-bold text-gray-400 uppercase tracking-wider">品牌</label>
                <div className="text-lg text-gray-700 mt-1">{target?.brand}</div>
               </div>
               <div>
                <label className="text-xs font-bold text-gray-400 uppercase tracking-wider">类目</label>
                <div className="text-lg text-gray-700 mt-1">{target?.category}</div>
               </div>
            </div>

            <div>
              <label className="text-xs font-bold text-gray-400 uppercase tracking-wider">关注维度</label>
              <div className="flex flex-wrap gap-2 mt-2">
                {target?.features.map((f, i) => (
                  <span key={i} className="px-3 py-1 bg-blue-50 text-blue-600 rounded-full text-sm font-medium border border-blue-100">
                    {f}
                  </span>
                ))}
              </div>
            </div>
          </div>

          <button
            onClick={handleConfirm}
            className="w-full mt-8 py-4 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-bold rounded-xl shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all cursor-pointer"
          >
            开始深度分析
          </button>
        </div>
      </div>
    </div>
  )
}

export default Confirm
