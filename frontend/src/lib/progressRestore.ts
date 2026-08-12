export type HistorySnapshot = {
  status: string
  reportId?: number
  progressMsg?: string
}

export type RestoreAction =
  | { kind: 'report'; reportId: number }
  | { kind: 'error'; message: string }
  | { kind: 'sse' }

export function resolveProgressRestore(snapshot: HistorySnapshot | null): RestoreAction {
  if (!snapshot) {
    return { kind: 'sse' }
  }
  if (snapshot.status === 'completed' && snapshot.reportId && snapshot.reportId > 0) {
    return { kind: 'report', reportId: snapshot.reportId }
  }
  if (snapshot.status === 'failed') {
    return {
      kind: 'error',
      message: snapshot.progressMsg || '任务执行失败',
    }
  }
  return { kind: 'sse' }
}

export const PRODUCT_PROGRESS_STEPS = [
  '搜索相关视频',
  '抓取视频评论',
  'AI 分析评论内容',
  '生成分析报告',
] as const

export const VIDEO_PROGRESS_STEPS = [
  '解析视频信息',
  '抓取视频评论',
  'AI 分析评论内容',
  '生成分析报告',
] as const

export function progressStepsForMode(mode: 'product' | 'video'): readonly string[] {
  return mode === 'video' ? VIDEO_PROGRESS_STEPS : PRODUCT_PROGRESS_STEPS
}
