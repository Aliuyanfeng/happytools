<template>
  <div class="db">

    <!-- ══ 左侧：核心指标竖栏 ══ -->
    <aside class="db-aside">

      <!-- 主机名 -->
      <div class="aside-host">
        <div class="host-dot" :class="{ online: hardwareInfo.hostname }" />
        <div>
          <div class="host-name">{{ hardwareInfo.hostname || 'Loading…' }}</div>
          <div class="host-os">{{ hardwareInfo.platform || '—' }}</div>
        </div>
      </div>

      <div class="aside-divider" />

      <!-- 内存环形 -->
      <div class="metric-block">
        <div class="metric-header">
          <span class="metric-label">内存</span>
          <span class="metric-val">{{ memoryUsed }}</span>
        </div>
        <div class="ring-wrap">
          <a-progress
            type="circle"
            :percent="memoryUsage"
            :stroke-color="ringColor(memoryUsage)"
            :stroke-width="8"
            :width="52"
            :format="p => p + '%'"
          />
        </div>
        <div class="metric-sub">共 {{ memoryTotal }}</div>
      </div>

      <div class="aside-divider" />

      <!-- 硬盘环形 -->
      <div class="metric-block">
        <div class="metric-header">
          <span class="metric-label">硬盘</span>
          <span class="metric-val">{{ formatBytes(diskUsage.used) }}</span>
        </div>
        <div class="ring-wrap">
          <a-progress
            type="circle"
            :percent="diskUsage.used_percent"
            :stroke-color="ringColor(diskUsage.used_percent)"
            :stroke-width="8"
            :width="52"
            :format="p => p + '%'"
          />
        </div>
        <div class="metric-sub">共 {{ formatBytes(diskUsage.total) }}</div>
      </div>

      <div class="aside-divider" />

      <!-- 网卡 -->
      <div class="metric-block net-block">
        <div class="metric-header">
          <span class="metric-label">网卡</span>
          <div class="net-ctrl">
            <a-button type="text" size="small" @click="isCurrentFavorite() ? handleUnfavorite() : handleFavorite()" :class="{ 'fav-on': isCurrentFavorite() }">
              {{ isCurrentFavorite() ? '⭐' : '☆' }}
            </a-button>
            <a-select v-model:value="currentNetworkInterface" size="small" style="width:90px" @change="handleChange">
              <a-select-option v-for="n in networkInterfaces" :key="n.name" :value="n.name">{{ n.name }}</a-select-option>
            </a-select>
          </div>
        </div>
        <template v-for="n in networkInterfaces" :key="n.name">
          <div v-if="n.name === currentNetworkInterface" class="net-detail">
            <!-- 网卡名称 + 状态 -->
            <div class="nd-row">
              <span class="nd-k">名称</span>
              <span class="nd-v mono">{{ n.name }}</span>
              <span class="nd-tag" :class="n.isUp ? 'up' : 'down'">{{ n.isUp ? '在线' : '离线' }}</span>
            </div>
            <!-- IP -->
            <div class="nd-row">
              <span class="nd-k">IP</span>
              <span class="nd-v">{{ n.ipv4 }}</span>
            </div>
            <!-- MAC -->
            <div class="nd-row">
              <span class="nd-k">MAC</span>
              <span class="nd-v mono">{{ n.mac }}</span>
            </div>

            <!-- 实时速率 -->
            <div class="nd-section-label">实时速率</div>
            <div class="nd-traffic">
              <div class="traf-item rx">
                <span class="traf-arrow">↓</span>
                <div>
                  <div class="traf-num">{{ formatBytes(n.bytesRecvRate) }}/s</div>
                  <div class="traf-label">接收速率</div>
                </div>
              </div>
              <div class="traf-item tx">
                <span class="traf-arrow">↑</span>
                <div>
                  <div class="traf-num">{{ formatBytes(n.bytesSentRate) }}/s</div>
                  <div class="traf-label">发送速率</div>
                </div>
              </div>
            </div>

            <!-- 累计流量 -->
            <div class="nd-section-label">累计流量</div>
            <div class="nd-traffic">
              <div class="traf-item rx">
                <span class="traf-arrow">↓</span>
                <div>
                  <div class="traf-num">{{ formatBytes(n.bytesRecv) }}</div>
                  <div class="traf-label">已接收</div>
                </div>
              </div>
              <div class="traf-item tx">
                <span class="traf-arrow">↑</span>
                <div>
                  <div class="traf-num">{{ formatBytes(n.bytesSent) }}</div>
                  <div class="traf-label">已发送</div>
                </div>
              </div>
            </div>

            <!-- 包数统计 -->
            <div class="nd-section-label">数据包</div>
            <div class="nd-traffic">
              <div class="traf-item rx">
                <span class="traf-arrow">↓</span>
                <div>
                  <div class="traf-num">{{ n.packetsRecv }}</div>
                  <div class="traf-label">接收包</div>
                </div>
              </div>
              <div class="traf-item tx">
                <span class="traf-arrow">↑</span>
                <div>
                  <div class="traf-num">{{ n.packetsSent }}</div>
                  <div class="traf-label">发送包</div>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>

    </aside>

    <!-- ══ 右侧：CPU + 主机详情 ══ -->
    <div class="db-main">

      <!-- CPU 热力图区 -->
      <div class="panel cpu-panel">
        <div class="panel-header">
          <span class="panel-title">CPU 核心负载</span>
          <span class="panel-sub">{{ hardwareInfo.modelName }}</span>
          <div class="view-toggle">
            <button class="vt-btn" :class="{ active: cpuView === 'large' }" @click="cpuView = 'large'" title="大视图">
              <span class="vt-icon">⊞</span>
            </button>
            <button class="vt-btn" :class="{ active: cpuView === 'small' }" @click="cpuView = 'small'" title="小视图">
              <span class="vt-icon">⊟</span>
            </button>
          </div>
        </div>
        <div class="cpu-heatmap" :class="cpuView === 'large' ? 'view-large' : 'view-small'">
          <div
            v-for="(pct, idx) in cpuUsage"
            :key="idx"
            class="cpu-core"
            :class="cpuView"
            :style="{ '--pct': pct, '--color': heatColor(pct) }"
          >
            <template v-if="cpuView === 'large'">
              <!-- 大视图：圆形进度 + 标签 -->
              <div class="core-ring">
                <svg viewBox="0 0 36 36" class="core-svg">
                  <circle cx="18" cy="18" r="15" fill="none" stroke="#e2e8f0" stroke-width="3" />
                  <circle
                    cx="18" cy="18" r="15" fill="none"
                    :stroke="heatColor(pct)" stroke-width="3"
                    stroke-linecap="round"
                    :stroke-dasharray="`${pct * 0.942} 94.2`"
                    stroke-dashoffset="23.55"
                    style="transition: stroke-dasharray 0.6s ease"
                  />
                </svg>
                <span class="core-ring-pct">{{ pct }}%</span>
              </div>
              <span class="core-label-lg">C{{ idx }}</span>
            </template>
            <template v-else>
              <!-- 小视图：条形 + 文字 -->
              <div class="core-bar-bg">
                <div class="core-bar-fill" />
              </div>
              <div class="core-meta">
                <span class="core-idx">C{{ idx }}</span>
                <span class="core-pct">{{ pct }}%</span>
              </div>
            </template>
          </div>
        </div>
      </div>

      <!-- 主机详情卡片 -->
      <div class="panel host-panel">
        <div class="panel-header">
          <span class="panel-title">主机详情</span>
        </div>
        <div class="host-grid">
          <div class="hg-item">
            <span class="hg-k">设备名称</span>
            <span class="hg-v">{{ hardwareInfo.hostname }}</span>
          </div>
          <div class="hg-item">
            <span class="hg-k">操作系统</span>
            <span class="hg-v">{{ hardwareInfo.platform }}</span>
          </div>
          <div class="hg-item">
            <span class="hg-k">处理器</span>
            <span class="hg-v">{{ hardwareInfo.modelName }}</span>
          </div>
          <div class="hg-item">
            <span class="hg-k">内存</span>
            <span class="hg-v">{{ memoryTotal }}</span>
          </div>
          <div class="hg-item">
            <span class="hg-k">内核版本</span>
            <span class="hg-v">{{ hardwareInfo.kernel_version }}</span>
          </div>
          <div class="hg-item">
            <span class="hg-k">系统架构</span>
            <span class="hg-v">{{ hardwareInfo.architecture }}</span>
          </div>
          <div class="hg-item">
            <span class="hg-k">存储总量</span>
            <span class="hg-v">{{ formatBytes(diskUsage.total) }}</span>
          </div>
          <div class="hg-item">
            <span class="hg-k">已用存储</span>
            <span class="hg-v">{{ formatBytes(diskUsage.used) }}</span>
          </div>
        </div>
      </div>

    </div>

    <a-modal v-model:open="settingsVisible" title="设置" @ok="saveSettings" @cancel="closeSettings">
      <settings-panel :settings="settings" @update:settings="updateSettings" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import {onMounted, onUnmounted, reactive, ref} from 'vue';
