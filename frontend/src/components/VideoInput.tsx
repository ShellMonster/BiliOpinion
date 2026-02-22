import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

// VideoInput 组件：B站视频链接输入组件
// 功能：接收用户输入的视频URL，验证格式，通过后跳转到确认页
const VideoInput = () => {
  // url: 输入框的值
  const [url, setUrl] = useState('')
  // error: 错误提示信息，为空时表示没有错误
  const [error, setError] = useState('')
  // navigate: React Router的导航函数，用于跳转页面
  const navigate = useNavigate()

  /**
   * 验证URL是否为有效的B站视频链接
   * 支持的格式：
   * - bilibili.com/video/BV... (PC端)
   * - bilibili.com/video/av... (PC端旧格式)
   * - m.bilibili.com/video/BV... (移动端)
   * 
   * 注意：不支持 b23.tv 短链接
   * 
   * @param inputUrl 用户输入的URL字符串
   * @returns 是否有效
   */
  const isValidBilibiliUrl = (inputUrl: string): boolean => {
    // 去除首尾空格
    const trimmed = inputUrl.trim()
    
    // 定义有效的URL匹配模式
    const validPatterns = [
      /bilibili\.com\/video\/BV/i,     // PC端 BV号
      /bilibili\.com\/video\/av/i,     // PC端 av号
      /m\.bilibili\.com\/video\/BV/i,  // 移动端 BV号
    ]
    
    // 只要匹配任意一个模式就认为是有效的
    return validPatterns.some(pattern => pattern.test(trimmed))
  }

  /**
   * 处理表单提交
   * 1. 阻止默认表单行为
   * 2. 验证URL格式
   * 3. 验证失败显示错误
   * 4. 验证成功跳转到确认页
   * 
   * @param e 表单提交事件
   */
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    // 清空之前的错误信息
    setError('')
    
    // 检查输入是否为空
    if (!url.trim()) {
      setError('请输入视频链接')
      return
    }
    
    // 验证URL格式
    if (!isValidBilibiliUrl(url)) {
      setError('请输入有效的B站视频链接')
      return
    }
    
    // 验证通过，跳转到确认页
    // 使用 encodeURIComponent 对URL进行编码，防止特殊字符导致的问题
    navigate(`/video-confirm?url=${encodeURIComponent(url.trim())}`)
  }

  return (
    <div className="w-full">
      {/* 表单容器 */}
      <form onSubmit={handleSubmit}>
        {/* 输入框容器：相对定位，用于放置按钮 */}
        <div className="relative">
          {/* 文本输入框 */}
          <input
            type="text"
            value={url}
            onChange={(e) => {
              setUrl(e.target.value)
              // 用户输入时清空错误提示
              if (error) setError('')
            }}
            placeholder="粘贴B站视频链接..."
            className="w-full px-5 py-4 pr-14 text-base rounded-xl border border-gray-200 focus:outline-none focus:border-gray-400 focus:ring-1 focus:ring-gray-400 bg-white transition-colors"
          />
          
          {/* 提交按钮：绝对定位在输入框右侧 */}
          <button
            type="submit"
            disabled={!url.trim()}
            className="absolute right-2 top-1/2 -translate-y-1/2 p-2.5 bg-black text-white rounded-lg hover:bg-gray-800 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors cursor-pointer"
            aria-label="提交"
          >
            {/* 箭头图标 */}
            <svg 
              xmlns="http://www.w3.org/2000/svg" 
              fill="none" 
              viewBox="0 0 24 24" 
              strokeWidth={2} 
              stroke="currentColor" 
              className="w-5 h-5"
            >
              <path 
                strokeLinecap="round" 
                strokeLinejoin="round" 
                d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" 
              />
            </svg>
          </button>
        </div>
      </form>
      
      {/* 支持的格式提示 */}
      <p className="mt-2 text-sm text-gray-400">
        支持格式：bilibili.com/video/BV... 或 m.bilibili.com/video/BV...
      </p>
      
      {/* 错误提示：当有错误时显示红色文字 */}
      {error && (
        <p className="mt-2 text-sm text-red-500">
          {error}
        </p>
      )}
    </div>
  )
}

export default VideoInput
