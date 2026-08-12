import { describe, expect, it } from 'vitest'
import { addListItem, buildConfirmPayload, removeListItem } from './confirmPlan'

describe('buildConfirmPayload', () => {
  it('sends the edited brands, dimensions, keywords and user constraints', () => {
    const payload = buildConfirmPayload({
      requirement: '吸尘器，有宠物',
      budget: '2000元',
      scenario: '家庭',
      special_needs: ['宠物毛发'],
      brands: ['戴森', '石头'],
      dimensions: [{ name: '吸力', description: '吸尘效果' }],
      keywords: ['吸尘器评测'],
    })
    expect(payload.brands).toEqual(['戴森', '石头'])
    expect(payload.dimensions).toEqual([{ name: '吸力', description: '吸尘效果' }])
    expect(payload.keywords).toEqual(['吸尘器评测'])
    expect(payload.budget).toBe('2000元')
    expect(payload.scenario).toBe('家庭')
    expect(payload.special_needs).toEqual(['宠物毛发'])
  })
})

describe('list edits', () => {
  it('adds and removes brands', () => {
    expect(addListItem(['戴森'], '小米')).toEqual(['戴森', '小米'])
    expect(removeListItem(['戴森', '小米'], '戴森')).toEqual(['小米'])
  })
})
