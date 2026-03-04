<template>
  <div class="vt p-6">
    <a-card>
      <template #title>
        <div class="flex items-center justify-between">
          <div>
            任务列表
          </div>
          <a-button type="primary" class="ml-2" :icon="h(PlusOutlined)" @click="openCreateTaskModal">创建任务</a-button>
        </div>
      </template>
      <a-table :dataSource="tasks" :columns="columns" rowKey="id" :loading="loading" :pagination="{ pageSize: 10 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'taskType'">
            <a-tag :color="record.taskType === 'directory' ? 'blue' : 'green'">
              {{ record.taskType === 'directory' ? '目录扫描' : '单文件' }}
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
                查看详情
              </a-button>
              <a-button
                type="link"
                size="small"
                @click="refreshTaskStatus(record)"
                :loading="record.refreshing"
                v-if="record.status !== 'completed' && record.status !== 'failed'"
              >
                刷新状态
              </a-button>
              <a-popconfirm
                title="确定要删除此任务吗？"
                @confirm="deleteTask(record.id)"
                okText="确定"
                cancelText="取消"
              >
                <a-button type="link" size="small" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 创建任务弹窗 -->
    <a-modal
      v-model:open="createModalVisible"
      title="创建扫描任务"
      centered
      :maskClosable="false"
      @ok="createTask"
      @cancel="createCancel"
      :confirmLoading="createLoading"
    >
      <a-alert
        v-if="!settingsStore.vtApiKey"
        message="未配置 API Key"
        description="请先在全局设置中配置 VirusTotal API Key"
        type="warning"
        show-icon
        class="mb-4"
      />

      <a-form :model="formState" layout="vertical">
        <a-form-item label="扫描类型" required>
          <a-radio-group v-model:value="formState.scanType">
            <a-radio value="file">单个文件</a-radio>
            <a-radio value="directory">批量目录</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item
          label="文件路径"
          v-if="formState.scanType === 'file'"
          required
        >
          <a-input-group compact>
            <a-input v-model:value="formState.filePath" readonly style="width: calc(100% - 100px)" />
            <a-button type="primary" @click="chooseFile">
              选择文件
            </a-button>
          </a-input-group>
        </a-form-item>

        <a-form-item
          label="目录路径"
          v-if="formState.scanType === 'directory'"
          required
        >
          <a-input-group compact>
            <a-input v-model:value="formState.directoryPath" readonly style="width: calc(100% - 100px)" />
            <a-button type="primary" @click="chooseDirectory">
              选择目录
            </a-button>
          </a-input-group>
          <div class="text-gray-500 text-xs mt-1">将递归扫描目录下所有文件</div>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 任务详情弹窗 - 统一显示文件列表 -->
    <a-modal
      v-model:open="detailModalVisible"
      :title="currentTask?.taskType === 'directory' ? '目录扫描结果' : '扫描结果'"
      :width="currentTask?.taskType === 'directory' ? '1200px' : '1000px'"
      :footer="null"
    >
      <div v-if="currentTask" class="task-detail">
        <!-- 基本信息 -->
        <a-descriptions :column="currentTask?.taskType === 'directory' ? 3 : 2" bordered size="small">
          <a-descriptions-item label="任务类型">
            <a-tag :color="currentTask.taskType === 'directory' ? 'blue' : 'green'">
              {{ currentTask.taskType === 'directory' ? '目录扫描' : '单文件' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="名称">
            <a-tooltip :title="currentTask.filePath">
              <span class="mono-text">{{ currentTask.fileName || '未知' }}</span>
            </a-tooltip>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentTask.status)">
              {{ getStatusText(currentTask.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="文件总数">{{ currentTask.taskType === 'directory' ? (currentTask.totalFiles || 0) : 1 }}</a-descriptions-item>
          <a-descriptions-item label="扫描进度" v-if="currentTask.taskType === 'directory'">
            {{ currentTask.completedFiles || 0 }} / {{ currentTask.totalFiles || 0 }}
            <a-progress
              :percent="directoryProgress"
              size="small"
              style="width: 100px; margin-left: 8px"
            />
          </a-descriptions-item>
          <a-descriptions-item label="扫描时间">{{ currentTask.scanTime || '-' }}</a-descriptions-item>
        </a-descriptions>

        <!-- 统计概览 - 目录任务 -->
        <template v-if="currentTask.taskType === 'directory'">
          <a-divider orientation="left">扫描概览</a-divider>
          <div class="stats-container">
            <div class="stat-item stat-malicious">
              <div class="stat-count">{{ currentTask.maliciousFiles || 0 }}</div>
              <div class="stat-label">恶意文件</div>
            </div>
            <div class="stat-item stat-suspicious">
              <div class="stat-count">{{ currentTask.suspiciousFiles || 0 }}</div>
              <div class="stat-label">可疑文件</div>
            </div>
            <div class="stat-item stat-safe">
              <div class="stat-count">{{ (currentTask.completedFiles || 0) - (currentTask.maliciousFiles || 0) - (currentTask.suspiciousFiles || 0) }}</div>
              <div class="stat-label">安全文件</div>
            </div>
            <div class="stat-item stat-pending">
              <div class="stat-count">{{ (currentTask.totalFiles || 0) - (currentTask.completedFiles || 0) }}</div>
              <div class="stat-label">待扫描</div>
            </div>
          </div>
        </template>

        <!-- 统计概览 - 单文件任务 -->
        <template v-else>
          <a-divider orientation="left">扫描概览</a-divider>
          <div class="stats-container">
            <div class="stat-item stat-malicious">
              <div class="stat-count">{{ (currentTask.stats?.malicious || 0) > 0 ? 1 : 0 }}</div>
              <div class="stat-label">恶意文件</div>
            </div>
            <div class="stat-item stat-suspicious">
              <div class="stat-count">{{ (currentTask.stats?.suspicious || 0) > 0 && (currentTask.stats?.malicious || 0) === 0 ? 1 : 0 }}</div>
              <div class="stat-label">可疑文件</div>
            </div>
            <div class="stat-item stat-safe">
              <div class="stat-count">{{ (currentTask.stats?.malicious || 0) === 0 && (currentTask.stats?.suspicious || 0) === 0 && currentTask.status === 'completed' ? 1 : 0 }}</div>
              <div class="stat-label">安全文件</div>
            </div>
            <div class="stat-item stat-pending">
              <div class="stat-count">{{ currentTask.status !== 'completed' ? 1 : 0 }}</div>
              <div class="stat-label">待扫描</div>
            </div>
          </div>
        </template>

        <!-- 文件列表 -->
        <a-divider orientation="left">
          文件列表
          <span class="result-count">(共 {{ taskFiles.length }} 个文件)</span>
        </a-divider>

        <!-- 文件筛选 - 仅目录任务显示 -->
        <div class="filter-section" v-if="currentTask.taskType === 'directory'">
          <a-radio-group v-model:value="fileFilter" button-style="solid" size="small">
            <a-radio-button value="all">全部</a-radio-button>
            <a-radio-button value="malicious">恶意</a-radio-button>
            <a-radio-button value="suspicious">可疑</a-radio-button>
            <a-radio-button value="safe">安全</a-radio-button>
            <a-radio-button value="pending">扫描中</a-radio-button>
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
                查看详情
              </a-button>
              <a-button
                type="link"
                size="small"
                @click="refreshFileStatus(record)"
                :loading="record.refreshing"
                v-else-if="record.status !== 'failed'"
              >
                刷新
              </a-button>
            </template>
          </template>
        </a-table>
      </div>
    </a-modal>

    <!-- 文件详情弹窗 -->
    <a-modal
      v-model:open="fileDetailModalVisible"
      title="文件扫描详情"
      width="900px"
      :footer="null"
    >
      <div v-if="currentFile" class="task-detail">
        <!-- 基本信息 -->
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="文件名">{{ currentFile.fileName || '未知' }}</a-descriptions-item>
          <a-descriptions-item label="文件大小">{{ formatFileSize(currentFile.fileSize) }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getFileStatusColor(currentFile.status)">
              {{ getFileStatusText(currentFile.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="检出率">
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
            <span class="mono-text hash-value">{{ currentFile.md5 || '等待分析...' }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="SHA256" :span="2">
            <span class="mono-text hash-value">{{ currentFile.sha256 || '等待分析...' }}</span>
          </a-descriptions-item>
        </a-descriptions>

        <!-- 统计信息 -->
        <a-divider orientation="left">检测统计</a-divider>
        <div class="stats-container" v-if="currentFile.stats">
          <div class="stat-item stat-malicious">
            <div class="stat-count">{{ currentFile.stats.malicious || 0 }}</div>
            <div class="stat-label">恶意</div>
          </div>
          <div class="stat-item stat-suspicious">
            <div class="stat-count">{{ currentFile.stats.suspicious || 0 }}</div>
            <div class="stat-label">可疑</div>
          </div>
          <div class="stat-item stat-harmless">
            <div class="stat-count">{{ currentFile.stats.harmless || 0 }}</div>
            <div class="stat-label">安全</div>
          </div>
          <div class="stat-item stat-undetected">
            <div class="stat-count">{{ currentFile.stats.undetected || 0 }}</div>
            <div class="stat-label">未检出</div>
          </div>
          <div class="stat-item stat-unsupported">
            <div class="stat-count">{{ currentFile.stats.typeUnsupported || 0 }}</div>
            <div class="stat-label">不支持</div>
          </div>
          <div class="stat-item stat-timeout">
            <div class="stat-count">{{ currentFile.stats.timeout || 0 }}</div>
            <div class="stat-label">超时</div>
          </div>
          <div class="stat-item stat-failure">
            <div class="stat-count">{{ currentFile.stats.failure || 0 }}</div>
            <div class="stat-label">失败</div>
          </div>
        </div>

        <!-- 检测结果 -->
        <a-divider orientation="left">
          检测结果
          <span class="result-count" v-if="currentFile.results">
            (共 {{ currentFile.results.length }} 个引擎)
          </span>
        </a-divider>

        <!-- 分类筛选 -->
        <div class="filter-section" v-if="currentFile.results && currentFile.results.length > 0">
          <a-radio-group v-model:value="fileResultFilter" button-style="solid" size="small">
            <a-radio-button value="all">全部</a-radio-button>
            <a-radio-button value="malicious">恶意</a-radio-button>
            <a-radio-button value="suspicious">可疑</a-radio-button>
            <a-radio-button value="harmless">安全</a-radio-button>
            <a-radio-button value="undetected">未检出</a-radio-button>
            <a-radio-button value="type-unsupported">不支持</a-radio-button>
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
                <div><strong>引擎:</strong> {{ result.engine }}</div>
                <div><strong>分类:</strong> {{ getCategoryText(result.category) }}</div>
                <div><strong>结果:</strong> {{ result.result || '无' }}</div>
                <div><strong>方法:</strong> {{ result.method || '未知' }}</div>
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
        <a-empty v-else description="暂无检测结果" />
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import {
  PlusOutlined,
  FileFilled,
  FolderFilled,
} from '@ant-design/icons-vue';
import { h, onMounted, onUnmounted, ref, reactive, computed } from 'vue';
import { message } from 'ant-design-vue';
import { VTService } from "../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/vt";
import { useSettingsStore } from '@/stores/settings';

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
    title: '任务类型',
    dataIndex: 'taskType',
    key: 'taskType',
    width: 100
  },
  {
    title: '总文件数',
    dataIndex: 'totalFiles',
    key: 'totalFiles',
    width: 100
  },
  {
    title: '任务状态',
    dataIndex: 'status',
    key: 'status',
    width: 130
  },
  {
    title: '检出情况',
    dataIndex: 'detectionRate',
    key: 'detectionRate',
    width: 120
  },
  {
    title: '扫描时间',
    dataIndex: 'scanTime',
    key: 'scanTime',
    width: 180
  },
  {
    title: '操作',
    key: 'action',
    width: 200
  }
];

const fileColumns = [
  {
    title: '文件名',
    dataIndex: 'fileName',
    key: 'fileName',
    ellipsis: true
  },
  {
    title: '大小',
    dataIndex: 'fileSize',
    key: 'fileSize',
    width: 100
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100
  },
  {
    title: '检出率',
    dataIndex: 'detectionRate',
    key: 'detectionRate',
    width: 150
  },
  {
    title: '操作',
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
    'malicious': '恶意',
    'suspicious': '可疑',
    'harmless': '安全',
    'undetected': '未检出',
    'type-unsupported': '不支持'
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
    pending: '等待中',
    queued: '排队中',
    completed: '已完成',
    failed: '失败',
    timeout: '超时'
  };
  return textMap[status] || status;
};

const openCreateTaskModal = () => {
  if (!settingsStore.vtApiKey) {
    message.warning('请先在全局设置中配置 VirusTotal API Key');
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
    message.error('选择文件失败');
  }
};

const chooseDirectory = async () => {
  try {
    formState.directoryPath = await VTService.OpenFileDialogs();
  } catch (error) {
    console.error("Failed to open directory dialog:", error);
    message.error('选择目录失败');
  }
};

const createTask = async () => {
  if (!settingsStore.vtApiKey) {
    message.error('请先配置 API Key');
    return;
  }

  if (formState.scanType === 'file' && !formState.filePath) {
    message.error('请选择文件');
    return;
  }

  if (formState.scanType === 'directory' && !formState.directoryPath) {
    message.error('请选择目录');
    return;
  }

  await VTService.SetAPIKey(settingsStore.vtApiKey);

  createLoading.value = true;
  try {
    if (formState.scanType === 'file') {
      const task = await VTService.CreateScanTask(formState.filePath);
      if (task) {
        message.success('任务创建成功，正在上传文件...');
        tasks.value.unshift(task);
        startPolling(task.id);
      }
    } else {
      const task = await VTService.CreateDirectoryScanTask(formState.directoryPath);
      if (task) {
        message.success(`目录扫描任务创建成功，共 ${task.totalFiles} 个文件`);
        tasks.value.unshift(task);
        startPolling(task.id);
      }
    }

    createModalVisible.value = false;
  } catch (error: any) {
    console.error('Failed to create task:', error);
    message.error('任务创建失败: ' + (error.message || '未知错误'));
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
    message.error('加载任务列表失败');
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
    message.error('加载文件列表失败');
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
    message.error('获取文件详情失败');
  }
};

const refreshFileStatus = async (file: any) => {
  if (!settingsStore.vtApiKey) {
    message.warning('请先配置 API Key');
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
    message.error('刷新状态失败');
  } finally {
    file.refreshing = false;
  }
};

const refreshTaskStatus = async (task: any) => {
  if (!settingsStore.vtApiKey) {
    message.warning('请先配置 API Key');
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
        message.success('分析完成');
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
          message.success('目录扫描完成');
        }
      }
    }
  } catch (error: any) {
    console.error('Failed to refresh task status:', error);
    message.error('刷新状态失败: ' + (error.message || '未知错误'));
  } finally {
    task.refreshing = false;
  }
};

const deleteTask = async (id: string) => {
  try {
    await VTService.DeleteTask(id);
    tasks.value = tasks.value.filter(t => t.id !== id);
    stopPolling(id);
    message.success('任务已删除');
  } catch (error: any) {
    console.error('Failed to delete task:', error);
    message.error('删除失败: ' + (error.message || '未知错误'));
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
    pending: '等待中',
    queued: '排队中',
    'in-progress': '分析中',
    scanning: '扫描中',
    completed: '已完成',
    failed: '失败'
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

onMounted(() => {
  loadTasks();
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

/* 深色主题 */
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
