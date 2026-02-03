import { useEffect, useState, useRef } from 'react'
import { useParams, useNavigate, useSearchParams } from 'react-router-dom'

interface Step {
  id: number
  label: string
  status: 'pending' | 'processing' | 'completed' | 'error'
}

interface SSEData {
  task_id: string
  status: string
  message?: string
  error?: string
  progress?: {
    current: number
    total: number
    stage?: string
  }
}

const Progress = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const title = searchParams.get('title') || '分析任务'
  const eventSourceRef = useRef<EventSource | null>(null)

  const [progress, setProgress] = useState(0)
  const [message, setMessage] = useState('正在连接服务器...')
  const [error, setError] = useState<string | null>(null)
  const [steps, setSteps] = useState<Step[]>([
    { id: 1, label: '搜索相关视频', status: 'pending' },
    { id: 2, label: '抓取视频评论', status: 'pending' },
    { id: 3, label: 'AI 分析评论内容', status: 'pending' },
    { id: 4, label: '生成分析报告', status: 'pending' }
  ])

  useEffect(() => {
    if (!id) {
      console.error('[Progress] No task_id provided')
      setError('缺少任务ID')
      return
    }

    const eventSource = new EventSource(`http://localhost:8080/api/sse?task_id=${id}`)
    eventSourceRef.current = eventSource

    eventSource.onopen = () => {
    }

    eventSource.onmessage = (event) => {
      try {
        const data: SSEData = JSON.parse(event.data)
        
        if (data.message) {
          setMessage(data.message)
        }

        if (data.progress) {
          setProgress(data.progress.current)
        }

        updateStepsFromStatus(data.status, data.progress?.current || 0)

        if (data.status === 'completed') {
          const reportId = data.progress?.stage
          eventSource.close()
          setTimeout(() => {
            navigate(`/report/${reportId}`)
          }, 1000)
        }

        if (data.status === 'error') {
          console.error('[Progress] Task error:', data.error || data.message)
          setError(data.error || data.message || '任务执行失败')
          eventSource.close()
        }
      } catch (e) {
        console.error('[Progress] Failed to parse SSE data:', e, event.data)
      }
    }

    eventSource.onerror = (err) => {
      console.error('[Progress] SSE connection error:', err)
      setError('连接中断，请刷新页面重试')
      eventSource.close()
    }

    return () => {
      eventSource.close()
    }
  }, [id, navigate])

  const updateStepsFromStatus = (status: string, progressValue: number) => {
    setSteps(prev => {
      const newSteps = [...prev]
      
      if (status === 'searching' || progressValue < 20) {
        newSteps[0].status = 'processing'
      } else if (status === 'scraping' || (progressValue >= 20 && progressValue < 50)) {
        newSteps[0].status = 'completed'
        newSteps[1].status = 'processing'
      } else if (status === 'analyzing' || (progressValue >= 50 && progressValue < 85)) {
        newSteps[0].status = 'completed'
        newSteps[1].status = 'completed'
        newSteps[2].status = 'processing'
      } else if (status === 'generating' || progressValue >= 85) {
        newSteps[0].status = 'completed'
        newSteps[1].status = 'completed'
        newSteps[2].status = 'completed'
        newSteps[3].status = 'processing'
      }
      
      if (status === 'completed') {
        return newSteps.map(s => ({ ...s, status: 'completed' as const }))
      }
      
      if (status === 'error') {
        return newSteps.map(s => 
          s.status === 'processing' ? { ...s, status: 'error' as const } : s
        )
      }
      
      return newSteps
    })
  }

  if (error) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-12">
        <div className="glass-card text-center py-12">
          <div className="text-6xl mb-6">❌</div>
          <h2 className="text-2xl font-bold text-red-600 mb-4">任务执行失败</h2>
          <p className="text-gray-600 mb-6">{error}</p>
          <button
            onClick={() => navigate('/')}
            className="px-6 py-2 bg-gray-800 text-white rounded-lg hover:bg-gray-700 transition-colors cursor-pointer"
          >
            返回首页
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-3xl mx-auto px-4 py-12">
      <div className="glass-card text-center py-12">
        <div className="mb-8">
          <h1 className="text-2xl font-semibold text-gray-900 mb-6">{title}</h1>
          <div className="relative w-48 h-48 mx-auto mb-6 flex items-center justify-center">
            <svg className="w-full h-full transform -rotate-90" viewBox="0 0 100 100">
              <circle cx="50" cy="50" r="45" fill="none" stroke="#eee" strokeWidth="8" />
              <circle 
                cx="50" cy="50" r="45" fill="none" stroke="url(#gradient)" strokeWidth="8" 
                strokeDasharray="283"
                strokeDashoffset={283 - (283 * progress) / 100}
                className="transition-all duration-500 ease-out"
                strokeLinecap="round"
              />
              <defs>
                <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                  <stop offset="0%" stopColor="#3b82f6" />
                  <stop offset="100%" stopColor="#8b5cf6" />
                </linearGradient>
              </defs>
            </svg>
            <div className="absolute inset-0 flex flex-col items-center justify-center">
              <span className="text-4xl font-bold text-gray-800">{Math.round(progress)}%</span>
              <span className="text-sm text-gray-500 mt-1">处理中</span>
            </div>
          </div>
          
          <h2 className="text-2xl font-bold text-gray-800 mb-2">{message}</h2>
          <p className="text-gray-500 text-sm">Task ID: {id}</p>
        </div>

        <div className="max-w-md mx-auto space-y-4 text-left">
          {steps.map((step, index) => (
            <div key={step.id} className="flex items-center gap-4">
              <div className={`
                w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold transition-colors duration-300
                ${step.status === 'completed' ? 'bg-green-500 text-white' : 
                  step.status === 'processing' ? 'bg-blue-500 text-white animate-pulse' : 
                  step.status === 'error' ? 'bg-red-500 text-white' :
                  'bg-gray-200 text-gray-400'}
              `}>
                {step.status === 'completed' ? '✓' : 
                 step.status === 'error' ? '!' : index + 1}
              </div>
              <div className="flex-1">
                <div className={`font-medium ${
                  step.status === 'pending' ? 'text-gray-400' : 
                  step.status === 'error' ? 'text-red-600' : 'text-gray-800'
                }`}>
                  {step.label}
                </div>
              </div>
              {step.status === 'processing' && (
                <div className="w-4 h-4 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export default Progress
