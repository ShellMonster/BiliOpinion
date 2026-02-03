import { Link } from 'react-router-dom'
import { type ReactNode, useState } from 'react'
import SettingsModal from '../Settings/SettingsModal'

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const [isSettingsOpen, setIsSettingsOpen] = useState(false)

  return (
    <div className="min-h-screen bg-[#fafafa]">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Link to="/" className="text-xl font-semibold text-gray-900">
              B站商品评论分析
            </Link>
            <nav className="flex gap-6">
              <Link to="/" className="text-gray-500 hover:text-gray-900 font-medium transition-colors">
                首页
              </Link>
              <Link to="/history" className="text-gray-500 hover:text-gray-900 font-medium transition-colors">
                历史记录
              </Link>
              <button 
                type="button"
                onClick={() => setIsSettingsOpen(true)} 
                className="text-gray-500 hover:text-gray-900 font-medium cursor-pointer bg-transparent border-none transition-colors"
              >
                设置
              </button>
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {children}
      </main>
      
      <SettingsModal 
        isOpen={isSettingsOpen} 
        onClose={() => setIsSettingsOpen(false)} 
      />
    </div>
  )
}
