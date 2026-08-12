import { describe, expect, it } from 'vitest'
import { progressStepsForMode, resolveProgressRestore } from './progressRestore'

describe('resolveProgressRestore', () => {
  it('sends completed tasks to the report', () => {
    expect(resolveProgressRestore({ status: 'completed', reportId: 42 })).toEqual({
      kind: 'report',
      reportId: 42,
    })
  })

  it('shows failed tasks as an error instead of hanging on SSE', () => {
    const result = resolveProgressRestore({
      status: 'failed',
      progressMsg: '请先配置AI API Key',
    })
    expect(result).toEqual({ kind: 'error', message: '请先配置AI API Key' })
  })

  it('resumes SSE for in-progress tasks', () => {
    expect(resolveProgressRestore({ status: 'processing', reportId: 0 })).toEqual({
      kind: 'sse',
    })
  })

  it('does not treat completed-without-report as a report navigation', () => {
    expect(resolveProgressRestore({ status: 'completed', reportId: 0 })).toEqual({
      kind: 'sse',
    })
  })
})

describe('progressStepsForMode', () => {
  it('uses video pipeline labels, not product search', () => {
    const steps = progressStepsForMode('video')
    expect(steps).toContain('解析视频信息')
    expect(steps).not.toContain('搜索相关视频')
  })
})
