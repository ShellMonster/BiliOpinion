import { test, expect } from '@playwright/test'

test.describe('Settings Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/settings')
  })

  test('should display all configuration sections', async ({ page }) => {
    await expect(page.getByRole('heading', { name: '系统设置' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'AI配置' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'B站Cookie' })).toBeVisible()
    await expect(page.getByRole('heading', { name: '并发配置' })).toBeVisible()
  })

  test('should display concurrency input fields with default values', async ({ page }) => {
    const scrapeInput = page.getByLabel('抓取并发数')
    const aiInput = page.getByLabel('AI并发数')

    await expect(scrapeInput).toBeVisible()
    await expect(aiInput).toBeVisible()

    await expect(scrapeInput).toHaveValue('5')
    await expect(aiInput).toHaveValue('10')
  })

  test('should display warning messages for concurrency inputs', async ({ page }) => {
    await expect(page.getByText('⚠️ 并发数过高可能触发B站反爬机制，建议保持默认值5')).toBeVisible()
    await expect(page.getByText('⚠️ 并发数过高可能触发API频率限制，建议根据API配额调整')).toBeVisible()
  })

  test('should enforce min/max limits on scrape concurrency input', async ({ page }) => {
    const scrapeInput = page.getByLabel('抓取并发数')

    await scrapeInput.fill('0')
    await scrapeInput.blur()
    await expect(scrapeInput).toHaveValue('1')

    await scrapeInput.fill('15')
    await scrapeInput.blur()
    await expect(scrapeInput).toHaveValue('10')

    await scrapeInput.fill('7')
    await scrapeInput.blur()
    await expect(scrapeInput).toHaveValue('7')
  })

  test('should enforce min/max limits on AI concurrency input', async ({ page }) => {
    const aiInput = page.getByLabel('AI并发数')

    await aiInput.fill('0')
    await aiInput.blur()
    await expect(aiInput).toHaveValue('1')

    await aiInput.fill('25')
    await aiInput.blur()
    await expect(aiInput).toHaveValue('20')

    await aiInput.fill('15')
    await aiInput.blur()
    await expect(aiInput).toHaveValue('15')
  })

  test('should load configuration from backend API on mount', async ({ page }) => {
    await page.route('http://localhost:8080/api/config', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          ai_base_url: 'https://custom-api.com/v1',
          ai_api_key: 'test-key-123',
          ai_model: 'gpt-4',
          bilibili_cookie: 'test-cookie',
          scrape_max_concurrency: '8',
          ai_max_concurrency: '15'
        })
      })
    })

    await page.reload()

    await expect(page.getByLabel('API Base URL')).toHaveValue('https://custom-api.com/v1')
    await expect(page.getByLabel('Model')).toHaveValue('gpt-4')
    await expect(page.getByLabel('抓取并发数')).toHaveValue('8')
    await expect(page.getByLabel('AI并发数')).toHaveValue('15')
  })

  test('should save configuration to backend API when clicking save button', async ({ page }) => {
    let savedConfig: any = null

    await page.route('http://localhost:8080/api/config', async route => {
      if (route.request().method() === 'POST') {
        savedConfig = await route.request().postDataJSON()
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ success: true })
        })
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            ai_base_url: 'https://api.openai.com/v1',
            ai_api_key: '',
            ai_model: 'gemini-3-flash-preview',
            bilibili_cookie: '',
            scrape_max_concurrency: '5',
            ai_max_concurrency: '10'
          })
        })
      }
    })

    await page.reload()

    await page.getByLabel('抓取并发数').fill('7')
    await page.getByLabel('AI并发数').fill('12')
    await page.getByRole('button', { name: '保存设置' }).click()

    await page.waitForTimeout(500)

    expect(savedConfig).toBeTruthy()
    expect(savedConfig.scrape_max_concurrency).toBe('7')
    expect(savedConfig.ai_max_concurrency).toBe('12')
  })

  test('should show success toast after saving', async ({ page }) => {
    await page.route('http://localhost:8080/api/config', async route => {
      if (route.request().method() === 'POST') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ success: true })
        })
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            ai_base_url: 'https://api.openai.com/v1',
            ai_api_key: '',
            ai_model: 'gemini-3-flash-preview',
            bilibili_cookie: '',
            scrape_max_concurrency: '5',
            ai_max_concurrency: '10'
          })
        })
      }
    })

    await page.reload()
    await page.getByRole('button', { name: '保存设置' }).click()

    await expect(page.getByText('设置已保存')).toBeVisible()
  })

  test('should show error toast when save fails', async ({ page }) => {
    await page.route('http://localhost:8080/api/config', async route => {
      if (route.request().method() === 'POST') {
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal server error' })
        })
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            ai_base_url: 'https://api.openai.com/v1',
            ai_api_key: '',
            ai_model: 'gemini-3-flash-preview',
            bilibili_cookie: '',
            scrape_max_concurrency: '5',
            ai_max_concurrency: '10'
          })
        })
      }
    })

    await page.reload()
    await page.getByRole('button', { name: '保存设置' }).click()

    await expect(page.getByText('保存失败')).toBeVisible()
  })
})
