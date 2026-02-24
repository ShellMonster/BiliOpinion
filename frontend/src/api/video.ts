import { apiClient } from './client'

export interface Dimension {
  name: string
  description: string
}

export interface VideoParseResponse {
  bvid: string
  title: string
  author: string
  play_count: number
  comment_count: number
  pub_date: string
  cover: string
  description?: string
}

export interface VideoAnalyzeResponse {
  task_id: string
  message: string
}

export interface DimensionsRequest {
  bvid: string
}

export interface DimensionsResponse {
  dimensions: Dimension[]
}

export async function getDimensions(bvid: string): Promise<DimensionsResponse> {
  return apiClient.post<DimensionsResponse>('/video/dimensions', { bvid })
}

export async function parseVideo(video_url: string): Promise<VideoParseResponse> {
  return apiClient.post<VideoParseResponse>('/video/parse', { video_url })
}

export async function analyzeVideo(
  video_url: string, 
  max_comments: number,
  dimensions?: Dimension[]  // 可选的分析维度
): Promise<VideoAnalyzeResponse> {
  return apiClient.post<VideoAnalyzeResponse>('/video/analyze', { 
    video_url, 
    max_comments,
    dimensions  // 传递维度到后端
  })
}