import {useRouter} from 'vue-router';
import SettingsPanel from '@/components/SettingsPanel.vue';
import {useSettingsStore} from '@/stores/settings';
import {SysInfoService} from "../../bindings/github.com/Aliuyanfeng/happytools/backend/services/monitor";
import {AppSettingsService} from "../../bindings/github.com/Aliuyanfeng/happytools/backend/services/appsettings";
import {Events} from "@wailsio/runtime";
import {message} from 'ant-design-vue';

const router = useRouter();
const settingsStore = useSettingsStore();
// 状态变量
const cpuUsage = ref<number[] | undefined>([]);
const memoryUsage = ref<number | undefined>(0); // 已使用内存百分比
const memoryTotal = ref<string | undefined>(""); // 总内存
const memoryUsed = ref<string | undefined>(""); // 已使用内存
const diskUsage = reactive({ // 硬盘使用情况
  total: 0,
  used: 0,
  used_percent: 0,
})
const hardwareInfo = reactive({ // 主机信息
  hostname: '' as string | undefined,
  platform: '' as string | undefined,
  modelName: '' as string | undefined,
  kernel_version: '' as string | undefined,
  architecture: '' as string | undefined,
})
const networkInterfaces = ref<any[] | undefined>([]); // 网卡信息
const currentNetworkInterface = ref('') // 当前选中的网卡
const favoriteNetworkInterface = ref('') // 收藏的网卡
const settingsVisible = ref(false);
const settings = reactive({
  fontFamily: '',
  fontSize: 16,
  cacheDir: ''
});

