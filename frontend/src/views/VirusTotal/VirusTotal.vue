<template>
  <div class="vt p-6">
    <a-card :bodyStyle="{ padding: 0 }">
      <!-- 配额信息条 -->
      <div class="quota-header" v-if="apiQuota || quotaLoading">
        <div class="quota-header-left">
          <span class="quota-header-title">
            <SafetyCertificateOutlined class="mr-1" />
            API Quota
          </span>
          <a-spin :spinning="quotaLoading" size="small" />
        </div>
        <div class="quota-header-items" v-if="apiQuota">
          <div class="quota-stat-item" v-for="(item, key) in quotaDisplayList" :key="key">
            <div class="quota-stat-header">
              <span class="quota-stat-label">{{ item.label }}</span>
              <span class="quota-stat-nums">
                <span :style="{ color: quotaColor(item.quota) }">{{ formatQuotaNum(item.quota.used) }}</span>
                <span class="quota-stat-sep">/</span>
                <span class="quota-stat-total">{{ formatQuotaNum(item.quota.allowed) }}</span>
              </span>
            </div>
            <a-progress
              :percent="quotaPercent(item.quota)"
              :stroke-color="quotaColor(item.quota)"
              :show-info="false"
              size="small"
              class="quota-progress"
            />
          </div>
        </div>
      </div>

      <!-- 任务列表 header -->
      <div class="task-list-header mt-5">
        <span class="task-list-title">{{ t('virustotal.taskList') }}</span>
        <a-button type="primary" :icon="h(PlusOutlined)" @click="openCreateTaskModal">{{ t('virustotal.createTask') }}</a-button>
      </div>
      <a-table class="quota-table" :dataSource="tasks" :columns="columns" rowKey="id" :loading="loading" :pagination="{ pageSize: 10 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'taskType'">
            <a-tag :color="record.taskType === 'directory' ? 'blue' : 'green'">
              {{ record.taskType === 'directory' ? t('virustotal.directoryScan') : t('virustotal.singleFile') }}
            </a-tag>
          </template>
          <template v-if="column.key === 'totalFiles'">
            <span>{{ record.taskType === 'directory' ? (record.totalFiles || 0) : 1 }}</span>
          </template>
          <template v-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
            <span v-if="record.taskType === 'directory' && record.status === 'scanning'" class="ml-1 text-gray-500 text-xs">
              ({{ record.completedFiles || 0 }}/{{ record.totalFiles || 0 }})
            </span>
          </template>
          <template v-if="column.key === 'detectionRate'">
            <!-- 统一格式: 恶意数/总数 -->
            <template v-if="record.status === 'completed'">
              <span :class="getDetectionClass(record)">
                {{ getDetectionText(record) }}
              </span>
            </template>
            <span v-else class="text-gray-400">-</span>
          </template>
          <template v-if="column.key === 'scanTime'">
            {{ record.scanTime || '-' }}
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button type="link" size="small" @click="viewTaskDetail(record)">
                {{ t('virustotal.viewDetail') }}
              </a-button>
              <a-button
                type="link"
                size="small"
                @click="refreshTaskStatus(record)"
                :loading="record.refreshing"
                v-if="record.status !== 'completed' && record.status !== 'failed'"
              >
                {{ t('virustotal.refreshStatus') }}
              </a-button>
              <a-popconfirm
                :title="t('virustotal.deleteConfirm')"
                @confirm="deleteTask(record.id)"
                :okText="t('common.confirm')"
                :cancelText="t('common.cancel')"
              >
                <a-button type="link" size="small" danger>{{ t('common.delete') }}</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 创建任务弹窗 -->
    <a-modal
      v-model:open="createModalVisible"
      :title="t('virustotal.createTask')"
      centered
      :maskClosable="false"
      @ok="createTask"
      @cancel="createCancel"
      :confirmLoading="createLoading"
    >
      <a-alert
        v-if="!settingsStore.vtApiKey"
        :message="t('settings.apiKeyNotConfigured')"
        :description="t('settings.apiKeyNotConfiguredDesc')"
        type="warning"
        show-icon
        class="mb-4"
      />

      <a-form :model="formState" layout="vertical">
        <a-form-item :label="t('virustotal.scanType')" required>
          <a-radio-group v-model:value="formState.scanType">
            <a-radio value="file">{{ t('virustotal.singleFileScan') }}</a-radio>
            <a-radio value="directory">{{ t('virustotal.batchDirectory') }}</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item
          :label="t('virustotal.filePath')"
          v-if="formState.scanType === 'file'"
          required
        >
          <a-input-group compact>
            <a-input v-model:value="formState.filePath" readonly style="width: calc(100% - 100px)" />
            <a-button type="primary" @click="chooseFile">
              {{ t('virustotal.selectFile') }}
            </a-button>
          </a-input-group>
        </a-form-item>

        <a-form-item
          :label="t('virustotal.directoryPath')"
          v-if="formState.scanType === 'directory'"
          required
        >
          <a-input-group compact>
            <a-input v-model:value="formState.directoryPath" readonly style="width: calc(100% - 100px)" />
            <a-button type="primary" @click="chooseDirectory">
              {{ t('virustotal.selectDirectory') }}
            </a-button>
          </a-input-group>
          <div class="text-gray-500 text-xs mt-1">{{ t('virustotal.recursiveScanHint') }}</div>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 任务详情弹窗 - 统一显示文件列表 -->
    <a-modal
      v-model:open="detailModalVisible"
      :title="currentTask?.taskType === 'directory' ? t('virustotal.directoryScanResult') : t('virustotal.scanResult')"
      :width="currentTask?.taskType === 'directory' ? '1200px' : '1000px'"
      :footer="null"
    >
      <div v-if="currentTask" class="task-detail">
        <!-- 基本信息 -->
        <a-descriptions :column="currentTask?.taskType === 'directory' ? 3 : 2" bordered size="small">
          <a-descriptions-item :label="t('virustotal.taskType')">
            <a-tag :color="currentTask.taskType === 'directory' ? 'blue' : 'green'">
              {{ currentTask.taskType === 'directory' ? t('virustotal.directoryScan') : t('virustotal.singleFile') }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item :label="t('virustotal.fileInfo')">
            <a-tooltip :title="currentTask.filePath">
              <span class="mono-text">{{ currentTask.fileName || '未知' }}</span>
            </a-tooltip>
          </a-descriptions-item>
          <a-descriptions-item :label="t('virustotal.taskStatus')">
            <a-tag :color="getStatusColor(currentTask.status)">
              {{ getStatusText(currentTask.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item :label="t('virustotal.totalFiles')">{{ currentTask.taskType === 'directory' ? (currentTask.totalFiles || 0) : 1 }}</a-descriptions-item>
          <a-descriptions-item :label="t('virustotal.scanTime')" v-if="currentTask.taskType === 'directory'">
            {{ currentTask.completedFiles || 0 }} / {{ currentTask.totalFiles || 0 }}
            <a-progress
              :percent="directoryProgress"
              size="small"
              style="width: 100px; margin-left: 8px"
            />
          </a-descriptions-item>
          <a-descriptions-item :label="t('virustotal.scanTime')">{{ currentTask.scanTime || '-' }}</a-descriptions-item>
        </a-descriptions>

        <!-- 统计概览 - 目录任务 -->
        <template v-if="currentTask.taskType === 'directory'">
          <a-divider orientation="left">{{ t('virustotal.scanOverview') }}</a-divider>
          <div class="stats-container">
            <div class="stat-item stat-malicious">
              <div class="stat-count">{{ currentTask.maliciousFiles || 0 }}</div>
              <div class="stat-label">{{ t('virustotal.maliciousFiles') }}</div>
            </div>
            <div class="stat-item stat-suspicious">
              <div class="stat-count">{{ currentTask.suspiciousFiles || 0 }}</div>
              <div class="stat-label">{{ t('virustotal.suspiciousFiles') }}</div>
            </div>
            <div class="stat-item stat-safe">
              <div class="stat-count">{{ (currentTask.completedFiles || 0) - (currentTask.maliciousFiles || 0) - (currentTask.suspiciousFiles || 0) }}</div>
              <div class="stat-label">{{ t('virustotal.safeFiles') }}</div>
            </div>
            <div class="stat-item stat-pending">
              <div class="stat-count">{{ (currentTask.totalFiles || 0) - (currentTask.completedFiles || 0) }}</div>
              <div class="stat-label">{{ t('virustotal.pendingScan') }}</div>
            </div>
          </div>
        </template>

        <!-- 统计概览 - 单文件任务 -->
        <template v-else>
          <a-divider orientation="left">{{ t('virustotal.scanOverview') }}</a-divider>
          <div class="stats-container">
            <div class="stat-item stat-malicious">
              <div class="stat-count">{{ (currentTask.stats?.malicious || 0) > 0 ? 1 : 0 }}</div>
              <div class="stat-label">{{ t('virustotal.maliciousFiles') }}</div>
            </div>
            <div class="stat-item stat-suspicious">
              <div class="stat-count">{{ (currentTask.stats?.suspicious || 0) > 0 && (currentTask.stats?.malicious || 0) === 0 ? 1 : 0 }}</div>
              <div class="stat-label">{{ t('virustotal.suspiciousFiles') }}</div>
            </div>
            <div class="stat-item stat-safe">
              <div class="stat-count">{{ (currentTask.stats?.malicious || 0) === 0 && (currentTask.stats?.suspicious || 0) === 0 && currentTask.status === 'completed' ? 1 : 0 }}</div>
              <div class="stat-label">{{ t('virustotal.safeFiles') }}</div>
            </div>
            <div class="stat-item stat-pending">
              <div class="stat-count">{{ currentTask.status !== 'completed' ? 1 : 0 }}</div>
              <div class="stat-label">{{ t('virustotal.pendingScan') }}</div>
            </div>
          </div>
        </template>

        <!-- 文件列表 -->
        <a-divider orientation="left">
          {{ t('virustotal.fileList') }}
          <span class="result-count">({{ t('virustotal.filesTotal', { count: taskFiles.length }) }})</span>
        </a-divider>

        <!-- 文件筛选 - 仅目录任务显示 -->
        <div class="filter-section" v-if="currentTask.taskType === 'directory'">
          <a-radio-group v-model:value="fileFilter" button-style="solid" size="small">
            <a-radio-button value="all">{{ t('virustotal.filter.all') }}</a-radio-button>
            <a-radio-button value="malicious">{{ t('virustotal.filter.malicious') }}</a-radio-button>
            <a-radio-button value="suspicious">{{ t('virustotal.filter.suspicious') }}</a-radio-button>
            <a-radio-button value="safe">{{ t('virustotal.filter.safe') }}</a-radio-button>
            <a-radio-button value="pending">{{ t('virustotal.filter.scanning') }}</a-radio-button>
          </a-radio-group>
        </div>

        <a-table
          :dataSource="filteredFiles"
          :columns="fileColumns"
          rowKey="id"
          size="small"
          :pagination="currentTask.taskType === 'directory' ? { pageSize: 10 } : false"
          :loading="filesLoading"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'fileName'">
              <a-tooltip :title="record.filePath">
                <span class="cursor-pointer text-blue-500">
                  <FileFilled class="mr-1" style="color: #91d5ff" />
                  {{ record.fileName || '未知' }}
                </span>
              </a-tooltip>
            </template>
            <template v-if="column.key === 'fileSize'">
              {{ formatFileSize(record.fileSize) }}
            </template>
            <template v-if="column.key === 'status'">
              <a-tag :color="getFileStatusColor(record.status)">
                {{ getFileStatusText(record.status) }}
              </a-tag>
            </template>
            <template v-if="column.key === 'detectionRate'">
              <template v-if="record.status === 'completed'">
                <a-progress
                  :percent="record.detectionRate || 0"
                  :status="record.detectionRate > 0 ? 'exception' : 'success'"
                  size="small"
                  style="width: 80px"
                />
                <span class="ml-1 text-gray-500 text-xs">
                  ({{ record.stats?.malicious || 0 }}/{{ record.totalEngines || 0 }})
                </span>
              </template>
              <span v-else class="text-gray-400">-</span>
            </template>
            <template v-if="column.key === 'action'">
              <a-button
                type="link"
                size="small"
                @click="viewFileDetail(record)"
                v-if="record.status === 'completed'"
              >
                {{ t('virustotal.viewDetail') }}
              </a-button>
              <a-button
                type="link"
                size="small"
                @click="refreshFileStatus(record)"
                :loading="record.refreshing"
                v-else-if="record.status !== 'failed'"
              >
                {{ t('virustotal.refreshStatus') }}
              </a-button>
            </template>
          </template>
        </a-table>
      </div>
    </a-modal>

    <!-- 文件详情弹窗 -->
    <a-modal
      v-model:open="fileDetailModalVisible"
      :title="t('virustotal.fileScanDetail')"
      width="900px"
      :footer="null"
    >
      <div v-if="currentFile" class="task-detail">
        <!-- 基本信息 -->
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item :label="t('virustotal.fileInfo')">{{ currentFile.fileName || '未知' }}</a-descriptions-item>
          <a-descriptions-item label="Size">{{ formatFileSize(currentFile.fileSize) }}</a-descriptions-item>
          <a-descriptions-item :label="t('virustotal.taskStatus')">
            <a-tag :color="getFileStatusColor(currentFile.status)">
              {{ getFileStatusText(currentFile.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item :label="t('virustotal.detectionRate')">
            <a-progress
              :percent="currentFile.detectionRate || 0"
              :status="currentFile.detectionRate > 0 ? 'exception' : 'success'"
              style="width: 150px"
            />
            <span class="ml-2 text-gray-500">
              ({{ currentFile.stats?.malicious || 0 }}/{{ currentFile.totalEngines || 0 }})
            </span>
          </a-descriptions-item>
          <a-descriptions-item label="MD5" :span="2">
            <span class="mono-text hash-value">{{ currentFile.md5 || '...' }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="SHA256" :span="2">
            <span class="mono-text hash-value">{{ currentFile.sha256 || '...' }}</span>
          </a-descriptions-item>
        </a-descriptions>

        <!-- 统计信息 -->
        <a-divider orientation="left">{{ t('virustotal.detectionStats') }}</a-divider>
        <div class="stats-container" v-if="currentFile.stats">
          <div class="stat-item stat-malicious">
            <div class="stat-count">{{ currentFile.stats.malicious || 0 }}</div>
            <div class="stat-label">{{ t('virustotal.malicious') }}</div>
          </div>
          <div class="stat-item stat-suspicious">
            <div class="stat-count">{{ currentFile.stats.suspicious || 0 }}</div>
            <div class="stat-label">{{ t('virustotal.suspicious') }}</div>
          </div>
          <div class="stat-item stat-harmless">
            <div class="stat-count">{{ currentFile.stats.harmless || 0 }}</div>
            <div class="stat-label">{{ t('virustotal.harmless') }}</div>
          </div>
          <div class="stat-item stat-undetected">
            <div class="stat-count">{{ currentFile.stats.undetected || 0 }}</div>
            <div class="stat-label">{{ t('virustotal.undetected') }}</div>
          </div>
          <div class="stat-item stat-unsupported">
            <div class="stat-count">{{ currentFile.stats.typeUnsupported || 0 }}</div>
            <div class="stat-label">{{ t('virustotal.unsupported') }}</div>
          </div>
          <div class="stat-item stat-timeout">
            <div class="stat-count">{{ currentFile.stats.timeout || 0 }}</div>
            <div class="stat-label">{{ t('virustotal.timeout') }}</div>
          </div>
          <div class="stat-item stat-failure">
            <div class="stat-count">{{ currentFile.stats.failure || 0 }}</div>
            <div class="stat-label">{{ t('virustotal.failure') }}</div>
          </div>
        </div>

        <!-- 检测结果 -->
        <a-divider orientation="left">
          {{ t('virustotal.detectionResults') }}
          <span class="result-count" v-if="currentFile.results">
            ({{ t('virustotal.enginesTotal', { count: currentFile.results.length }) }})
          </span>
        </a-divider>

        <!-- 分类筛选 -->
        <div class="filter-section" v-if="currentFile.results && currentFile.results.length > 0">
          <a-radio-group v-model:value="fileResultFilter" button-style="solid" size="small">
            <a-radio-button value="all">{{ t('virustotal.filter.all') }}</a-radio-button>
            <a-radio-button value="malicious">{{ t('virustotal.filter.malicious') }}</a-radio-button>
            <a-radio-button value="suspicious">{{ t('virustotal.filter.suspicious') }}</a-radio-button>
            <a-radio-button value="harmless">{{ t('virustotal.filter.safe') }}</a-radio-button>
            <a-radio-button value="undetected">{{ t('virustotal.undetected') }}</a-radio-button>
            <a-radio-button value="type-unsupported">{{ t('virustotal.unsupported') }}</a-radio-button>
          </a-radio-group>
        </div>

        <!-- 厂商检测结果 - Tag 展示 -->
        <div class="results-container" v-if="filteredFileResults.length > 0">
          <a-tooltip
            v-for="result in filteredFileResults"
            :key="result.engine"
            placement="top"
          >
            <template #title>
              <div class="result-tooltip">
                <div><strong>Engine:</strong> {{ result.engine }}</div>
                <div><strong>Category:</strong> {{ getCategoryText(result.category) }}</div>
                <div><strong>Result:</strong> {{ result.result || '-' }}</div>
                <div><strong>Method:</strong> {{ result.method || '-' }}</div>
              </div>
            </template>
            <a-tag
              :color="getCategoryColor(result.category)"
              class="engine-tag"
            >
              {{ result.engine }}
            </a-tag>
          </a-tooltip>
        </div>
        <a-empty v-else :description="t('virustotal.noResults')" />
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import {
  PlusOutlined,
  FileFilled,
  FolderFilled,
  SafetyCertificateOutlined,
} from '@ant-design/icons-vue';
import { h, onMounted, onUnmounted, ref, reactive, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { message } from 'ant-design-vue';
import { VTService } from "../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/vt";
import { useSettingsStore } from '@/stores/settings';

const { t } = useI18n();
const settingsStore = useSettingsStore();
const loading = ref(false);
const createModalVisible = ref(false);
const detailModalVisible = ref(false);
const fileDetailModalVisible = ref(false);
const createLoading = ref(false);
const resultFilter = ref('all');
const fileFilter = ref('all');
const fileResultFilter = ref('all');
const filesLoading = ref(false);

const formState = reactive({
  scanType: 'file',
  filePath: '',
  directoryPath: ''
});

const currentTask = ref<any>(null);
const currentFile = ref<any>(null);
const taskFiles = ref<any[]>([]);

// 任务列表
const tasks = ref<any[]>([]);

// 轮询定时器
let pollingTimer: ReturnType<typeof setInterval> | null = null;

const columns = [
  {
    title: t('virustotal.taskType'),
    dataIndex: 'taskType',
    key: 'taskType',
    width: 100
  },
  {
    title: t('virustotal.totalFiles'),
    dataIndex: 'totalFiles',
    key: 'totalFiles',
    width: 100
  },
  {
    title: t('virustotal.taskStatus'),
    dataIndex: 'status',
    key: 'status',
    width: 130
  },
  {
    title: t('virustotal.detectionStatus'),
    dataIndex: 'detectionRate',
    key: 'detectionRate',
    width: 120
  },
  {
    title: t('virustotal.scanTime'),
    dataIndex: 'scanTime',
    key: 'scanTime',
    width: 180
  },
  {
    title: t('virustotal.actions'),
    key: 'action',
    width: 200
  }
];

const fileColumns = [
  {
    title: t('virustotal.fileInfo'),
    dataIndex: 'fileName',
    key: 'fileName',
    ellipsis: true
  },
  {
    title: 'Size',
    dataIndex: 'fileSize',
    key: 'fileSize',
    width: 100
  },
  {
    title: t('virustotal.taskStatus'),
    dataIndex: 'status',
    key: 'status',
    width: 100
  },
  {
    title: t('virustotal.detectionRate'),
    dataIndex: 'detectionRate',
    key: 'detectionRate',
    width: 150
  },
  {
    title: t('virustotal.actions'),
    key: 'action',
    width: 100
  }
];

// 目录扫描进度
const directoryProgress = computed(() => {
  if (!currentTask.value || !currentTask.value.totalFiles) return 0;
  return Math.round((currentTask.value.completedFiles || 0) / currentTask.value.totalFiles * 100);
});

// 筛选后的结果
const filteredResults = computed(() => {
  if (!currentTask.value?.results) return [];
  if (resultFilter.value === 'all') return currentTask.value.results;
  return currentTask.value.results.filter((r: any) => r.category === resultFilter.value);
});

// 筛选后的文件检测结果
const filteredFileResults = computed(() => {
  if (!currentFile.value?.results) return [];
  if (fileResultFilter.value === 'all') return currentFile.value.results;
  return currentFile.value.results.filter((r: any) => r.category === fileResultFilter.value);
});

// 筛选后的文件列表
const filteredFiles = computed(() => {
  if (!taskFiles.value) return [];
  if (fileFilter.value === 'all') return taskFiles.value;

  return taskFiles.value.filter((f: any) => {
    if (fileFilter.value === 'malicious') {
      return f.status === 'completed' && f.stats?.malicious > 0;
    }
    if (fileFilter.value === 'suspicious') {
      return f.status === 'completed' && f.stats?.suspicious > 0 && f.stats?.malicious === 0;
    }
    if (fileFilter.value === 'safe') {
      return f.status === 'completed' && f.stats?.malicious === 0 && f.stats?.suspicious === 0;
    }
    if (fileFilter.value === 'pending') {
      return f.status !== 'completed' && f.status !== 'failed';
    }
    return true;
  });
});

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (!bytes) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

// 获取分类颜色
const getCategoryColor = (category: string): string => {
  const colorMap: Record<string, string> = {
    'malicious': '#ff4d4f',
    'suspicious': '#faad14',
    'harmless': '#52c41a',
    'undetected': '#1890ff',
    'type-unsupported': '#8c8c8c'
  };
  return colorMap[category] || '#d9d9d9';
};

// 获取分类文本
const getCategoryText = (category: string): string => {
  const textMap: Record<string, string> = {
    'malicious': t('virustotal.malicious'),
    'suspicious': t('virustotal.suspicious'),
    'harmless': t('virustotal.harmless'),
    'undetected': t('virustotal.undetected'),
    'type-unsupported': t('virustotal.unsupported')
  };
  return textMap[category] || category;
};

// 获取文件状态颜色
const getFileStatusColor = (status: string): string => {
  const colorMap: Record<string, string> = {
    pending: 'default',
    queued: 'processing',
    completed: 'success',
    failed: 'error',
    timeout: 'warning'
  };
  return colorMap[status] || 'default';
};

// 获取文件状态文本
const getFileStatusText = (status: string): string => {
  const textMap: Record<string, string> = {
    pending: t('virustotal.statusPending'),
    queued: t('virustotal.statusQueued'),
    completed: t('virustotal.statusCompleted'),
    failed: t('virustotal.statusFailed'),
    timeout: t('virustotal.timeout')
  };
  return textMap[status] || status;
};

const openCreateTaskModal = () => {
  if (!settingsStore.vtApiKey) {
    message.warning(t('settings.apiKeyNotConfiguredDesc'));
    return;
  }
  formState.filePath = '';
  formState.directoryPath = '';
  formState.scanType = 'file';
  createModalVisible.value = true;
};

const chooseFile = async () => {
  try {
    formState.filePath = await VTService.OpenFileDialog();
  } catch (error) {
    console.error("Failed to open file dialog:", error);
    message.error(t('virustotal.selectFileFailed'));
  }
};

const chooseDirectory = async () => {
  try {
    formState.directoryPath = await VTService.OpenFileDialogs();
  } catch (error) {
    console.error("Failed to open directory dialog:", error);
    message.error(t('virustotal.selectDirectoryFailed'));
  }
};

const createTask = async () => {
  if (!settingsStore.vtApiKey) {
    message.error(t('virustotal.pleaseConfigApiKey'));
    return;
  }

  if (formState.scanType === 'file' && !formState.filePath) {
    message.error(t('virustotal.pleaseSelectFile'));
    return;
  }

  if (formState.scanType === 'directory' && !formState.directoryPath) {
    message.error(t('virustotal.pleaseSelectDirectory'));
    return;
  }

  await VTService.SetAPIKey(settingsStore.vtApiKey);

  createLoading.value = true;
  try {
    if (formState.scanType === 'file') {
      const task = await VTService.CreateScanTask(formState.filePath);
      if (task) {
        message.success(t('virustotal.taskCreated'));
        tasks.value.unshift(task);
        startPolling(task.id);
      }
    } else {
      const task = await VTService.CreateDirectoryScanTask(formState.directoryPath);
      if (task) {
        message.success(t('virustotal.directoryTaskCreated', { count: task.totalFiles }));
        tasks.value.unshift(task);
        startPolling(task.id);
      }
    }

    createModalVisible.value = false;
  } catch (error: any) {
    console.error('Failed to create task:', error);
    message.error(t('virustotal.taskCreateFailed') + ': ' + (error.message || ''));
  } finally {
    createLoading.value = false;
  }
};

const createCancel = () => {
  createModalVisible.value = false;
};

const loadTasks = async () => {
  loading.value = true;
  try {
    const result = await VTService.GetAllTasks();
    tasks.value = result || [];

    tasks.value.forEach(task => {
      if (task.status !== 'completed' && task.status !== 'failed') {
        startPolling(task.id);
      }
    });
  } catch (error) {
    console.error('Failed to load tasks:', error);
    message.error(t('virustotal.loadTaskListFailed'));
  } finally {
    loading.value = false;
  }
};

const viewTaskDetail = async (task: any) => {
  currentTask.value = task;
  resultFilter.value = 'all';
  fileFilter.value = 'all';
  detailModalVisible.value = true;

  // 加载文件列表 - 统一处理单文件和目录任务
  await loadTaskFiles(task.id);

  if (task.status !== 'completed' && task.taskType === 'single') {
    await refreshTaskStatus(task);
  }
};

const loadTaskFiles = async (taskId: string) => {
  filesLoading.value = true;
  try {
    const files = await VTService.GetTaskFiles(taskId);
    taskFiles.value = files || [];
  } catch (error) {
    console.error('Failed to load task files:', error);
    message.error(t('virustotal.loadFileListFailed'));
  } finally {
    filesLoading.value = false;
  }
};

const viewFileDetail = async (file: any) => {
  try {
    const fileDetail = await VTService.GetFileDetail(file.id);
    currentFile.value = fileDetail;
    fileResultFilter.value = 'all';
    fileDetailModalVisible.value = true;
  } catch (error) {
    console.error('Failed to get file detail:', error);
    message.error(t('virustotal.getFileDetailFailed'));
  }
};

const refreshFileStatus = async (file: any) => {
  if (!settingsStore.vtApiKey) {
    message.warning(t('virustotal.pleaseConfigApiKey'));
    return;
  }

  file.refreshing = true;
  try {
    await VTService.SetAPIKey(settingsStore.vtApiKey);
    const updatedFile = await VTService.RefreshFileStatus(file.id);

    const index = taskFiles.value.findIndex(f => f.id === file.id);
    if (index !== -1) {
      taskFiles.value[index] = { ...taskFiles.value[index], ...updatedFile, refreshing: false };
    }

    // 更新目录任务统计
    if (currentTask.value) {
      const updatedTask = await VTService.GetTask(currentTask.value.id);
      if (updatedTask) {
        currentTask.value = { ...currentTask.value, ...updatedTask };
      }
    }
  } catch (error: any) {
    console.error('Failed to refresh file status:', error);
    message.error(t('virustotal.refreshFailed'));
  } finally {
    file.refreshing = false;
  }
};

const refreshTaskStatus = async (task: any) => {
  if (!settingsStore.vtApiKey) {
    message.warning(t('virustotal.pleaseConfigApiKey'));
    return;
  }

  task.refreshing = true;
  try {
    await VTService.SetAPIKey(settingsStore.vtApiKey);

    if (task.taskType === 'single') {
      const updatedTask = await VTService.CheckAnalysisStatus(task.id);

      const index = tasks.value.findIndex(t => t.id === task.id);
      if (index !== -1) {
        tasks.value[index] = { ...tasks.value[index], ...updatedTask, refreshing: false };
      }

      if (updatedTask?.status === 'completed') {
        stopPolling(task.id);
        message.success(t('virustotal.statusCompleted'));
      }

      if (currentTask.value?.id === task.id) {
        currentTask.value = tasks.value[index];
      }
    } else {
      // 目录任务刷新
      const updatedTask = await VTService.GetTask(task.id);
      if (updatedTask) {
        const index = tasks.value.findIndex(t => t.id === task.id);
        if (index !== -1) {
          tasks.value[index] = { ...tasks.value[index], ...updatedTask, refreshing: false };
        }
        if (currentTask.value?.id === task.id) {
          currentTask.value = { ...currentTask.value, ...updatedTask };
        }
        if (updatedTask.status === 'completed') {
          stopPolling(task.id);
          message.success(t('virustotal.statusCompleted'));
        }
      }
    }
  } catch (error: any) {
    console.error('Failed to refresh task status:', error);
    message.error(t('virustotal.refreshFailed') + ': ' + (error.message || ''));
  } finally {
    task.refreshing = false;
  }
};

const deleteTask = async (id: string) => {
  try {
    await VTService.DeleteTask(id);
    tasks.value = tasks.value.filter(t => t.id !== id);
    stopPolling(id);
    message.success(t('virustotal.deleteSuccess'));
  } catch (error: any) {
    console.error('Failed to delete task:', error);
    message.error(t('virustotal.deleteFailed') + ': ' + (error.message || ''));
  }
};

// 轮询管理
const pollingTasks = new Set<string>();

const startPolling = (taskId: string) => {
  if (pollingTasks.has(taskId)) return;
  pollingTasks.add(taskId);

  if (!pollingTimer) {
    pollingTimer = setInterval(() => {
      pollingTasks.forEach(id => {
        const task = tasks.value.find(t => t.id === id);
        if (task && task.status !== 'completed' && task.status !== 'failed') {
          refreshTaskStatus(task);
        } else {
          stopPolling(id);
        }
      });
    }, 10000);
  }
};

const stopPolling = (taskId: string) => {
  pollingTasks.delete(taskId);

  if (pollingTasks.size === 0 && pollingTimer) {
    clearInterval(pollingTimer);
    pollingTimer = null;
  }
};

const getStatusColor = (status: string) => {
  const colorMap: Record<string, string> = {
    pending: 'processing',
    queued: 'processing',
    'in-progress': 'warning',
    scanning: 'warning',
    completed: 'success',
    failed: 'error'
  };
  return colorMap[status] || 'default';
};

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    pending: t('virustotal.statusPending'),
    queued: t('virustotal.statusQueued'),
    'in-progress': t('virustotal.statusAnalyzing'),
    scanning: t('virustotal.statusScanning'),
    completed: t('virustotal.statusCompleted'),
    failed: t('virustotal.statusFailed')
  };
  return textMap[status] || status;
};

// 获取检出情况文本 - 统一格式: 检出文件数/总文件数
const getDetectionText = (task: any): string => {
  if (task.taskType === 'directory') {
    // 目录任务: 恶意文件数/总文件数
    const malicious = task.maliciousFiles || 0;
    const total = task.totalFiles || 0;
    return `${malicious}/${total}`;
  } else {
    // 单文件任务: 是否有恶意检出 (1/1 或 0/1)
    const hasMalicious = (task.stats?.malicious || 0) > 0 ? 1 : 0;
    return `${hasMalicious}/1`;
  }
};

// 获取检出情况样式类
const getDetectionClass = (task: any): string => {
  if (task.taskType === 'directory') {
    const malicious = task.maliciousFiles || 0;
    if (malicious > 0) return 'detection-malicious';
    return 'detection-safe';
  } else {
    // 单文件任务: 根据是否有恶意检出判断
    const hasMalicious = (task.stats?.malicious || 0) > 0;
    if (hasMalicious) return 'detection-malicious';
    return 'detection-safe';
  }
};

// API 配额
interface QuotaItem { allowed: number; used: number }
interface APIQuota { hourly: QuotaItem; daily: QuotaItem; monthly: QuotaItem }
const apiQuota = ref<APIQuota | null>(null)
const quotaLoading = ref(false)

const quotaDisplayList = computed(() => {
  if (!apiQuota.value) return []
  return [
    { label: t('virustotal.quota.hourly'), quota: apiQuota.value.hourly },
    { label: t('virustotal.quota.daily'),  quota: apiQuota.value.daily  },
    { label: t('virustotal.quota.monthly'), quota: apiQuota.value.monthly },
  ]
})

const quotaPercent = (q: QuotaItem) => q.allowed > 0 ? Math.min(100, Math.round(q.used / q.allowed * 100)) : 0
const quotaColor = (q: QuotaItem) => {
  const p = quotaPercent(q)
  if (p >= 90) return '#ff4d4f'
  if (p >= 70) return '#faad14'
  return '#52c41a'
}
const formatQuotaNum = (n: number) => n >= 1000000 ? (n / 1000000).toFixed(1) + 'M' : n >= 1000 ? (n / 1000).toFixed(0) + 'K' : String(n)

const loadQuota = async () => {
  if (!settingsStore.vtApiKey) return
  quotaLoading.value = true
  try {
    const q = await VTService.GetAPIQuota()
    apiQuota.value = q as APIQuota
    quotaLoading.value = false
  } catch (e) {
    console.warn('Failed to load API quota:', e)
  }
}

onMounted(() => {
  loadTasks();
  loadQuota();
});

onUnmounted(() => {
  if (pollingTimer) {
    clearInterval(pollingTimer);
    pollingTimer = null;
  }
  pollingTasks.clear();
});
</script>

<style scoped>
.task-detail {
  max-height: 70vh;
  overflow-y: auto;
}

.mono-text {
  font-family: 'Consolas', 'Monaco', monospace;
}

.hash-value {
  font-size: 12px;
  word-break: break-all;
}

.ml-1 {
  margin-left: 4px;
}

.ml-2 {
  margin-left: 8px;
}

.text-gray-400 {
  color: #9ca3af;
}

.text-gray-500 {
  color: #6b7280;
}

.text-xs {
  font-size: 12px;
}

/* 统计信息样式 */
.stats-container {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 20px;
  border-radius: 8px;
  min-width: 80px;
}

.stat-count {
  font-size: 24px;
  font-weight: bold;
}

.stat-label {
  font-size: 12px;
  margin-top: 4px;
}

.stat-malicious {
  background-color: #fff1f0;
  color: #cf132d;
}

.stat-suspicious {
  background-color: #fff7e6;
  color: #d46b08;
}

.stat-harmless, .stat-safe {
  background-color: #f6ffed;
  color: #389e0d;
}

.stat-undetected {
  background-color: #e6f7ff;
  color: #1890ff;
}

.stat-unsupported {
  background-color: #f5f5f5;
  color: #595959;
}

.stat-timeout {
  background-color: #fff0f6;
  color: #c41d7f;
}

.stat-failure {
  background-color: #fff1f0;
  color: #a8071a;
}

.stat-pending {
  background-color: #f0f0f0;
  color: #8c8c8c;
}

/* 检出情况样式 */
.detection-malicious {
  color: #cf132d;
  font-weight: 600;
}

.detection-safe {
  color: #389e0d;
  font-weight: 600;
}

/* 筛选区域 */
.filter-section {
  margin-bottom: 16px;
}

.result-count {
  font-size: 12px;
  color: #8c8c8c;
  font-weight: normal;
  margin-left: 8px;
}

/* 检测结果容器 */
.results-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
  padding: 8px;
  background-color: #fafafa;
  border-radius: 8px;
}

.engine-tag {
  cursor: pointer;
  font-size: 12px;
  margin: 2px;
}

.engine-tag:hover {
  opacity: 0.8;
  transform: scale(1.05);
  transition: all 0.2s;
}

/* Tooltip 样式 */
.result-tooltip {
  font-size: 12px;
  line-height: 1.6;
}

.result-tooltip div {
  margin: 2px 0;
}

/* quota 信息条 */
.quota-header {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #f0f5ff 0%, #e6f7ff 100%);
  border-bottom: 1px solid #d6e4ff;
}

.quota-header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.quota-header-title {
  font-size: 13px;
  font-weight: 600;
  color: #2f54eb;
  display: flex;
  align-items: center;
  white-space: nowrap;
}

.quota-header-items {
  display: flex;
  flex: 1;
  gap: 24px;
}

.quota-stat-item {
  flex: 1;
  min-width: 120px;
}

.quota-stat-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 4px;
}

.quota-stat-label {
  font-size: 12px;
  color: #595959;
  font-weight: 500;
}

.quota-stat-nums {
  font-size: 12px;
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.quota-stat-sep {
  color: #bfbfbf;
  font-size: 11px;
}

.quota-stat-total {
  color: #8c8c8c;
}

.quota-progress {
  margin: 0 !important;
}

/* 任务列表 header */
.task-list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.task-list-title {
  font-size: 15px;
  font-weight: 600;
  color: rgba(0, 0, 0, 0.85);
}

.quota-table {
  padding: 0 8px;
}

/* 深色主题 */
[data-theme="dark"] .quota-header {
  background: linear-gradient(135deg, #1a1f3a 0%, #111d2c 100%);
  border-bottom-color: #1e3a5f;
}

[data-theme="dark"] .quota-header-title {
  color: #69b1ff;
}

[data-theme="dark"] .quota-stat-label {
  color: #a6a6a6;
}

[data-theme="dark"] .task-list-header {
  border-bottom-color: #303030;
}

[data-theme="dark"] .task-list-title {
  color: rgba(255, 255, 255, 0.85);
}

[data-theme="dark"] .results-container {
  background-color: #1f1f1f;
}

[data-theme="dark"] .stat-malicious {
  background-color: #2a1215;
}

[data-theme="dark"] .stat-suspicious {
  background-color: #2b2117;
}

[data-theme="dark"] .stat-harmless,
[data-theme="dark"] .stat-safe {
  background-color: #162312;
}

[data-theme="dark"] .stat-undetected {
  background-color: #111d2c;
}

[data-theme="dark"] .stat-unsupported {
  background-color: #262626;
}

[data-theme="dark"] .stat-timeout {
  background-color: #2a1a24;
}

[data-theme="dark"] .stat-failure {
  background-color: #2a1215;
}

[data-theme="dark"] .stat-pending {
  background-color: #262626;
}
</style>
