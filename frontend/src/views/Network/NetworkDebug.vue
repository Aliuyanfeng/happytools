<template>
  <div class="network-debug p-6">
    <a-tabs v-model:activeKey="activeTab" type="card">
      <!-- 文件传输协议 -->
      <a-tab-pane key="file">
        <template #tab>
          <span class="flex items-center">
            <FolderOutlined />
            文件传输协议
          </span>
        </template>
        <div class="protocol-section">
          <a-radio-group v-model:value="fileProtocol" button-style="solid" class="mb-4">
            <a-radio-button value="ftp">FTP</a-radio-button>
            <a-radio-button value="sftp">SFTP</a-radio-button>
            <a-radio-button value="smb">SMB</a-radio-button>
            <a-radio-button value="s3">S3</a-radio-button>
          </a-radio-group>

          <!-- FTP/SFTP/SMB 配置 -->
          <div v-if="fileProtocol !== 's3'" class="config-form mb-4">
            <a-form layout="inline">
              <a-form-item label="主机">
                <a-input v-model:value="fileConfig.host" placeholder="例如: 192.168.1.100" style="width: 200px" />
              </a-form-item>
              <a-form-item label="端口">
                <a-input-number v-model:value="fileConfig.port" :min="1" :max="65535" style="width: 100px" />
              </a-form-item>
              <a-form-item label="用户名">
                <a-input v-model:value="fileConfig.username" placeholder="用户名" style="width: 150px" />
              </a-form-item>
              <a-form-item label="密码">
                <a-input-password v-model:value="fileConfig.password" placeholder="密码" style="width: 150px" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="testFileConnection" :loading="testingConnection">
                  测试连接
                </a-button>
              </a-form-item>
            </a-form>
          </div>

          <!-- S3 配置 -->
          <div v-else class="config-form mb-4">
            <a-form layout="inline">
              <a-form-item label="端点">
                <a-input v-model:value="s3Config.endpoint" placeholder="例如: https://s3.amazonaws.com" style="width: 300px" />
              </a-form-item>
              <a-form-item label="Access Key">
                <a-input v-model:value="s3Config.accessKeyId" placeholder="Access Key ID" style="width: 200px" />
              </a-form-item>
              <a-form-item label="Secret Key">
                <a-input-password v-model:value="s3Config.secretAccessKey" placeholder="Secret Access Key" style="width: 200px" />
              </a-form-item>
              <a-form-item label="Bucket">
                <a-input v-model:value="s3Config.bucket" placeholder="存储桶名称" style="width: 150px" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="testS3Connection" :loading="testingConnection">
                  测试连接
                </a-button>
              </a-form-item>
            </a-form>
          </div>

          <!-- 连接状态 -->
          <a-alert v-if="connectionResult" :type="connectionResult.success ? 'success' : 'error'" class="mb-4" show-icon>
            <template #message>
              {{ connectionResult.message }}
              <span v-if="connectionResult.success"> (延迟: {{ connectionResult.latency }}ms)</span>
            </template>
          </a-alert>

          <!-- 操作按钮 -->
          <a-space class="mb-4">
            <a-button @click="selectLocalFile">选择本地文件</a-button>
            <a-button @click="selectLocalDir">选择本地目录</a-button>
            <a-button type="primary" @click="uploadFile" :disabled="!localFilePath">上传文件</a-button>
            <a-button type="primary" @click="downloadFile">下载文件</a-button>
            <a-button @click="listFiles" :disabled="!localDirPath">列出文件</a-button>
            <a-button @click="countFiles" :disabled="!localDirPath">统计文件数</a-button>
          </a-space>

          <!-- 本地文件路径 -->
          <div v-if="localFilePath" class="mb-4">
            <a-tag color="blue">本地文件: {{ localFilePath }}</a-tag>
          </div>

          <!-- 本地目录路径 -->
          <div v-if="localDirPath" class="mb-4">
            <a-tag color="green">本地目录: {{ localDirPath }}</a-tag>
          </div>

          <!-- 文件列表 -->
          <a-table v-if="fileList.length > 0" :dataSource="fileList" :columns="fileColumns" size="small" class="mb-4" rowKey="path">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'size'">
                {{ formatFileSize(record.size) }}
              </template>
              <template v-if="column.key === 'isDir'">
                <a-tag :color="record.isDir ? 'blue' : 'green'">
                  {{ record.isDir ? '目录' : '文件' }}
                </a-tag>
              </template>
            </template>
          </a-table>

          <!-- 文件数量统计 -->
          <a-alert v-if="fileCount !== null" type="info" class="mb-4" show-icon>
            <template #message>
              文件数量: {{ fileCount }}
            </template>
          </a-alert>
        </div>
      </a-tab-pane>

      <!-- TCP/IP协议 -->
      <a-tab-pane key="tcpip">
        <template #tab>
          <span class="flex items-center">
            <ApiOutlined />
            TCP/IP协议
          </span>
        </template>
        <div class="protocol-section">
          <a-radio-group v-model:value="tcpProtocol" button-style="solid" class="mb-4">
            <a-radio-button value="tcp">TCP</a-radio-button>
            <a-radio-button value="udp">UDP</a-radio-button>
          </a-radio-group>

          <!-- TCP/UDP 配置 -->
          <div class="config-form mb-4">
            <a-form layout="inline">
              <a-form-item label="主机">
                <a-input v-model:value="tcpConfig.host" placeholder="例如: 127.0.0.1" style="width: 200px" />
              </a-form-item>
              <a-form-item label="端口">
                <a-input-number v-model:value="tcpConfig.port" :min="1" :max="65535" style="width: 100px" />
              </a-form-item>
              <a-form-item v-if="tcpProtocol === 'tcp'">
                <a-button type="primary" @click="connectTCP" :loading="connecting">
                  {{ tcpConnected ? '断开连接' : '连接' }}
                </a-button>
              </a-form-item>
              <a-form-item>
                <a-button @click="testTCPConnection" :loading="testingConnection">
                  测试连接
                </a-button>
              </a-form-item>
            </a-form>
          </div>

          <!-- 连接状态 -->
          <a-alert v-if="tcpConnectionStatus" :type="tcpConnectionStatus.isConnected ? 'success' : 'info'" class="mb-4" show-icon>
            <template #message>
              {{ tcpConnectionStatus.isConnected ? '已连接' : '未连接' }}
              <span v-if="tcpConnectionStatus.isConnected">
                (本地: {{ tcpConnectionStatus.localAddr }} → 远程: {{ tcpConnectionStatus.remoteAddr }})
              </span>
            </template>
          </a-alert>

          <!-- 数据发送/接收 -->
          <a-card title="数据调试" class="mb-4">
            <a-form layout="vertical">
              <a-form-item label="发送数据">
                <a-textarea v-model:value="sendData" :rows="4" placeholder="输入要发送的数据" />
              </a-form-item>
              <a-form-item>
                <a-checkbox v-model:checked="sendAsHex">以十六进制发送</a-checkbox>
                <a-button type="primary" class="ml-4" @click="sendDataToServer" :disabled="tcpProtocol === 'tcp' && !tcpConnected">
                  发送
                </a-button>
                <a-button class="ml-2" @click="receiveDataFromServer" :disabled="tcpProtocol === 'tcp' && !tcpConnected">
                  接收
                </a-button>
                <a-button class="ml-2" @click="convertToHex">
                  转HEX
                </a-button>
                <a-button class="ml-2" @click="convertFromHex">
                  HEX转文本
                </a-button>
              </a-form-item>
              <a-form-item label="接收数据">
                <a-textarea v-model:value="receivedData" :rows="4" readonly placeholder="接收到的数据将显示在这里" />
              </a-form-item>
              <a-form-item v-if="receivedHexData" label="十六进制数据">
                <a-input v-model:value="receivedHexData" readonly />
              </a-form-item>
            </a-form>
          </a-card>
        </div>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { message } from 'ant-design-vue';