const systemInfoInterval = ref<number | null>(null)
const cpuView = ref<'large' | 'small'>('large')
// 生命周期
onMounted(async () => {
  console.log('mount');

  // 加载收藏的网卡
  try {
    const favorite = await AppSettingsService.GetFavoriteNetworkInterface();
    if (favorite) {
      favoriteNetworkInterface.value = favorite;
    }
  } catch (error) {
    console.error('加载收藏网卡失败:', error);
  }

  Events.On('monitor:sysInfo', (event) => {
    console.log(event.data);
    // cpu信息
    const cpuInfo = event.data.cpu_info;
    cpuUsage.value = cpuInfo?.core_usages
    // 主机信息
    hardwareInfo.hostname = event.data.host_info.hostname || ''
    hardwareInfo.platform = event.data.host_info.platform || ''
    hardwareInfo.kernel_version = event.data.host_info.kernel_version || ''
    hardwareInfo.modelName = cpuInfo?.core_info.modelName || ''
    hardwareInfo.architecture = event.data.host_info.architecture || ''
    // 内存信息
    memoryUsage.value = event.data.memory_info.used_percent
    memoryTotal.value = event.data.memory_info.total
    memoryUsed.value = event.data.memory_info.used
    // 硬盘信息
    diskUsage.total = event.data.disk_info.total
    diskUsage.used = event.data.disk_info.used
    diskUsage.used_percent = event.data.disk_info.used_percent
    // 网卡信息
    networkInterfaces.value = event.data.network_interfaces
    if (currentNetworkInterface.value === '') {
      // 优先使用收藏的网卡,如果收藏的网卡存在且可用
      if (favoriteNetworkInterface.value && networkInterfaces.value?.some(n => n.name === favoriteNetworkInterface.value)) {
        currentNetworkInterface.value = favoriteNetworkInterface.value;
      } else {
        currentNetworkInterface.value = networkInterfaces.value?.[0]?.name
      }
    }
  });
});
onUnmounted(() => {
  console.log('unmount');
  Events.Off('monitor:sysInfo');
});

