import { test, expect } from '@playwright/test'

const parseBody = {
  understanding: '你想买一台适合有宠物家庭的吸尘器',
  product_type: '无线吸尘器',
  budget: '2000元左右',
  scenario: '家庭使用',
  special_needs: ['宠物毛发清理'],
  brands: ['戴森', '小米'],
  dimensions: [{ name: '吸力', description: '吸尘效果' }],
  keywords: ['吸尘器评测'],
}

const reportBody = {
  id: 99,
  history_id: 1,
  category: '无线吸尘器',
  data: {
    category: '无线吸尘器',
    brands: ['戴森', '小米'],
    dimensions: [{ name: '吸力', description: '吸尘效果' }],
    scores: { 戴森: { 吸力: 8.5 }, 小米: { 吸力: 7.2 } },
    rankings: [
      { brand: '戴森', overall_score: 8.5, rank: 1, scores: { 吸力: 8.5 } },
    ],
    recommendation: '戴森吸力更稳，适合有宠物的家庭。',
    stats: { total_videos: 3, total_comments: 40, comments_by_brand: { 戴森: 20 } },
  },
  created_at: '2026-08-13',
}

test.describe('商品主路径', () => {
  test('需求 → 确认 → 进度完成 → 报告', async ({ page }) => {
    await page.route('**/api/parse', async (route) => {
      await route.fulfill({ json: parseBody })
    })
    await page.route('**/api/confirm', async (route) => {
      await route.fulfill({ json: { task_id: 'task-e2e', message: 'ok' } })
    })
    await page.route('**/api/history/task-e2e', async (route) => {
      await route.fulfill({
        json: {
          id: 1,
          taskId: 'task-e2e',
          status: 'completed',
          reportId: 99,
          progressMsg: '分析完成',
        },
      })
    })
    await page.route('**/api/report/99', async (route) => {
      await route.fulfill({ json: reportBody })
    })
    await page.route('**/api/history/1', async (route) => {
      await route.fulfill({ json: { brands: ['戴森', '小米'] } })
    })

    await page.goto('/')
    await page.getByPlaceholder(/描述你的需求/).fill('吸尘器，有宠物')
    await page.locator('form').locator('button[type="submit"]').click()

    await expect(page.getByText('确认分析方案')).toBeVisible()
    await expect(page.getByText('戴森')).toBeVisible()

    await page.getByRole('button', { name: /确认开始分析/ }).click()

    await expect(page).toHaveURL(/\/report\/99/)
    await expect(page.getByRole('heading', { name: '戴森' }).first()).toBeVisible()
    await expect(page.getByText('推荐', { exact: true })).toBeVisible()
  })
})
