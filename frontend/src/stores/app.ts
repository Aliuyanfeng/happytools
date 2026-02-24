/*
 * @Author: LiuYanFeng
 * @Date: 2026-02-06
 * @Description: 应用状态管理
 */
import { defineStore } from 'pinia'
import { ref, onMounted } from 'vue'
import { AppSettingsService } from '../../bindings/github.com/Aliuyanfeng/happytools/backend/services/appsettings'
export const useAppStore = defineStore('app', () => {
  const currentTime = new Date().toLocaleString()
  var lastUsedTime = ref<string>(currentTime)

  function updateLastUsedTime(timestr) {
    lastUsedTime.value = timestr
  }
  return {
    lastUsedTime,
    updateLastUsedTime
  }
})