import { FolderOutlined, ApiOutlined } from '@ant-design/icons-vue';
import {
  FileTransferService,
  TCPUDPService,
  type FTPConfig,
  type SFTPConfig,
  type SMBConfig,
  type S3Config,
  type TCPUDPConfig
} from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/network';

const activeTab = ref('file');
const fileProtocol = ref('ftp');
const tcpProtocol = ref('tcp');

// 文件传输配置
const fileConfig = reactive({
  host: '',
  port: 21,
  username: '',
  password: ''
});

const s3Config = reactive({
  endpoint: '',
  accessKeyId: '',
  secretAccessKey: '',
  bucket: '',
  region: 'us-east-1',
  useSSL: true
});

// TCP/UDP配置
const tcpConfig = reactive({
  host: '127.0.0.1',
  port: 8080
});

// 状态
const testingConnection = ref(false);
const connecting = ref(false);
const connectionResult = ref<any>(null);
const tcpConnectionStatus = ref<any>(null);
const tcpConnected = ref(false);

// 文件操作
const localFilePath = ref('');
const localDirPath = ref('');
const fileList = ref<any[]>([]);
const fileCount = ref<number | null>(null);

// TCP/UDP数据
const sendData = ref('');
const receivedData = ref('');
const receivedHexData = ref('');
const sendAsHex = ref(false);

