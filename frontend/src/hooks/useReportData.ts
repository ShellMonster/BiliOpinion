import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'

export interface BrandRanking {
  brand: string
  overall_score: number
  rank: number
  scores: Record<string, number>
}

export interface ModelRanking {
  model: string
  brand: string
  overall_score: number
  rank: number
  scores: Record<string, number>
  comment_count: number
}

export interface Dimension {
  name: string
  description: string
}

export interface ReportStats {
  total_videos: number
  total_comments: number
  comments_by_brand: Record<string, number>
}

export interface TypicalComment {
  content: string
  score: number
}

export interface BrandAnalysis {
  strengths: string[]
  weaknesses: string[]
}

export interface VideoSource {
  bvid: string
  title: string
  author: string
  play: number
  video_review: number
}

export interface KeywordItem {
  word: string
  count: number
}

export interface SentimentStats {
  positive_count: number
  neutral_count: number
  negative_count: number
  positive_pct: number
  neutral_pct: number
  negative_pct: number
}

export interface ReportData {
  category: string
  brands: string[]
  dimensions: Dimension[]
  scores: Record<string, Record<string, number>>
  rankings: BrandRanking[]
  model_rankings?: ModelRanking[]
  recommendation: string
  stats?: ReportStats
  top_comments?: Record<string, TypicalComment[]>
  bad_comments?: Record<string, TypicalComment[]>
  brand_analysis?: Record<string, BrandAnalysis>
  sentiment_distribution?: Record<string, SentimentStats>
  video_sources?: VideoSource[]
  keywords?: KeywordItem[]
}

export interface ApiResponse {
  id: number
  history_id: number
  category: string
  data: ReportData
  created_at: string
}

/**
 * 获取报告数据的 Hook
 * 包含报告详情和历史记录中的品牌信息
 */
export function useReportData() {
  const { id } = useParams()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [report, setReport] = useState<ApiResponse | null>(null)
  const [specifiedBrands, setSpecifiedBrands] = useState<string[]>([])

  useEffect(() => {
    if (!id) return

    const fetchReport = async () => {
      try {
        setLoading(true)
        // 获取报告详情
        const response = await fetch(`http://localhost:8080/api/report/${id}`)
        if (!response.ok) {
          throw new Error('报告不存在')
        }
        const data = await response.json()
        setReport(data)

        // 如果有关联的历史记录，获取历史记录信息（主要是品牌列表）
        if (data.history_id) {
          try {
            const historyRes = await fetch(`http://localhost:8080/api/history/${data.history_id}`)
            if (historyRes.ok) {
              const historyData = await historyRes.json()
              setSpecifiedBrands(historyData.brands || [])
            }
          } catch (e) {
            console.error("Failed to fetch history brands", e)
          }
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : '加载报告失败')
      } finally {
        setLoading(false)
      }
    }

    fetchReport()
  }, [id])

  return {
    loading,
    error,
    report,
    specifiedBrands,
    id
  }
}
