import { Link } from 'react-router-dom'
import { type ReactNode } from 'react'

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="min-h-screen bg-[#f8fafc]">
      {/* Header */}
      <header className="bg-white/70 backdrop-blur-xl border-b border-white/40 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Link to="/" className="text-2xl font-black text-slate-800">
              B站商品评论分析
            </Link>
            <nav className="flex gap-6">
              <Link to="/" className="text-slate-600 hover:text-slate-900 font-medium">
                首页
              </Link>
              <Link to="/history" className="text-slate-600 hover:text-slate-900 font-medium">
                历史记录
              </Link>
              <Link to="/settings" className="text-slate-600 hover:text-slate-900 font-medium">
                设置
              </Link>
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {children}
      </main>
    </div>
  )
}