// 方法
const openSettings = () => {
  settingsVisible.value = true;
};

const closeSettings = () => {
  settingsVisible.value = false;
};

const saveSettings = () => {
  // 保存设置（这里可以添加保存逻辑）
  closeSettings();
};

const updateSettings = (newSettings: any) => {
  Object.assign(settings, newSettings);
};

const navigateToTool = (path: string) => {
  router.push(path);
};
// 监听器
const focus = () => {
  console.log('focus');
};
// 选择网卡
const handleChange = (value: string) => {
  console.log(`selected ${value}`);
};

// 收藏网卡
const handleFavorite = async () => {
  if (!currentNetworkInterface.value) {
    message.warning('请先选择一个网卡');
    return;
  }

  try {
    await AppSettingsService.SetFavoriteNetworkInterface(currentNetworkInterface.value);
    favoriteNetworkInterface.value = currentNetworkInterface.value;
    message.success(`已收藏网卡: ${currentNetworkInterface.value}`);
  } catch (error) {
    message.error('收藏失败');
    console.error('收藏网卡失败:', error);
  }
};

// 取消收藏
const handleUnfavorite = async () => {
  try {
    await AppSettingsService.SetFavoriteNetworkInterface('');
    favoriteNetworkInterface.value = '';
    message.success('已取消收藏');
  } catch (error) {
    message.error('取消收藏失败');
    console.error('取消收藏失败:', error);
  }
};

// 判断当前网卡是否已收藏
const isCurrentFavorite = () => {
  return currentNetworkInterface.value === favoriteNetworkInterface.value && favoriteNetworkInterface.value !== '';
};

// 获取系统信息
const updateSystemInfo = async () => {
  let cpuInfo = await SysInfoService.GetCPUInfo()
  cpuUsage.value = cpuInfo?.core_usages
  hardwareInfo.modelName = cpuInfo?.core_info.modelName || ''

  let MemInfo = await SysInfoService.GetMemoryInfo()
  memoryUsage.value = MemInfo?.used_percent
  memoryTotal.value = MemInfo?.total
  memoryUsed.value = MemInfo?.used;
  console.log(MemInfo);
  let hostInfo = await SysInfoService.GetHostInfo()
  hardwareInfo.hostname = hostInfo?.hostname || ''
  hardwareInfo.platform = hostInfo?.platform || ''
  hardwareInfo.kernel_version = hostInfo?.kernel_version || ''

  let diskInfo = await SysInfoService.GetDiskInfo()
  diskUsage.total = diskInfo?.total || 0
  diskUsage.used = diskInfo?.used || 0
  diskUsage.used_percent = diskInfo?.used_percent || 0

  let networkInfo = await SysInfoService.GetNetworkInterfaces()
  networkInterfaces.value = networkInfo || []
  // console.log(currentNetworkInterface.value[0])
  if (currentNetworkInterface.value.length === 0) {
    // console.log("no network interface")
    currentNetworkInterface.value = networkInterfaces.value?.[0]?.name
  }
};

