import { describe, expect, it } from 'vitest'
import { settingsSaveOutcome } from './settingsSave'

describe('settingsSaveOutcome', () => {
  it('treats a non-OK response as failure', () => {
    expect(settingsSaveOutcome(false)).toEqual({ kind: 'error', message: '保存失败' })
  })

  it('treats OK as success', () => {
    expect(settingsSaveOutcome(true)).toEqual({ kind: 'ok' })
  })
})
