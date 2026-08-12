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

  it('rejects blank task_id', () => {
    expect(confirmNavigationTarget(true, { task_id: '   ' }, 'x').kind).toBe('error')
    expect(confirmNavigationTarget(true, { task_id: '' }, 'x').kind).toBe('error')
  })

  it('includes video mode in the href', () => {
    const result = confirmNavigationTarget(true, { task_id: 'vid-1' }, '评测', 'video')
    expect(result.kind).toBe('progress')
    if (result.kind === 'progress') {
      expect(result.href).toContain('mode=video')
      expect(result.href).not.toContain('搜索')
    }
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

  it('accepts a successful parse with brands', () => {
    const result = parseOutcomeFromResponse(true, { brands: ['戴森'], product_type: '吸尘器' })
    expect(result.kind).toBe('ok')
  })

  it('rejects empty brands and 200+error payloads', () => {
    expect(parseOutcomeFromResponse(true, { brands: [] }).kind).toBe('error')
    expect(parseOutcomeFromResponse(true, { error: '模型超时', brands: ['戴森'] }).kind).toBe('error')
  })
})

describe('beginSubmit', () => {
  it('rejects a second start while already submitting', () => {
    expect(beginSubmit(false)).toBe(true)
    expect(beginSubmit(true)).toBe(false)
  })
})
