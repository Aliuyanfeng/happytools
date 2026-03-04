import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'auto'
export type FontSize = 'small' | 'medium' | 'large'

export const useSettingsStore = defineStore('settings', () => {
  // 主题模式
  const themeMode = ref<ThemeMode>((localStorage.getItem('themeMode') as ThemeMode) || 'light')

  // 字体大小
  const fontSize = ref<FontSize>((localStorage.getItem('fontSize') as FontSize) || 'medium')

  // 自定义字体
  const customFont = ref<string>(localStorage.getItem('customFont') || '')

  // VirusTotal API Key
  const vtApiKey = ref<string>(localStorage.getItem('vtApiKey') || '')

  // VirusTotal 并发扫描数
  const vtConcurrency = ref<number>(parseInt(localStorage.getItem('vtConcurrency') || '5'))

  // 监听变化并保存到本地存储
  watch(themeMode, (newVal) => {
    localStorage.setItem('themeMode', newVal)
    applyTheme(newVal)
  })

  watch(fontSize, (newVal) => {
    localStorage.setItem('fontSize', newVal)
    applyFontSize(newVal)
  })

  watch(customFont, (newVal) => {
    localStorage.setItem('customFont', newVal)
    applyCustomFont(newVal)
  })

  watch(vtApiKey, (newVal) => {
    localStorage.setItem('vtApiKey', newVal)
  })

  watch(vtConcurrency, (newVal) => {
    localStorage.setItem('vtConcurrency', String(newVal))
  })

  // 应用主题
  function applyTheme(mode: ThemeMode) {
    const root = document.documentElement

    if (mode === 'auto') {
      // 跟随系统
      const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      root.setAttribute('data-theme', isDark ? 'dark' : 'light')
    } else {
      root.setAttribute('data-theme', mode)
    }
  }

  // 应用字体大小
  function applyFontSize(size: FontSize) {
    const root = document.documentElement
    const sizeMap = {
      small: '14px',
      medium: '16px',
      large: '18px'
    }
    root.style.fontSize = sizeMap[size]
  }

  // 应用自定义字体
  function applyCustomFont(font: string) {
    const root = document.documentElement
    if (font) {
      root.style.fontFamily = font
    } else {
      root.style.fontFamily = ''
    }
  }

  // 初始化应用设置
  function initSettings() {
    applyTheme(themeMode.value)
    applyFontSize(fontSize.value)
    applyCustomFont(customFont.value)
  }

  return {
    themeMode,
    fontSize,
    customFont,
    vtApiKey,
    vtConcurrency,
    initSettings
  }
})
