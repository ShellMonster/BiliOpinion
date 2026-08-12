import { describe, expect, it } from 'vitest'
import { filterReportByDimensions, overviewRecommendation, RADAR_SCORE_MAX } from './reportView'
import type { ReportData } from '../types/report'

const sample: ReportData = {
  category: '吸尘器',
  brands: ['戴森', '小米'],
  dimensions: [
    { name: '吸力', description: '' },
    { name: '续航', description: '' },
  ],
  scores: { 戴森: { 吸力: 8.5, 续航: 7 }, 小米: { 吸力: 7, 续航: 8 } },
  rankings: [
    { brand: '戴森', overall_score: 8.5, rank: 1, scores: { 吸力: 8.5, 续航: 7 } },
  ],
  recommendation: '戴森吸力更稳，适合有宠物的家庭。小米更便宜。',
}

describe('RADAR_SCORE_MAX', () => {
  it('is 10 on the 1-10 scale', () => {
    expect(RADAR_SCORE_MAX).toBe(10)
  })
})

describe('filterReportByDimensions', () => {
  it('keeps only selected dimensions for chart inputs', () => {
    const filtered = filterReportByDimensions(sample, ['吸力'])
    expect(filtered.dimensions.map((d) => d.name)).toEqual(['吸力'])
    expect(filtered.dimensions.map((d) => d.name)).not.toContain('续航')
  })
})

describe('overviewRecommendation', () => {
  it('puts the top brand and a short reason first', () => {
    const rec = overviewRecommendation(sample)
    expect(rec?.brand).toBe('戴森')
    expect(rec?.reason).toContain('戴森')
  })
})
