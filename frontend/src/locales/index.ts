import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import enUS from './en-US'

export type AppLocale = 'zh-CN' | 'en-US'

// 获取系统语言
function getSystemLocale(): AppLocale {
  // 从 navigator 获取系统语言
  const systemLang = navigator.language || (navigator as any).userLanguage
  if (systemLang) {
    // 标准化语言代码
    const lang = systemLang.toLowerCase()
    if (lang.startsWith('zh')) {
      return 'zh-CN'
    }
    if (lang.startsWith('en')) {
      return 'en-US'
    }
  }
  return 'zh-CN' // 默认中文
}

// 从 localStorage 获取保存的语言设置
function getSavedLocale(): AppLocale | null {
  const saved = localStorage.getItem('settings')
  if (saved) {
    try {
      const settings = JSON.parse(saved)
      if (settings.language && settings.language !== 'auto') {
        return settings.language as AppLocale
      }
    } catch {
      // ignore
    }
  }
  return null
}

// 获取初始语言
function getInitialLocale(): AppLocale {
  const saved = getSavedLocale()
  if (saved) {
    return saved
  }
  return getSystemLocale()
}

// 创建 i18n 实例
const i18n = createI18n({
  legacy: false, // 使用 Composition API 模式
  locale: getInitialLocale(),
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
  },
})

// 切换语言的函数
export function setLocale(locale: string): void {
  let targetLocale: AppLocale
  if (locale === 'auto') {
    targetLocale = getSystemLocale()
  } else {
    targetLocale = locale as AppLocale
  }
  i18n.global.locale.value = targetLocale
  // 更新 HTML lang 属性
  document.documentElement.lang = targetLocale
}

// 获取当前语言
export function getLocale(): AppLocale {
  return i18n.global.locale.value as AppLocale
}

// 获取系统语言
export function getSystemLanguage(): AppLocale {
  return getSystemLocale()
}

export default i18n
