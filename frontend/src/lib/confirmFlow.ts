export type ConfirmNav =
  | { kind: 'progress'; href: string }
  | { kind: 'error'; message: string }

export function confirmNavigationTarget(
  responseOk: boolean,
  payload: { task_id?: unknown; error?: unknown },
  title: string,
  mode: 'product' | 'video' = 'product',
): ConfirmNav {
  if (!responseOk) {
    const message = typeof payload.error === 'string' && payload.error ? payload.error : '创建任务失败'
    return { kind: 'error', message }
  }
  const taskId = typeof payload.task_id === 'string' ? payload.task_id.trim() : ''
  if (!taskId) {
    return { kind: 'error', message: '服务器未返回任务ID' }
  }
  const params = new URLSearchParams()
  if (title) {
    params.set('title', title)
  }
  if (mode === 'video') {
    params.set('mode', 'video')
  }
  const qs = params.toString()
  return { kind: 'progress', href: qs ? `/progress/${taskId}?${qs}` : `/progress/${taskId}` }
}

export type ParseOutcome =
  | { kind: 'ok'; data: Record<string, unknown> }
  | { kind: 'error'; message: string }

export function parseOutcomeFromResponse(
  responseOk: boolean,
  payload: Record<string, unknown> | null,
): ParseOutcome {
  if (!responseOk || !payload) {
    const message =
      payload && typeof payload.error === 'string' && payload.error
        ? payload.error
        : '解析需求失败，请检查设置中的 AI 配置'
    return { kind: 'error', message }
  }
  if (typeof payload.error === 'string' && payload.error) {
    return { kind: 'error', message: payload.error }
  }
  if (!Array.isArray(payload.brands) || payload.brands.length === 0) {
    return { kind: 'error', message: '解析结果缺少品牌，请重试或先配置 AI' }
  }
  return { kind: 'ok', data: payload }
}

export function beginSubmit(alreadySubmitting: boolean): boolean {
  return !alreadySubmitting
}
