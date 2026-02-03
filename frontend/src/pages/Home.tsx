import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

const Home = () => {
  const [requirement, setRequirement] = useState('')
  const navigate = useNavigate()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!requirement.trim()) return
    navigate(`/confirm?requirement=${encodeURIComponent(requirement.trim())}`)
  }

  const features = [
    {
      icon: '🤖',
      title: 'AI 智能解析',
      desc: '自动提取品牌、维度和搜索关键词'
    },
    {
      icon: '📊',
      title: '多维度分析',
      desc: '6大评价维度，全面了解产品表现'
    },
    {
      icon: '🏆',
      title: '品牌排名',
      desc: '综合评分排序，最佳选择一目了然'
    },
    {
      icon: '📈',
      title: '型号对比',
      desc: '精准到具体型号，不再选择困难'
    }
  ]

  const examples = [
    { text: "机械键盘，打游戏用", icon: "⌨️" },
    { text: "投影仪，卧室用", icon: "📽️" },
    { text: "空气炸锅，一个人用", icon: "🍳" },
    { text: "吸尘器，有宠物", icon: "🧹" },
    { text: "蓝牙耳机，通勤降噪", icon: "🎧" }
  ]

  return (
    <div className="min-h-[85vh] bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 -mx-4 sm:-mx-6 lg:-mx-8 px-4 sm:px-6 lg:px-8 py-12">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-12">
          <div className="text-6xl mb-6 animate-pulse">
            📊
          </div>
          
          <h1 className="text-4xl md:text-5xl font-bold text-gray-800 mb-4">
            Bilibili 商品评论分析
          </h1>
          
          <p className="text-lg md:text-xl text-gray-500 max-w-2xl mx-auto leading-relaxed">
            基于 AI 的真实用户评价分析工具
            <br />
            <span className="text-base">帮你从海量B站评论中找到最值得买的产品</span>
          </p>
        </div>

        <form onSubmit={handleSubmit} className="w-full max-w-2xl mx-auto mb-8">
          <div className="relative group">
            <div className="absolute -inset-1 bg-gradient-to-r from-blue-400 to-purple-400 rounded-2xl blur opacity-20 group-hover:opacity-40 transition duration-1000 group-hover:duration-200"></div>
            <div className="relative">
              <input
                type="text"
                value={requirement}
                onChange={(e) => setRequirement(e.target.value)}
                placeholder="描述你的需求，比如：想买个吸尘器，预算2000，家里有宠物..."
                className="w-full p-6 pr-16 text-lg rounded-2xl border border-gray-200 shadow-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 bg-white/80 backdrop-blur-xl transition-all"
              />
              <button
                type="submit"
                disabled={!requirement.trim()}
                className="absolute right-3 top-3 bottom-3 aspect-square bg-gray-900 text-white rounded-xl hover:bg-black disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center cursor-pointer"
              >
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-6 h-6">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
                </svg>
              </button>
            </div>
          </div>
        </form>

        <div className="flex flex-wrap gap-3 justify-center mb-16">
          <span className="text-gray-400 text-sm">💡 试试这些:</span>
          {examples.map((ex, i) => (
            <button
              key={i}
              onClick={() => setRequirement(ex.text)}
              className="px-4 py-2 bg-white/70 hover:bg-white rounded-full shadow-sm hover:shadow-md transition-all text-sm text-gray-600 hover:text-gray-800 border border-transparent hover:border-gray-200"
            >
              {ex.icon} {ex.text}
            </button>
          ))}
        </div>

        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
          {features.map((f, i) => (
            <div
              key={i}
              className="bg-white rounded-2xl p-6 shadow-lg hover:shadow-xl hover:-translate-y-1 transition-all duration-300"
            >
              <div className="text-4xl mb-4">{f.icon}</div>
              <h3 className="font-bold text-gray-800 mb-2">{f.title}</h3>
              <p className="text-sm text-gray-500">{f.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export default Home
