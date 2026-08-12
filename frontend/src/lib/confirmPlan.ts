export type ConfirmDimension = { name: string; description: string }

export type ConfirmPlan = {
  requirement: string
  budget?: string
  scenario?: string
  special_needs?: string[]
  brands: string[]
  dimensions: ConfirmDimension[]
  keywords: string[]
}

export function buildConfirmPayload(
  plan: ConfirmPlan,
  extras: Record<string, unknown> = {},
): Record<string, unknown> {
  return {
    requirement: plan.requirement,
    budget: plan.budget || '',
    scenario: plan.scenario || '',
    special_needs: plan.special_needs || [],
    brands: plan.brands,
    dimensions: plan.dimensions,
    keywords: plan.keywords,
    ...extras,
  }
}

export function addListItem(list: string[], item: string): string[] {
  const value = item.trim()
  if (!value || list.includes(value)) return list
  return [...list, value]
}

export function removeListItem(list: string[], item: string): string[] {
  return list.filter((entry) => entry !== item)
}
