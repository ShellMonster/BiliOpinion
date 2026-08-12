import { test, expect } from '@playwright/test'

test.describe('设置弹窗', () => {
  test('顶栏打开设置而不是 /settings 死路由', async ({ page }) => {
    await page.route('**/api/config', async (route) => {
      if (route.request().method() === 'GET') {
        await route.fulfill({
          json: {
            ai_base_url: 'https://api.openai.com/v1',
            ai_api_key: '',
            ai_model: 'gemini-3-flash-preview',
            bilibili_cookie: '',
            scrape_max_concurrency: '5',
            ai_max_concurrency: '10',
          },
        })
        return
      }
      await route.fulfill({ json: { message: 'ok' } })
    })

    await page.goto('/')
    await page.getByRole('button', { name: '设置' }).click()
    await expect(page.getByText('系统设置')).toBeVisible()
    await expect(page.getByText('AI配置')).toBeVisible()
  })
})
