import { useState, useEffect } from 'react'
import Modal from '../common/Modal'
import Input from '../common/Input'
import Button from '../common/Button'
import { useToast } from '../../hooks/useToast'

interface SettingsModalProps {
  isOpen: boolean
  onClose: () => void
}

interface SettingsData {
  aiApiBase: string
  aiApiKey: string
  aiModel: string
  bilibiliCookie: string
}

export default function SettingsModal({ isOpen, onClose }: SettingsModalProps) {
  const { showToast } = useToast()
  const [settings, setSettings] = useState<SettingsData>({
    aiApiBase: 'https://api.openai.com/v1',
    aiApiKey: '',
    aiModel: 'gemini-3-flash-preview',
    bilibiliCookie: ''
  })
  const [scrapeMaxConcurrency, setScrapeMaxConcurrency] = useState(5)
  const [aiMaxConcurrency, setAiMaxConcurrency] = useState(10)

  // Load settings from backend API when modal opens
  useEffect(() => {
    if (isOpen) {
      fetch('http://localhost:8080/api/config')
        .then(res => res.json())
        .then(data => {
          setSettings({
            aiApiBase: data.ai_base_url || 'https://api.openai.com/v1',
            aiApiKey: data.ai_api_key || '',
            aiModel: data.ai_model || 'gemini-3-flash-preview',
            bilibiliCookie: data.bilibili_cookie || ''
          })
          setScrapeMaxConcurrency(parseInt(data.scrape_max_concurrency) || 5)
          setAiMaxConcurrency(parseInt(data.ai_max_concurrency) || 10)
        })
        .catch(err => {
          console.error('加载配置失败:', err)
          const saved = localStorage.getItem('settings')
          if (saved) {
            setSettings(JSON.parse(saved))
          }
        })
    }
  }, [isOpen])

  const handleSave = async () => {
    localStorage.setItem('settings', JSON.stringify(settings))
    
    try {
      await fetch('http://localhost:8080/api/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          ai_base_url: settings.aiApiBase,
          ai_api_key: settings.aiApiKey,
          ai_model: settings.aiModel,
          bilibili_cookie: settings.bilibiliCookie,
          scrape_max_concurrency: String(scrapeMaxConcurrency),
          ai_max_concurrency: String(aiMaxConcurrency)
        })
      })
    } catch (error) {
      console.error('Failed to save to backend:', error)
    }
    
    showToast('设置已保存', 'success')
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="系统设置">
      {/* AI配置 */}
      <div className="space-y-4 mb-8">
        <h3 className="text-lg font-bold text-slate-700">AI配置</h3>
        
        <Input
          label="API Base URL"
          value={settings.aiApiBase}
          onChange={(e) => setSettings({...settings, aiApiBase: e.target.value})}
          placeholder="https://api.openai.com/v1"
        />

        <Input
          label="API Key"
          type="password"
          value={settings.aiApiKey}
          onChange={(e) => setSettings({...settings, aiApiKey: e.target.value})}
          placeholder="sk-..."
        />

        <Input
          label="Model"
          value={settings.aiModel}
          onChange={(e) => setSettings({...settings, aiModel: e.target.value})}
          placeholder="gemini-3-flash-preview"
        />
      </div>

      {/* B站Cookie */}
      <div className="space-y-4 mb-8">
        <h3 className="text-lg font-bold text-slate-700">B站Cookie</h3>
        
        <div>
          <label className="block text-sm font-bold text-slate-700 mb-2">
            完整Cookie字符串
          </label>
          <textarea
            className="w-full h-32 rounded-2xl bg-slate-100 px-4 py-3 text-sm text-slate-900 
                       placeholder:text-slate-400 focus:bg-white focus:ring-2 focus:ring-blue-500/20 
                       transition-all duration-200 outline-none resize-none"
            value={settings.bilibiliCookie}
            onChange={(e) => setSettings({...settings, bilibiliCookie: e.target.value})}
            placeholder="SESSDATA=xxx; buvid3=xxx; ..."
          />
          <p className="text-xs text-slate-500 mt-2">
            从浏览器开发者工具中复制完整的Cookie字符串
          </p>
        </div>
      </div>

      {/* 并发配置 */}
      <div className="space-y-4 mb-8">
        <h3 className="text-lg font-bold text-slate-700">并发配置</h3>
        
        <div>
          <label className="block text-sm font-bold text-slate-700 mb-2">
            抓取并发数
          </label>
          <input
            type="number"
            min={1}
            max={10}
            value={scrapeMaxConcurrency}
            onChange={(e) => setScrapeMaxConcurrency(Math.min(10, Math.max(1, parseInt(e.target.value) || 1)))}
            className="w-full rounded-2xl bg-slate-100 px-4 py-3 text-sm text-slate-900 
                       placeholder:text-slate-400 focus:bg-white focus:ring-2 focus:ring-blue-500/20 
                       transition-all duration-200 outline-none"
          />
          <p className="text-xs text-amber-600 mt-2">
            ⚠️ 并发数过高可能触发B站反爬机制，建议保持默认值5
          </p>
        </div>

        <div>
          <label className="block text-sm font-bold text-slate-700 mb-2">
            AI并发数
          </label>
          <input
            type="number"
            min={1}
            max={20}
            value={aiMaxConcurrency}
            onChange={(e) => setAiMaxConcurrency(Math.min(20, Math.max(1, parseInt(e.target.value) || 1)))}
            className="w-full rounded-2xl bg-slate-100 px-4 py-3 text-sm text-slate-900 
                       placeholder:text-slate-400 focus:bg-white focus:ring-2 focus:ring-blue-500/20 
                       transition-all duration-200 outline-none"
          />
          <p className="text-xs text-amber-600 mt-2">
            ⚠️ 并发数过高可能触发API频率限制，建议根据API配额调整
          </p>
        </div>
      </div>

      {/* 保存按钮 */}
      <Button onClick={handleSave} className="w-full">
        保存设置
      </Button>
    </Modal>
  )
}
