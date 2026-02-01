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
    if (!data || !requirement) return

    try {
      const response = await fetch('http://localhost:8080/api/confirm', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          requirement: requirement,
          brands: data.brands,
          dimensions: data.dimensions,
          keywords: data.keywords
        })
      })
      const result = await response.json()
      navigate(`/progress/${result.task_id}`)
    } catch (error) {
      console.error('Failed to confirm:', error)
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

            {/* Brand Tags */}
            <div>
                <h4 className="text-sm font-bold text-gray-600 mb-4 flex items-center gap-2">
                    <span>🏷️</span> 将分析这些品牌
                </h4>
                <div className="flex flex-wrap gap-3">
                {data.brands.map(brand => (
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
                {data.dimensions.map(dim => (
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
                    {data.keywords.join(' | ')}
                </div>
            </div>
        </div>

        {/* Confirm Button */}
        <button
          onClick={handleConfirm}
          className="w-full py-4 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-bold rounded-xl shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all cursor-pointer flex items-center justify-center gap-2"
        >
          <span>✓</span> 确认开始分析
        </button>
      </div>
    </div>
  )
}

export default Confirm
