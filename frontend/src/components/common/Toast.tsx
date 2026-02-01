import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'

interface ToastProps {
  message: string
  type?: 'success' | 'error' | 'info'
  duration?: number
  onClose: () => void
}

export default function Toast({ message, type = 'success', duration = 3000, onClose }: ToastProps) {
  const [isVisible, setIsVisible] = useState(true)

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsVisible(false)
      setTimeout(onClose, 300) // Wait for fade out animation
    }, duration)

    return () => clearTimeout(timer)
  }, [duration, onClose])

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

  return createPortal(
    <div 
      className={`fixed top-6 left-1/2 -translate-x-1/2 z-[9999] transition-all duration-300 ${
        isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-4'
      }`}
    >
      <div className="flex items-center gap-3 px-6 py-4 bg-white/90 backdrop-blur-xl rounded-2xl shadow-2xl border border-white/50">
        <span className={`w-6 h-6 ${bgColor} text-white rounded-full flex items-center justify-center text-sm font-bold`}>
          {icon}
        </span>
        <span className="text-slate-800 font-medium">{message}</span>
      </div>
    </div>,
    document.body
  )
}
