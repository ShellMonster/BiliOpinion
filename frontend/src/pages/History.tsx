import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import Button from '../components/common/Button'

interface HistoryItem {
  id: number
  category: string
  videoCount: number
  commentCount: number
  status: string
  reportId: number
  createdAt: string
}

export default function History() {
  const navigate = useNavigate()
  const [histories, setHistories] = useState<HistoryItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    fetchHistories()
  }, [])

  const fetchHistories = async () => {
    try {
      setLoading(true)
      const response = await fetch('http://localhost:8080/api/history')
      if (!response.ok) throw new Error('Failed to fetch histories')
      const data = await response.json()
      setHistories(data || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这条历史记录吗？')) return

    try {
      const response = await fetch(`http://localhost:8080/api/history/${id}`, {
        method: 'DELETE'
      })
      if (!response.ok) throw new Error('Failed to delete history')
      
      setHistories(histories.filter(h => h.id !== id))
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Delete failed')
    }
  }

  const handleView = (reportId: number) => {
    navigate(`/report/${reportId}`)
  }

  const getStatusBadge = (status: string) => {
    const statusMap: Record<string, { text: string; color: string }> = {
      pending: { text: '待处理', color: 'bg-yellow-100 text-yellow-700' },
      processing: { text: '处理中', color: 'bg-blue-100 text-blue-700' },
      completed: { text: '已完成', color: 'bg-green-100 text-green-700' },
      failed: { text: '失败', color: 'bg-red-100 text-red-700' }
    }
    const { text, color } = statusMap[status] || { text: status, color: 'bg-gray-100 text-gray-700' }
    return <span className={`px-3 py-1 rounded-full text-sm font-medium ${color}`}>{text}</span>
  }

  if (loading) {
    return (
      <div className="text-center py-20">
        <div className="inline-block animate-spin rounded-full h-12 w-12 border-4 border-blue-600 border-t-transparent"></div>
        <p className="mt-4 text-slate-600">加载中...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-20">
        <p className="text-red-600 mb-4">❌ {error}</p>
        <Button onClick={fetchHistories}>重试</Button>
      </div>
    )
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-black text-slate-800 mb-2">历史记录</h1>
        <p className="text-slate-500">查看所有分析任务的历史记录</p>
      </div>

      {histories.length === 0 ? (
        <div className="text-center py-20 bg-slate-50 rounded-3xl">
          <p className="text-slate-500 text-lg">暂无历史记录</p>
          <Button className="mt-6" onClick={() => navigate('/')}>开始新分析</Button>
        </div>
      ) : (
        <div className="space-y-4">
          {histories.map(history => (
            <div 
              key={history.id}
              className="bg-white rounded-2xl shadow-sm border border-slate-200 p-6 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-3">
                    <h3 className="text-xl font-bold text-slate-800">{history.category}</h3>
                    {getStatusBadge(history.status)}
                  </div>
                  
                  <div className="flex gap-6 text-sm text-slate-600 mb-3">
                    <span>📹 视频: {history.videoCount}</span>
                    <span>💬 评论: {history.commentCount}</span>
                    <span>🕒 {history.createdAt}</span>
                  </div>
                </div>

                <div className="flex gap-2">
                  {history.status === 'completed' && history.reportId > 0 && (
                    <Button 
                      variant="primary" 
                      onClick={() => handleView(history.reportId)}
                      className="text-sm px-4 py-2"
                    >
                      查看报告
                    </Button>
                  )}
                  <Button 
                    variant="secondary" 
                    onClick={() => handleDelete(history.id)}
                    className="text-sm px-4 py-2"
                  >
                    删除
                  </Button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

