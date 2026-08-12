import type { Dimension, ReportData } from '../types/report'

export const RADAR_SCORE_MAX = 10

export function filterReportByDimensions(data: ReportData, selected: string[]): ReportData {
  const allow = new Set(selected)
  const dimensions: Dimension[] = selected.length
    ? data.dimensions.filter((d) => allow.has(d.name))
    : []
  return {
    ...data,
    dimensions,
  }
}

export function overviewRecommendation(data: ReportData): { brand: string; reason: string } | null {
  const top = data.rankings?.[0]
  if (!top) return null
  const fromAI = firstSentence(data.recommendation)
  const reason = fromAI || `${top.brand}综合得分 ${top.overall_score.toFixed(1)}，可作为首选参考`
  return { brand: top.brand, reason }
}

function firstSentence(text?: string): string {
  if (!text) return ''
  const cleaned = text.replace(/^#+\s+/gm, '').trim()
  const match = cleaned.match(/[^。！？\n]+[。！？]?/)
  return match ? match[0].trim() : cleaned.slice(0, 80)
}
