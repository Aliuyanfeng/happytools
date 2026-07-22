/**
 * 功能模块配置
 * 新增功能入口时，只需在此文件中添加配置即可
 */

import type { Component } from 'vue'
import { useI18n } from 'vue-i18n'

export type ModuleCategory = 'security' | 'dev' | 'efficiency' | 'network' | 'doc'

export interface ModuleConfig {
  id: string
  nameKey: string // i18n key
  description: string
  path: string
  icon: string
  theme: 'blue' | 'green' | 'purple' | 'orange' | 'red' | 'cyan'
  category: ModuleCategory
  badge?: number // 消息角标数
}

// 功能模块列表
export const modules: ModuleConfig[] = [
  {
    id: 'dashboard',
    nameKey: 'home.modules.dashboard',
    description: '实时监控系统运行状态与资源占用',
    path: '/dashboard',
    icon: 'DashboardOutlined',
    theme: 'blue',
    category: 'dev',
  },
  {
    id: 'todo',
    nameKey: 'home.modules.todo',
    description: '高效管理待办事项与任务提醒',
    path: '/todo',
    icon: 'CheckCircleOutlined',
    theme: 'green',
    category: 'efficiency',
  },
  {
    id: 'toolbox',
    nameKey: 'home.modules.toolbox',
    description: '常用开发工具集合，一键调用',
    path: '/toolbox',
    icon: 'ToolOutlined',
    theme: 'purple',
    category: 'dev',
  },
  {
    id: 'virusTotal',
    nameKey: 'home.modules.virusTotal',
    description: '在线病毒检测与文件安全分析',
    path: '/vt',
    icon: 'SafetyOutlined',
    theme: 'red',
    category: 'security',
  },
  {
    id: 'network',
    nameKey: 'home.modules.network',
    description: '网络连接调试与诊断工具',
    path: '/network',
    icon: 'ApartmentOutlined',
    theme: 'orange',
    category: 'network',
  },
  {
    id: 'dailyReport',
    nameKey: 'home.modules.dailyReport',
    description: '记录每日工作日报，日历视图展示',
    path: '/dailyReport',
    icon: 'CalendarOutlined',
    theme: 'cyan',
    category: 'efficiency',
  },
  {
    id: 'gitConfig',
    nameKey: 'home.modules.gitConfig',
    description: 'Git 用户配置与多账号管理',
    path: '/git-config',
    icon: 'BranchesOutlined',
    theme: 'green',
    category: 'dev',
  },
  {
    id: 'makefileEditor',
    nameKey: 'home.modules.makefileEditor',
    description: 'Makefile 语法高亮编辑与模板管理',
    path: '/makefile-editor',
    icon: 'FileTextOutlined',
    theme: 'orange',
    category: 'dev',
  },
  {
    id: 'nucleiParser',
    nameKey: 'home.modules.nucleiParser',
    description: 'Nuclei POC 模板可视化解析',
    path: '/nuclei-parser',
    icon: 'BugOutlined',
    theme: 'red',
    category: 'security',
  },
  {
    id: 'pdfParser',
    nameKey: 'home.modules.pdfParser',
    description: 'PDF 文件解析与元数据管理',
    path: '/pdf-parser',
    icon: 'FilePdfOutlined',
    theme: 'blue',
    category: 'doc',
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
