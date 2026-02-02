// 品牌排名信息
export interface BrandRanking {
  brand: string
  overall_score: number
  rank: number
  scores: Record<string, number>
}

// 型号排名信息
export interface ModelRanking {
  model: string
  brand: string
  overall_score: number
  rank: number
  scores: Record<string, number>
  comment_count: number
}

// 评价维度
export interface Dimension {
  name: string
  description: string
}

// 报告统计数据
export interface ReportStats {
  total_videos: number
  total_comments: number
  comments_by_brand: Record<string, number>
}

// 典型评论
export interface TypicalComment {
  content: string
  score: number
}

// 品牌优劣势分析
export interface BrandAnalysis {
  strengths: string[]
  weaknesses: string[]
}

// 视频来源信息
export interface VideoSource {
  bvid: string
  title: string
  author: string
  play: number
  video_review: number
}

// 关键词条目
export interface KeywordItem {
  word: string
  count: number
}

// 情感统计
export interface SentimentStats {
  positive_count: number
  neutral_count: number
  negative_count: number
  positive_pct: number
  neutral_pct: number
  negative_pct: number
}

// 报告数据结构
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
  video_sources?: VideoSource[]
}

// API 响应结构
export interface ApiResponse {
  id: number
  history_id: number
  category: string
  data: ReportData
  created_at: string
}