// 进度条根据百分比值自定义分段的颜色
const getProgressColor = (numberPercent: any) => {
  let color = ''
  if (numberPercent < 50) {
    color = '#87d068'
  } else if (numberPercent >= 50 && numberPercent < 80) {
    color = '#faad14'
  } else if (numberPercent >= 80) {
    color = '#f5222d'
  }
  return color
}

// 环形进度颜色（渐变）
const ringColor = (pct: any) => {
  if ((pct ?? 0) < 50)  return { '0%': '#34d399', '100%': '#10b981' }
  if ((pct ?? 0) < 80)  return { '0%': '#fbbf24', '100%': '#f59e0b' }
  return { '0%': '#f87171', '100%': '#ef4444' }
}

// CPU 热力颜色
const heatColor = (pct: number) => {
  if (pct < 30)  return '#34d399'
  if (pct < 60)  return '#60a5fa'
  if (pct < 80)  return '#fbbf24'
  return '#f87171'
}

const formatBytes = (bytes: number) => {
  if (bytes === 0) return "0B";

  const units = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
  let unitIndex = 0;
  let size = bytes;

  // 不断除以1024直到找到合适的单位
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }

  // 格式化数字，保留1位小数
  const formattedSize = size.toFixed(1);

  // 去除不必要的零和小数点
  const finalSize = formattedSize.includes('.0')
      ? formattedSize.replace('.0', '')
      : formattedSize;

  return finalSize + units[unitIndex];
}


</script>

<style scoped>
/* ── 根容器：左右分栏 ── */
.db {
  width: 100%;
  height: 100%;
  display: flex;
  background: #f0f4ff;
  overflow: hidden;
  gap: 0;
}

/* ══════════════════════════════
   左侧竖栏
══════════════════════════════ */
.db-aside {
  width: 240px;
  flex-shrink: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 0;
  background: #ffffff;
  border-right: 1px solid #eef0f6;
  overflow: hidden;
  padding: 14px 14px;
}

/* 主机名区 */
.aside-host {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}
.host-dot {
  width: 8px; height: 8px;
  border-radius: 50%;
  background: #cbd5e1;
  flex-shrink: 0;
  transition: background 0.3s;
}
.host-dot.online {
  background: #34d399;
  box-shadow: 0 0 0 3px rgba(52,211,153,0.2);
}
.host-name {
  font-size: 13px;
  font-weight: 700;
  color: #1e1b4b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.host-os {
  font-size: 11px;
  color: #94a3b8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.aside-divider {
  height: 1px;
  background: #f1f5f9;
  margin: 10px 0;
  flex-shrink: 0;
}

/* 指标块 */
.metric-block {
  display: flex;
  flex-direction: column;
  gap: 5px;
  flex-shrink: 0;
}
.metric-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}
.metric-label {
  font-size: 11px;
  font-weight: 600;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.8px;
}
.metric-val {
  font-size: 13px;
  font-weight: 700;
  color: #1e1b4b;
}
.ring-wrap {
  display: flex;
  justify-content: center;
  flex-shrink: 0;
}
.metric-sub {
  font-size: 11px;
  color: #cbd5e1;
  text-align: center;
}

