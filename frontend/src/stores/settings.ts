import { defineStore } from 'pinia';

interface SettingsState {
  fontFamily: string;
  fontSize: number;
  cacheDir: string;
}

export const useSettingsStore = defineStore('settings', {
  state: (): SettingsState => {
    // 从localStorage加载设置，如果没有则使用默认值
    const savedSettings = localStorage.getItem('happyToolsSettings');
    if (savedSettings) {
      return JSON.parse(savedSettings);
    }

    return {
      fontFamily: 'system',
      fontSize: 14,
      cacheDir: '',
    };
  },
  actions: {
    saveSettings(newSettings: Partial<SettingsState>) {
      // 更新状态
      this.$patch(newSettings);
      // 保存到localStorage
      localStorage.setItem('happyToolsSettings', JSON.stringify(this.$state));
    },
  },
});