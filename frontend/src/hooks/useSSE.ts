import { useEffect, useState } from 'react'

interface SSEMessage {
  task_id: string
  status: string
  progress?: { current: number; total: number }
  message?: string
}

export function useSSE(url: string | null) {
  const [data, setData] = useState<SSEMessage | null>(null)
  const [error, setError] = useState<Error | null>(null)
  const [isConnected, setIsConnected] = useState(false)

  useEffect(() => {
    if (!url) return

    let eventSource: EventSource | null = null

    const connect = () => {
      eventSource = new EventSource(url)

      eventSource.onopen = () => {
        setIsConnected(true)
        setError(null)
      }

      eventSource.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          setData(message)
        } catch (err) {
          console.error('Failed to parse SSE message:', err)
        }
      }

      eventSource.onerror = () => {
        setIsConnected(false)
        setError(new Error('SSE connection error'))
        eventSource?.close()
        
        // 自动重连（3秒后）
        setTimeout(connect, 3000)
      }
    }

    connect()

    return () => {
      eventSource?.close()
    }
  }, [url])

  return { data, error, isConnected }
}