/* 网卡 */
.net-block { gap: 6px; flex: 1; min-height: 0; overflow: hidden; }
.net-ctrl { display: flex; align-items: center; gap: 4px; }
.fav-on { color: #f59e0b; }

.net-detail { display: flex; flex-direction: column; gap: 4px; overflow: hidden; flex: 1; min-height: 0; }
.nd-section-label {
  font-size: 10px;
  font-weight: 600;
  color: #cbd5e1;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  margin-top: 0;
}
.nd-row {
  display: flex;
  align-items: center;
  gap: 6px;
}
.nd-k {
  font-size: 10px;
  color: #94a3b8;
  width: 28px;
  flex-shrink: 0;
}
.nd-v {
  font-size: 11px;
  color: #334155;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}
.mono { font-family: monospace; }
.nd-tag {
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 6px;
  font-weight: 600;
  flex-shrink: 0;
}
.nd-tag.up   { background: #dcfce7; color: #16a34a; }
.nd-tag.down { background: #fee2e2; color: #dc2626; }

.nd-traffic {
  display: flex;
  gap: 6px;
}
.traf-item {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  border-radius: 10px;
}
.traf-item.rx { background: #eff6ff; }
.traf-item.tx { background: #fff7ed; }
.traf-arrow {
  font-size: 14px;
  font-weight: 700;
}
.traf-item.rx .traf-arrow { color: #3b82f6; }
.traf-item.tx .traf-arrow { color: #f97316; }
.traf-num {
  font-size: 12px;
  font-weight: 700;
  color: #1e293b;
  white-space: nowrap;
}
.traf-label {
  font-size: 10px;
  color: #94a3b8;
}

/* ══════════════════════════════
   右侧主区
══════════════════════════════ */
.db-main {
  flex: 1;
  min-width: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  overflow: hidden;
}

/* 通用面板 */
.panel {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #eef0f6;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
}

.panel-header {
  display: flex;
  align-items: baseline;
  gap: 10px;
  flex-shrink: 0;
}
.panel-title {
  font-size: 13px;
  font-weight: 700;
  color: #1e1b4b;
  white-space: nowrap;
}
.panel-sub {
  font-size: 11px;
  color: #94a3b8;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* CPU 热力图 */
.cpu-panel { flex: 1; min-height: 0; }

/* 视图切换按钮 */
.view-toggle {
  display: flex;
  gap: 4px;
  margin-left: auto;
}
.vt-btn {
  width: 26px; height: 26px;
  border: 1px solid #e2e8f0;
  border-radius: 7px;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  font-size: 14px;
  transition: all 0.15s;
}
.vt-btn:hover { background: #f1f5f9; color: #475569; }
.vt-btn.active { background: #6366f1; border-color: #6366f1; color: #fff; }
.vt-icon { line-height: 1; }

/* 网格容器 */
.cpu-heatmap {
  flex: 1;
  display: grid;
  gap: 8px;
  overflow: auto;
  align-content: start;
}
.view-large { grid-template-columns: repeat(5, 1fr); }
.view-small { grid-template-columns: repeat(8, 1fr); }

/* 核心卡片基础 */
.cpu-core {
  border-radius: 12px;
  background: #f8fafc;
  border: 1px solid #f1f5f9;
  transition: border-color 0.2s, box-shadow 0.2s;
  overflow: hidden;
}
.cpu-core:hover {
  border-color: var(--color);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--color) 12%, transparent);
}

/* ── 大视图 ── */
.cpu-core.large {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0;
  padding: 10px 8px 10px;
}
.core-ring {
  position: relative;
  width: 52px; height: 52px;
}
.core-svg {
  width: 100%; height: 100%;
  transform: rotate(-90deg);
}
.core-ring-pct {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  color: var(--color);
}
.core-label-lg {
  position: absolute;
  top: 7px;
  left: 9px;
  font-size: 10px;
  font-weight: 700;
  color: #94a3b8;
  line-height: 1;
}

/* ── 小视图 ── */
.cpu-core.small {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 7px 7px 6px;
}
.core-bar-bg {
  height: 3px;
  border-radius: 3px;
  background: #e2e8f0;
  overflow: hidden;
}
.core-bar-fill {
  height: 100%;
  width: calc(var(--pct) * 1%);
  background: var(--color);
  border-radius: 3px;
  transition: width 0.6s ease;
}
.core-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.core-idx { font-size: 9px; color: #94a3b8; font-weight: 600; }
.core-pct { font-size: 10px; font-weight: 700; color: var(--color); }

/* 主机详情 */
.host-panel { flex-shrink: 0; }

.host-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0;
}
.hg-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 10px 12px;
  border-bottom: 1px solid #f1f5f9;
  border-right: 1px solid #f1f5f9;
}
.hg-item:nth-child(even) { border-right: none; }
.hg-item:nth-last-child(-n+2) { border-bottom: none; }

.hg-k {
  font-size: 10px;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.hg-v {
  font-size: 12px;
  font-weight: 600;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>