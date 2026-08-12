import { describe, expect, it } from 'vitest'
import { beginSubmit, confirmNavigationTarget, parseOutcomeFromResponse } from './confirmFlow'

describe('confirmNavigationTarget', () => {
  it('does not navigate when HTTP is not OK', () => {
    const result = confirmNavigationTarget(false, { error: '品牌列表不能为空' }, '吸尘器')
    expect(result.kind).toBe('error')
    if (result.kind === 'error') {
      expect(result.message).toContain('品牌列表不能为空')
    }
  })

  it('does not navigate to /progress/undefined when task_id is missing', () => {
    const result = confirmNavigationTarget(true, {}, '吸尘器')
    expect(result.kind).toBe('error')
    if (result.kind === 'progress') {
      throw new Error('must not produce a progress href')
    }
    expect(result.message).toContain('任务ID')
  })

  it('builds a progress href when task_id is present', () => {
    const result = confirmNavigationTarget(true, { task_id: 'abc-123' }, '吸尘器')
    expect(result).toEqual({
      kind: 'progress',
      href: '/progress/abc-123?title=%E5%90%B8%E5%B0%98%E5%99%A8',
    })
  })
})

describe('parseOutcomeFromResponse', () => {
  it('treats non-OK parse as an error, not a plan', () => {
    const result = parseOutcomeFromResponse(false, { error: 'AI API密钥未配置，请先在设置页面配置' })
    expect(result.kind).toBe('error')
    if (result.kind === 'error') {
      expect(result.message).toContain('AI API密钥未配置')
    }
  })
})

describe('beginSubmit', () => {
  it('rejects a second start while already submitting', () => {
    expect(beginSubmit(false)).toBe(true)
    expect(beginSubmit(true)).toBe(false)
  })
})