// 文件列表列
const fileColumns = [
  { title: '文件名', dataIndex: 'name', key: 'name' },
  { title: '路径', dataIndex: 'path', key: 'path', ellipsis: true },
  { title: '大小', dataIndex: 'size', key: 'size' },
  { title: '类型', dataIndex: 'isDir', key: 'isDir' }
];

// 测试文件传输连接
async function testFileConnection() {
  testingConnection.value = true;
  try {
    let result;
    if (fileProtocol.value === 'ftp') {
      const config: FTPConfig = {
        host: fileConfig.host,
        port: fileConfig.port,
        username: fileConfig.username,
        password: fileConfig.password,
        timeout: 10
      };
      result = await FileTransferService.TestFTPConnection(config);
    } else if (fileProtocol.value === 'sftp') {
      const config: SFTPConfig = {
        host: fileConfig.host,
        port: fileConfig.port || 22,
        username: fileConfig.username,
        password: fileConfig.password,
        keyFile: '',
        timeout: 10
      };
      result = await FileTransferService.TestSFTPConnection(config);
    } else if (fileProtocol.value === 'smb') {
      const config: SMBConfig = {
        host: fileConfig.host,
        port: fileConfig.port || 445,
        username: fileConfig.username,
        password: fileConfig.password,
        share: '',
        timeout: 10
      };
      result = await FileTransferService.TestSMBConnection(config);
    }
    
    connectionResult.value = result;
    if (result && result.success) {
      message.success('连接成功');
    } else if (result) {
      message.error(result.message);
    }
  } catch (error: any) {
    message.error(`测试失败: ${error.message}`);
  } finally {
    testingConnection.value = false;
  }
}

// 测试S3连接
async function testS3Connection() {
  testingConnection.value = true;
  try {
    const config: S3Config = {
      endpoint: s3Config.endpoint,
      accessKeyId: s3Config.accessKeyId,
      secretAccessKey: s3Config.secretAccessKey,
      region: s3Config.region,
      bucket: s3Config.bucket,
      useSSL: s3Config.useSSL
    };
    const result = await FileTransferService.TestS3Connection(config);
    connectionResult.value = result;
    if (result && result.success) {
      message.success('配置验证成功');
    } else if (result) {
      message.error(result.message);
    }
  } catch (error: any) {
    message.error(`测试失败: ${error.message}`);
  } finally {
    testingConnection.value = false;
  }
}

// 选择本地文件
async function selectLocalFile() {
  try {
    const path = await FileTransferService.OpenFileDialog();
    if (path) {
      localFilePath.value = path;
      message.success('已选择文件');
    }
  } catch (error: any) {
    message.error(`选择文件失败: ${error.message}`);
  }
}

// 选择本地目录
async function selectLocalDir() {
  try {
    const path = await FileTransferService.OpenDirectoryDialog();
    if (path) {
      localDirPath.value = path;
      message.success('已选择目录');
    }
  } catch (error: any) {
    message.error(`选择目录失败: ${error.message}`);
  }
}

// 上传文件
async function uploadFile() {
  if (!localFilePath.value) {
    message.warning('请先选择本地文件');
    return;
  }
  try {
    const result = await FileTransferService.UploadFile(
      localFilePath.value,
      '/remote/path/file.txt',
      fileProtocol.value
    );
    if (result && result.success) {
      message.success(`上传成功: ${result.speed}`);
    } else if (result) {
      message.error(result.message);
    }
  } catch (error: any) {
    message.error(`上传失败: ${error.message}`);
  }
}

// 下载文件
async function downloadFile() {
  try {
    const result = await FileTransferService.DownloadFile(
      '/remote/path/file.txt',
      'C:\\Downloads\\file.txt',
      fileProtocol.value
    );
    if (result && result.success) {
      message.success(`下载成功: ${result.speed}`);
    } else if (result) {
      message.error(result.message);
    }
  } catch (error: any) {
    message.error(`下载失败: ${error.message}`);
  }
}

// 列出文件
async function listFiles() {
  if (!localDirPath.value) {
    message.warning('请先选择本地目录');
    return;
  }
  try {
    const result = await FileTransferService.ListLocalFiles(localDirPath.value);
    if (result && result.success) {
      fileList.value = result.files;
      message.success(`找到 ${result.count} 个文件/目录`);
    } else if (result) {
      message.error(result.message);
    }
  } catch (error: any) {
    message.error(`列出文件失败: ${error.message}`);
  }
}

