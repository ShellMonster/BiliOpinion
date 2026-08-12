import type { ReportData } from '../types/report'

export const RADAR_SCORE_MAX = 10

export function filterReportByDimensions(data: ReportData, selected: string[]): ReportData {
  const allow = new Set(selected)
  return {
    ...data,
    dimensions: data.dimensions.filter((d) => allow.has(d.name)),
  }
}

export function overviewRecommendation(
  data: ReportData,
  options: { hideUnknown?: boolean } = {},
): { brand: string; reason: string } | null {
  const top = (data.rankings || []).find((r) => {
    if (options.hideUnknown && r.brand === '未知') return false
    if (r.overall_score === 0) return false
    return true
  })
  if (!top) return null
  const fromAI = firstSentence(data.recommendation)
  const reason = fromAI || `${top.brand}综合得分 ${top.overall_score.toFixed(1)}，可作为首选参考`
  return { brand: top.brand, reason }
}

function firstSentence(text?: string): string {
  if (!text) return ''
  const lines = text
    .split('\n')
    .map((line) => line.replace(/^#+\s*/, '').trim())
    .filter((line) => line.length > 0 && !/^(综合推荐|购买建议|总结|推荐)$/.test(line))
  const body = lines.join(' ')
  const match = body.match(/[^。！？]+[。！？]/)
  if (match) return match[0].trim()
  return body.slice(0, 80).trim()
}
