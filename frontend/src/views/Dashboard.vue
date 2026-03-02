<template>
  <div class="desktop-app">
    <!-- 主内容区 -->
    <main class="app-main">
      <a-row :gutter="24">
        <a-col :span="8">
          <a-spin :spinning="memoryUsage === 0 " tip="Loading...">
            <a-card :bordered="false">
              <template #title>
                <div class="flex items-center justify-between">
                  <span>💾 内存</span>
                  <span class="ml-5 text-xs">
                    Used：{{ memoryUsed }}
                  </span>
                </div>
              </template>
              <div class="relative text-center h-28">
                <a-progress type="dashboard" :stroke-color="getProgressColor(memoryUsage)" :strokeWidth="8"
                            :percent="memoryUsage"/>
              </div>
            </a-card>
          </a-spin>
        </a-col>
        <a-col :span="8">
          <a-spin :spinning="diskUsage.used_percent === 0 " tip="Loading...">
            <a-card :bordered="false">
              <template #title>
                <div class="flex items-center justify-between">
                  <span>💿 硬盘</span>
                  <span class="ml-5 text-xs">
                  Used：{{ formatBytes(diskUsage.used) }}
                </span>
                </div>
              </template>
              <div class="flex flex-col items-center justify-around h-28">
                <a-progress type="circle" :stroke-color="getProgressColor(diskUsage.used_percent)" :strokeWidth="8"
                            status="active" :percent="diskUsage.used_percent"/>
              </div>
            </a-card>
          </a-spin>
        </a-col>
        <a-col :span="8">
          <a-spin :spinning="memoryUsage === 0 " tip="Loading...">
            <a-card :bordered="false">
              <template #title>
                <div class="flex items-center justify-between">
                  <span>🚀 网卡</span>
                  <a-select
                      ref="select"
                      v-model:value="currentNetworkInterface"
                      style="width: 100px"
                      size="small"
                      @focus="focus"
                      @change="handleChange"
                  >
                    <a-select-option v-for="item in networkInterfaces" :key="item.name" :value="item.name">{{
                        item.name
                      }}
                    </a-select-option>
                  </a-select>
                </div>
              </template>
              <div class="h-28">
                <template v-for="(item,index) in networkInterfaces" :key="index">
                  <div v-if="item.name == currentNetworkInterface">
                    <a-descriptions size="small" :column="1">
                      <a-descriptions-item label="IP地址">
                        <a-tooltip :title=item.name>
                          <a-tag class="interface-name">{{ item.ipv4 }}</a-tag>
                        </a-tooltip>
                        <a-tag color="success" v-if="item.isUp">已连接</a-tag>
                        <a-tag color="error" v-else>未连接</a-tag>
                      </a-descriptions-item>
                      <a-descriptions-item label="MAC地址" class="font-bold">
                        {{ item.mac }}
                      </a-descriptions-item>
                      <a-descriptions-item label="Bytes">
                        <div class="text-[12px]">
                          <a-tag color="#87CEEB">RX</a-tag>
                          {{ formatBytes(item.bytesRecv) }}
                        </div>
                        <div class="ml-2 text-[12px]">
                          <a-tag color="#FF9800" >TX</a-tag>
                          {{ formatBytes(item.bytesSent) }}
                        </div>
                      </a-descriptions-item>
                      <a-descriptions-item label="Packets">
                        <div class="text-[12px]">
                          <a-tag color="#87CEEB">RX</a-tag>
                          {{ item.packetsRecv }}
                        </div>
                        <div class="ml-2 text-[12px]">
                          <a-tag color="#FF9800" class="text-[12px]">TX</a-tag>
                          {{ item.packetsSent }}
                        </div>
                      </a-descriptions-item>
                    </a-descriptions>
                  </div>
                </template>
              </div>
            </a-card>
          </a-spin>
        </a-col>
      </a-row>
      <a-row :gutter="24" class="mt-5">
        <a-col :span="12">
          <a-spin :spinning="cpuUsage?.length ===0" tip="Loading...">
            <a-card title="🖥️ CPU负载情况" :bordered="false" size="small" class="box-large" >
              <div class="flex flex-wrap items-center justify-between h-full">
                <div class="flex items-center justify-between w-3/6 mt-4" v-for="(item, index) in cpuUsage"
                     :key="index">
                  <span class="circle_number">{{ index }}</span>
                  <a-progress :steps="10" :size="[10, 10]" :percent="item"
                              :stroke-color="getProgressColor(item)" class="m-0 ml-2"/>
                </div>
              </div>
            </a-card>
          </a-spin>
        </a-col>
        <a-col :span="12">
          <a-spin :spinning="cpuUsage?.length ===0" tip="Loading...">
            <a-card title="ℹ️ 主机详情" :bordered="false" size="small" class="box-large" >
              <a-descriptions :column="1" bordered>
                <a-descriptions-item label="设备名称" :span="1">{{ hardwareInfo.hostname }}</a-descriptions-item>
                <a-descriptions-item label="版本" :span="1">{{ hardwareInfo.platform }}</a-descriptions-item>
                <a-descriptions-item label="处理器" :span="1">{{ hardwareInfo.modelName }}</a-descriptions-item>
                <a-descriptions-item label="内存" :span="1">{{ memoryTotal }}</a-descriptions-item>
                <a-descriptions-item label="系统版本" :span="1">{{ hardwareInfo.kernel_version }}</a-descriptions-item>
                <a-descriptions-item label="系统架构" :span="1">{{ hardwareInfo.architecture }}</a-descriptions-item>
                <a-descriptions-item label="存储" :span="1">{{ formatBytes(diskUsage.total) }}</a-descriptions-item>
              </a-descriptions>
            </a-card>
          </a-spin>
        </a-col>
      </a-row>
    </main>

    <!-- 底部状态栏 -->
    <!--    <footer class="app-footer">-->
    <!--      <div class="runtime-info">运行时间: {{ runtime }}</div>-->
    <!--      <div class="copyright">© 2023 Happy Tools. All rights reserved.</div>-->
    <!--    </footer>-->

    <!-- 设置对话框 -->
    <a-modal v-model:open="settingsVisible" title="设置" @ok="saveSettings" @cancel="closeSettings">
      <settings-panel :settings="settings" @update:settings="updateSettings"/>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import {onMounted, onUnmounted, reactive, ref} from 'vue';
import {useRouter} from 'vue-router';
import SettingsPanel from '@/components/SettingsPanel.vue';
import {useSettingsStore} from '@/stores/settings';
import {SysInfoService} from "../../bindings/github.com/Aliuyanfeng/happytools/backend/services/monitor";
import {Events} from "@wailsio/runtime";

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
const settingsVisible = ref(false);
const settings = reactive({
  fontFamily: '',
  fontSize: 16,
  cacheDir: ''
});

const systemInfoInterval = ref<number | null>(null)
// 生命周期
onMounted(() => {
  console.log('mount');
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
      currentNetworkInterface.value = networkInterfaces.value?.[0]?.name
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
/* 基础样式 */
.desktop-app {
  display: flex;
  flex-direction: column;
  background-color: #f5f5f5;
  justify-content: space-between;
  /* height: 100%; */
}

/* 主内容区 */
.app-main {
  flex: 1;
  padding: 16px;
  overflow-y: auto;

  .interface-name {
    max-width: 120px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .box-large {
    height: 445px;
  }
}

.circle_number {
  width: 20px;
  line-height: 20px;
  background-color: #fff;
  color: #999;
  border-radius: 50%;
}
</style>