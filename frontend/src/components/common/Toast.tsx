import { useEffect, useState } from 'react'

interface ToastProps {
  message: string
  type?: 'success' | 'error' | 'info'
  onClose: () => void
}

export default function Toast({ message, type = 'info', onClose }: ToastProps) {
  const [isVisible, setIsVisible] = useState(false)

  useEffect(() => {
    // Trigger enter animation
    requestAnimationFrame(() => setIsVisible(true))
  }, [])

  const bgColor = {
    success: 'bg-green-500',
    error: 'bg-red-500',
    info: 'bg-blue-500'
  }[type]

  const icon = {
    success: '✓',
    error: '✕',
    info: 'ℹ'
  }[type]

  return (
    <div 
      className={`transform transition-all duration-300 ease-in-out ${
        isVisible ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'
      } pointer-events-auto flex items-center gap-3 px-6 py-4 bg-white/90 backdrop-blur-xl rounded-2xl shadow-xl border border-white/50 min-w-[300px]`}
      onClick={onClose}
      role="alert"
    >
      <span className={`w-6 h-6 ${bgColor} text-white rounded-full flex items-center justify-center text-sm font-bold shrink-0`}>
        {icon}
      </span>
      <span className="text-slate-800 font-medium">{message}</span>
    </div>
  )
}
