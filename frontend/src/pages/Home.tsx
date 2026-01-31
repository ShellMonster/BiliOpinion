import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

const Home = () => {
  const [url, setUrl] = useState('')
  const navigate = useNavigate()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!url.trim()) return
    // Navigate to confirm page with the URL
    navigate(`/confirm?url=${encodeURIComponent(url)}`)
  }

  const examples = [
    "https://www.bilibili.com/video/BV1xx411c7mD",
    "https://www.bilibili.com/video/BV1GJ411x7h7",
    "BV1xx411c7mD"
  ]

  return (
    <div className="flex flex-col items-center justify-center min-h-[80vh] px-4 max-w-4xl mx-auto">
      <div className="text-center mb-12 animate-fade-in">
        <h1 className="text-5xl font-bold text-gray-800 mb-4 tracking-tight">
          Bilibili 评论分析
        </h1>
        <p className="text-xl text-gray-500 font-light">
          洞察用户反馈，提炼产品价值
        </p>
      </div>

      <form onSubmit={handleSubmit} className="w-full max-w-2xl relative">
        <div className="relative group">
          <div className="absolute -inset-1 bg-gradient-to-r from-blue-400 to-purple-400 rounded-2xl blur opacity-20 group-hover:opacity-40 transition duration-1000 group-hover:duration-200"></div>
          <div className="relative">
            <input
              type="text"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="输入 Bilibili 视频链接或 BV 号..."
              className="w-full p-6 pr-16 text-lg rounded-2xl border border-gray-200 shadow-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 bg-white/80 backdrop-blur-xl transition-all"
            />
            <button
              type="submit"
              disabled={!url.trim()}
              className="absolute right-3 top-3 bottom-3 aspect-square bg-gray-900 text-white rounded-xl hover:bg-black disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center cursor-pointer"
            >
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-6 h-6">
                <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
              </svg>
            </button>
          </div>
        </div>
      </form>

      <div className="mt-8 flex flex-wrap gap-3 justify-center text-sm text-gray-500">
        <span>试一试:</span>
        {examples.map((ex, i) => (
          <button
            key={i}
            onClick={() => setUrl(ex)}
            className="px-3 py-1 bg-white/50 hover:bg-white border border-transparent hover:border-gray-200 rounded-full transition-all cursor-pointer"
          >
            {ex.includes('http') ? '视频示例 ' + (i + 1) : ex}
          </button>
        ))}
      </div>
    </div>
  )
}

export default Home
