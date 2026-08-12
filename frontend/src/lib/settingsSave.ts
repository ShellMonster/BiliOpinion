export function settingsSaveOutcome(responseOk: boolean): { kind: 'ok' } | { kind: 'error'; message: string } {
  if (!responseOk) {
    return { kind: 'error', message: '保存失败' }
  }
  return { kind: 'ok' }
}