// 统计文件数
async function countFiles() {
  if (!localDirPath.value) {
    message.warning('请先选择本地目录');
    return;
  }
  try {
    const count = await FileTransferService.CountFilesInDirectory(localDirPath.value);
    fileCount.value = count;
    message.success(`统计完成: ${count} 个文件`);
  } catch (error: any) {
    message.error(`统计失败: ${error.message}`);
  }
}

// 测试TCP连接
async function testTCPConnection() {
  testingConnection.value = true;
  try {
    const config: TCPUDPConfig = {
      host: tcpConfig.host,
      port: tcpConfig.port,
      timeout: 10
    };
    const result = await TCPUDPService.TestTCPConnection(config);
    connectionResult.value = result;
    if (result && result.success) {
      message.success(`连接成功 (延迟: ${result.latency}ms)`);
    } else if (result) {
      message.error(result.message);
    }
  } catch (error: any) {
    message.error(`测试失败: ${error.message}`);
  } finally {
    testingConnection.value = false;
  }
}

// 连接TCP
async function connectTCP() {
  if (tcpConnected.value) {
    await TCPUDPService.DisconnectTCP();
    tcpConnected.value = false;
    tcpConnectionStatus.value = null;
    message.success('已断开连接');
  } else {
    connecting.value = true;
    try {
      const config: TCPUDPConfig = {
        host: tcpConfig.host,
        port: tcpConfig.port,
        timeout: 10
      };
      const status = await TCPUDPService.ConnectTCP(config);
      tcpConnectionStatus.value = status;
      tcpConnected.value = status ? status.isConnected : false;
      if (status && status.isConnected) {
        message.success('连接成功');
      } else {
        message.error('连接失败');
      }
    } catch (error: any) {
      message.error(`连接失败: ${error.message}`);
    } finally {
      connecting.value = false;
    }
  }
}

// 发送数据
async function sendDataToServer() {
  if (!sendData.value) {
    message.warning('请输入要发送的数据');
    return;
  }
  try {
    let result;
    if (tcpProtocol.value === 'tcp') {
      result = await TCPUDPService.SendTCP(sendData.value, sendAsHex.value);
    } else {
      const config: TCPUDPConfig = {
        host: tcpConfig.host,
        port: tcpConfig.port,
        timeout: 10
      };
      result = await TCPUDPService.SendUDP(config, sendData.value, sendAsHex.value);
    }
    if (result && result.success) {
      message.success(`发送成功 (${result.length} 字节)`);
    } else if (result) {
      message.error(result.message);
    }
  } catch (error: any) {
    message.error(`发送失败: ${error.message}`);
  }
}

// 接收数据
async function receiveDataFromServer() {
  try {
    let result;
    if (tcpProtocol.value === 'tcp') {
      result = await TCPUDPService.ReceiveTCP(5);
    } else {
      const config: TCPUDPConfig = {
        host: tcpConfig.host,
        port: tcpConfig.port,
        timeout: 10
      };
      result = await TCPUDPService.ReceiveUDP(config, 5);
    }
    if (result && result.success) {
      receivedData.value = result.data;
      receivedHexData.value = result.hexData;
      message.success(`接收成功 (${result.length} 字节)`);
    } else if (result) {
      message.error(result.message);
    }
  } catch (error: any) {
    message.error(`接收失败: ${error.message}`);
  }
}

// 转换为十六进制
async function convertToHex() {
  if (!sendData.value) {
    message.warning('请输入要转换的数据');
    return;
  }
  try {
    const hex = await TCPUDPService.StringToHex(sendData.value);
    sendData.value = hex;
    sendAsHex.value = true;
    message.success('已转换为十六进制');
  } catch (error: any) {
    message.error(`转换失败: ${error.message}`);
  }
}

// 从十六进制转换
async function convertFromHex() {
  if (!sendData.value) {
    message.warning('请输入十六进制数据');
    return;
  }
  try {
    const text = await TCPUDPService.HexToString(sendData.value);
    sendData.value = text;
    sendAsHex.value = false;
    message.success('已转换为文本');
  } catch (error: any) {
    message.error(`转换失败: ${error.message}`);
  }
}

// 格式化文件大小
function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}
</script>

<style scoped lang="postcss">
.network-debug {
  height: 100%;
  overflow: auto;
}

.protocol-section {
  padding: 16px;
}

.config-form {
  background: #f5f5f5;
  padding: 16px;
  border-radius: 4px;
}
</style>
