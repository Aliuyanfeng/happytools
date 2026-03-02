/**
 * 功能模块配置
 * 新增功能入口时，只需在此文件中添加配置即可
 */

import type { Component } from 'vue'

export interface ModuleConfig {
  id: string
  name: string
  description: string
  path: string
  icon: string
  theme: 'blue' | 'green' | 'purple' | 'orange' | 'red' | 'cyan'
  tags?: string[]
}

// 功能模块列表
export const modules: ModuleConfig[] = [
  {
    id: 'dashboard',
    name: '系统仪表盘',
    description: '',
    path: '/dashboard',
    icon: 'DashboardOutlined',
    theme: 'blue',
    tags: []
  },
  {
    id: 'todo',
    name: '智能待办',
    description: '',
    path: '/todo',
    icon: 'CheckCircleOutlined',
    theme: 'green',
    tags: []
  },
  {
    id: 'toolbox',
    name: '工具盒子',
    description: '',
    path: '/toolbox',
    icon: 'ToolOutlined',
    theme: 'purple',
    tags: []
  },
  {
    id: 'virusTotal',
    name: 'VirusTotal',
    description: '',
    path: '/vt',
    icon: 'SafetyOutlined',
    theme: 'red',
    tags: []
  },
  {
    id: 'network',
    name: '网络调试',
    description: '',
    path: '/network',
    icon: 'ApartmentOutlined',
    theme: 'orange',
    tags: []
  },
  {
    id: 'dailyReport',
    name: '日报管理',
    description: '记录每日工作日报，日历视图展示，支持周月切换',
    path: '/dailyReport',
    icon: 'CalendarOutlined',
    theme: 'cyan',
    tags: ['日报', '日历', '工作记录']
  }
]

// 主题颜色配置
export const themeColors = {
  blue: {
    primary: '#0ea5e9',
    secondary: '#06b6d4',
    bg: '#e0f2fe',
    border: '#0ea5e9'
  },
  green: {
    primary: '#22c55e',
    secondary: '#16a34a',
    bg: '#dcfce7',
    border: '#22c55e'
  },
  purple: {
    primary: '#a855f7',
    secondary: '#9333ea',
    bg: '#f3e8ff',
    border: '#a855f7'
  },
  orange: {
    primary: '#f97316',
    secondary: '#ea580c',
    bg: '#ffedd5',
    border: '#f97316'
  },
  red: {
    primary: '#ef4444',
    secondary: '#dc2626',
    bg: '#fee2e2',
    border: '#ef4444'
  },
  cyan: {
    primary: '#06b6d4',
    secondary: '#0891b2',
    bg: '#cffafe',
    border: '#06b6d4'
  }
}

// 欢迎信息配置
export const welcomeInfo = {
  emoji: '👋',
  title: '欢迎使用',
  subtitle: 'HappyTools',
  description: '',
  features: [
    {
      icon: '⚡',
      title: '轻量快速',
      description: ''
    },
    {
      icon: '🔒',
      title: '安全可靠',
      description: ''
    },
    {
      icon: '🎯',
      title: '简洁高效',
      description: ''
    }
  ],
  version: 'v1.0.0'
}

// 获取模块配置
export function getModule(id: string): ModuleConfig | undefined {
  return modules.find(module => module.id === id)
}

// 添加新模块
export function addModule(module: ModuleConfig): void {
  modules.push(module)
}

// 删除模块
export function removeModule(id: string): void {
  const index = modules.findIndex(module => module.id === id)
  if (index !== -1) {
    modules.splice(index, 1)
  }
}
