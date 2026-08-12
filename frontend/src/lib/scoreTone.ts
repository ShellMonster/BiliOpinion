export type ScoreTone = 'success' | 'info' | 'warning' | 'danger'

export function scoreTone(score: number): ScoreTone {
  if (score >= 8) return 'success'
  if (score >= 6) return 'info'
  if (score >= 4) return 'warning'
  return 'danger'
}

export function scoreToneTextClass(score: number): string {
  switch (scoreTone(score)) {
    case 'success':
      return 'text-emerald-500'
    case 'info':
      return 'text-blue-500'
    case 'warning':
      return 'text-amber-500'
    default:
      return 'text-rose-500'
  }
}

export function scoreToneBadgeClass(score: number): string {
  switch (scoreTone(score)) {
    case 'success':
      return 'text-emerald-600 bg-emerald-50 border-emerald-100'
    case 'info':
      return 'text-blue-600 bg-blue-50 border-blue-100'
    case 'warning':
      return 'text-amber-600 bg-amber-50 border-amber-100'
    default:
      return 'text-rose-600 bg-rose-50 border-rose-100'
  }
}

export function scoreToneBarClass(score: number): string {
  switch (scoreTone(score)) {
    case 'success':
      return 'bg-emerald-500'
    case 'info':
      return 'bg-blue-500'
    case 'warning':
      return 'bg-amber-500'
    default:
      return 'bg-rose-500'
  }
}
