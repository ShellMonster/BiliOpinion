import { describe, expect, it } from 'vitest'
import { scoreTone, scoreToneTextClass } from './scoreTone'

describe('scoreTone', () => {
  it('maps 8.5 to a success tone, not the 90-based red path', () => {
    expect(scoreTone(8.5)).toBe('success')
    expect(scoreToneTextClass(8.5)).toContain('emerald')
    expect(scoreToneTextClass(8.5)).not.toContain('rose')
  })

  it('uses 1-10 thresholds', () => {
    expect(scoreTone(8)).toBe('success')
    expect(scoreTone(6)).toBe('info')
    expect(scoreTone(4)).toBe('warning')
    expect(scoreTone(3.9)).toBe('danger')
  })
})
