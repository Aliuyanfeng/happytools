import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { setLocale } from '../locales'

export type ThemeMode = 'light' | 'dark' | 'auto'
export type FontSize = 'small' | 'medium' | 'large'
export type Language = 'auto' | 'zh-CN' | 'en-US'
export type CloseBehavior = 'exit' | 'hide'

export const useSettingsStore = defineStore('settings', () => {
  // 主题模式
  const themeMode = ref<ThemeMode>((localStorage.getItem('themeMode') as ThemeMode) || 'light')

  // 字体大小
  const fontSize = ref<FontSize>((localStorage.getItem('fontSize') as FontSize) || 'medium')

  // 自定义字体
  const customFont = ref<string>(localStorage.getItem('customFont') || '')

  // 界面语言
  const language = ref<Language>((localStorage.getItem('language') as Language) || 'auto')

  // 关闭按钮行为
  const closeBehavior = ref<CloseBehavior>((localStorage.getItem('closeBehavior') as CloseBehavior) || 'exit')

  // VirusTotal API Key
  const vtApiKey = ref<string>(localStorage.getItem('vtApiKey') || '')

  // VirusTotal 并发扫描数
  const vtConcurrency = ref<number>(parseInt(localStorage.getItem('vtConcurrency') || '5'))

  // 首页模块可见性（存储隐藏的模块 id 列表）
  const hiddenModules = ref<string[]>(
    JSON.parse(localStorage.getItem('hiddenModules') || '[]')
  )

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

  watch(language, (newVal) => {
    localStorage.setItem('language', newVal)
    setLocale(newVal)
  })

  watch(closeBehavior, (newVal) => {
    localStorage.setItem('closeBehavior', newVal)
  })

  watch(vtApiKey, (newVal) => {
    localStorage.setItem('vtApiKey', newVal)
  })

  watch(vtConcurrency, (newVal) => {
    localStorage.setItem('vtConcurrency', String(newVal))
  })

  watch(hiddenModules, (newVal) => {
    localStorage.setItem('hiddenModules', JSON.stringify(newVal))
  }, { deep: true })

  // 当前实际是否为暗色（用于 ConfigProvider）
  const isDark = ref(false)

  // 监听系统主题变化（用于 auto 模式）
  const systemDarkMQ = window.matchMedia('(prefers-color-scheme: dark)')
  systemDarkMQ.addEventListener('change', () => {
    if (themeMode.value === 'auto') {
      applyTheme('auto')
    }
  })

  // 应用主题
  function applyTheme(mode: ThemeMode) {
    const root = document.documentElement
    let dark = false

    if (mode === 'auto') {
      dark = window.matchMedia('(prefers-color-scheme: dark)').matches
    } else {
      dark = mode === 'dark'
    }

    isDark.value = dark
    root.setAttribute('data-theme', dark ? 'dark' : 'light')
  }

  // 应用字体大小
  function applyFontSize(size: FontSize) {
    const sizeMap = { small: '13px', medium: '14px', large: '16px' }
    const px = sizeMap[size]
    document.documentElement.style.fontSize = px
    document.body.style.fontSize = px
  }

  // 应用自定义字体
  function applyCustomFont(font: string) {
    const family = font
      ? `${font}, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif`
      : '-apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif'
    document.documentElement.style.fontFamily = family
    document.body.style.fontFamily = family
  }

  // 初始化应用设置
  function initSettings() {
    applyTheme(themeMode.value)
    applyFontSize(fontSize.value)
    applyCustomFont(customFont.value)
    setLocale(language.value)
  }

  return {
    themeMode,
    fontSize,
    customFont,
    language,
    closeBehavior,
    vtApiKey,
    vtConcurrency,
    hiddenModules,
    isDark,
    initSettings
  }
})
