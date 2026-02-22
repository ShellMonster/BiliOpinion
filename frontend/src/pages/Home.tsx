import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import VideoInput from '../components/VideoInput'

const Home = () => {
  const [activeTab, setActiveTab] = useState<'product' | 'video'>('product')
  const [requirement, setRequirement] = useState('')
  const navigate = useNavigate()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!requirement.trim()) return
    navigate(`/confirm?requirement=${encodeURIComponent(requirement.trim())}`)
  }

  const examples = [
    "机械键盘，打游戏用",
    "投影仪，卧室用",
    "空气炸锅，一个人用",
    "吸尘器，有宠物",
    "蓝牙耳机，通勤降噪"
  ]

  return (
    <div className="flex flex-col items-center justify-center min-h-[80vh] px-4">
      <div className="w-full max-w-3xl mx-auto">
        <div className="text-center mb-12">
          <h1 className="text-4xl md:text-5xl font-semibold text-gray-900 mb-4 tracking-tight">
            B站商品评论分析
          </h1>
          <p className="text-lg text-gray-500">
            告诉我你想买什么，AI 帮你分析真实用户评价
          </p>
        </div>

        {/* Tab 切换按钮区域 - 使用 flex 布局居中 */}
        <div className="flex justify-center mb-8">
          <div className="inline-flex bg-gray-100 p-1 rounded-xl">
            {/* 商品需求 Tab 按钮 */}
            <button
              onClick={() => setActiveTab('product')}
              className={`px-6 py-2.5 rounded-lg text-sm font-medium transition-all duration-200 cursor-pointer ${
                activeTab === 'product'
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <span className="mr-1.5">🔍</span>
              商品需求
            </button>
            {/* 视频链接 Tab 按钮 */}
            <button
              onClick={() => setActiveTab('video')}
              className={`px-6 py-2.5 rounded-lg text-sm font-medium transition-all duration-200 cursor-pointer ${
                activeTab === 'video'
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <span className="mr-1.5">📺</span>
              视频链接
            </button>
          </div>
        </div>

        {/* 根据选中的 Tab 显示不同内容 */}
        {activeTab === 'product' ? (
          /* 商品需求 Tab - 保持现有功能不变 */
          <form onSubmit={handleSubmit} className="mb-6">
          <div className="relative">
            <input
              type="text"
              value={requirement}
              onChange={(e) => setRequirement(e.target.value)}
              placeholder="描述你的需求，比如：想买个吸尘器，预算2000，家里有宠物..."
              className="w-full px-5 py-4 pr-14 text-base rounded-xl border border-gray-200 focus:outline-none focus:border-gray-400 focus:ring-1 focus:ring-gray-400 bg-white transition-colors"
            />
            <button
              type="submit"
              disabled={!requirement.trim()}
              className="absolute right-2 top-1/2 -translate-y-1/2 p-2.5 bg-black text-white rounded-lg hover:bg-gray-800 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors cursor-pointer"
            >
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-5 h-5">
                <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
              </svg>
            </button>
          </div>
        </form>
        ) : (
          /* 视频链接 Tab - 显示占位符组件 */
          <div className="mb-6">
            <VideoInput />
          </div>
        )}

        {activeTab === 'product' && (
        <div className="flex items-center justify-center gap-x-1 gap-y-2 flex-wrap mb-16">
          <span className="text-sm text-gray-400">试试这些:</span>
          {examples.map((text, i) => (
            <button
              key={i}
              onClick={() => setRequirement(text)}
              className="px-3 py-1.5 text-sm text-gray-500 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors cursor-pointer"
            >
              {text}
          </button>
        ))}
        </div>
        )}

        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
          <div>
            <div className="text-2xl mb-2">🤖</div>
            <div className="text-sm font-medium text-gray-900">AI 智能解析</div>
            <div className="text-xs text-gray-400 mt-1">自动提取品牌和维度</div>
          </div>
          <div>
            <div className="text-2xl mb-2">📊</div>
            <div className="text-sm font-medium text-gray-900">多维度分析</div>
            <div className="text-xs text-gray-400 mt-1">6大维度全面评估</div>
          </div>
          <div>
            <div className="text-2xl mb-2">🏆</div>
            <div className="text-sm font-medium text-gray-900">品牌排名</div>
            <div className="text-xs text-gray-400 mt-1">综合评分一目了然</div>
          </div>
          <div>
            <div className="text-2xl mb-2">📈</div>
            <div className="text-sm font-medium text-gray-900">型号对比</div>
            <div className="text-xs text-gray-400 mt-1">精准到具体型号</div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Home
