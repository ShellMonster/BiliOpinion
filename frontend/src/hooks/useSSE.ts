import { useEffect, useState, useRef } from 'react'

interface SSEMessage {
  task_id: string
  status: string
  progress?: { current: number; total: number }
  message?: string
}

const MAX_RECONNECT_ATTEMPTS = 5

export function useSSE(url: string | null) {
  const [data, setData] = useState<SSEMessage | null>(null)
  const [error, setError] = useState<Error | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const reconnectCountRef = useRef(0)

  useEffect(() => {
    if (!url) return

    let eventSource: EventSource | null = null
    reconnectCountRef.current = 0

    const connect = () => {
      eventSource = new EventSource(url)

      eventSource.onopen = () => {
        setIsConnected(true)
        setError(null)
        reconnectCountRef.current = 0
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
        eventSource?.close()
        
        if (reconnectCountRef.current < MAX_RECONNECT_ATTEMPTS) {
          reconnectCountRef.current++
          const delay = Math.min(3000 * Math.pow(2, reconnectCountRef.current - 1), 30000)
          setTimeout(connect, delay)
        } else {
          setError(new Error('SSE connection failed after maximum retry attempts'))
        }
      }
    }

    connect()

    return () => {
      eventSource?.close()
    }
  }, [url])

  return { data, error, isConnected }
}
