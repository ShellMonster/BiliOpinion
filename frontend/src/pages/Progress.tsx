import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
// import { useSSE } from '../hooks/useSSE' // Commented out to use mock implementation for now
// To demonstrate the UI without backend, I will use a local simulation that mimics the SSE structure

interface Step {
  id: number
  label: string
  status: 'pending' | 'processing' | 'completed'
}

const Progress = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  // const { data } = useSSE(\`/api/progress/\${id}\`) // Real implementation

  const [progress, setProgress] = useState(0)
  const [steps, setSteps] = useState<Step[]>([
    { id: 1, label: '抓取视频评论数据', status: 'pending' },
    { id: 2, label: '清洗与预处理文本', status: 'pending' },
    { id: 3, label: 'AI 情感与观点分析', status: 'pending' },
    { id: 4, label: '生成多维度报表', status: 'pending' }
  ])

  // Mock simulation of progress
  useEffect(() => {
    let p = 0
    let step = 0
    
    const interval = setInterval(() => {
      p += Math.random() * 5
      if (p > 100) p = 100
      setProgress(p)

      // Update steps based on progress
      if (p < 25) step = 0
      else if (p < 50) step = 1
      else if (p < 85) step = 2
      else step = 3
      
      setSteps(prev => prev.map((s, i) => ({
        ...s,
        status: i < step ? 'completed' : i === step ? 'processing' : 'pending'
      })))

      if (p >= 100) {
        clearInterval(interval)
        setTimeout(() => {
          navigate(`/report/${id}`)
        }, 1000)
      }
    }, 500)

    return () => clearInterval(interval)
  }, [id, navigate])

  return (
    <div className="max-w-3xl mx-auto px-4 py-12">
      <div className="glass-card text-center py-12">
        <div className="mb-8">
          <div className="relative w-48 h-48 mx-auto mb-6 flex items-center justify-center">
             {/* Circular Progress */}
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
          
          <h2 className="text-2xl font-bold text-gray-800 mb-2">正在深入分析评论内容</h2>
          <p className="text-gray-500">Task ID: {id}</p>
        </div>

        <div className="max-w-md mx-auto space-y-4 text-left">
          {steps.map((step, index) => (
            <div key={step.id} className="flex items-center gap-4">
              <div className={`
                w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold transition-colors duration-300
                ${step.status === 'completed' ? 'bg-green-500 text-white' : 
                  step.status === 'processing' ? 'bg-blue-500 text-white animate-pulse' : 
                  'bg-gray-200 text-gray-400'}
              `}>
                {step.status === 'completed' ? '✓' : index + 1}
              </div>
              <div className="flex-1">
                <div className={`font-medium ${step.status === 'pending' ? 'text-gray-400' : 'text-gray-800'}`}>
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
