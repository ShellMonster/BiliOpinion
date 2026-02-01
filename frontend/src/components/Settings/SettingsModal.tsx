import { useState, useEffect } from 'react'
import Modal from '../common/Modal'
import Input from '../common/Input'
import Button from '../common/Button'

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
  const [settings, setSettings] = useState<SettingsData>({
    aiApiBase: 'https://api.openai.com/v1',
    aiApiKey: '',
    aiModel: 'gpt-3.5-turbo',
    bilibiliCookie: ''
  })

  // Load settings from localStorage when modal opens
  useEffect(() => {
    if (isOpen) {
      const saved = localStorage.getItem('settings')
      if (saved) {
        setSettings(JSON.parse(saved))
      }
    }
  }, [isOpen])

  const handleSave = async () => {
    // Save to localStorage
    localStorage.setItem('settings', JSON.stringify(settings))
    
    // Also save to backend API
    try {
      await fetch('http://localhost:8080/api/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          ai_base_url: settings.aiApiBase,
          ai_api_key: settings.aiApiKey,
          ai_model: settings.aiModel,
          bilibili_cookie: settings.bilibiliCookie
        })
      })
    } catch (error) {
      console.error('Failed to save to backend:', error)
    }
    
    alert('设置已保存')
    onClose()
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
          placeholder="gpt-3.5-turbo"
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

      {/* 保存按钮 */}
      <Button onClick={handleSave} className="w-full">
        保存设置
      </Button>
    </Modal>
  )
}
