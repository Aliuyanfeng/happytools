<template>
  <div class="network-debug">
    <!-- 协议选择 -->
    <div class="protocol-selector">
      <div
        v-for="proto in protocols"
        :key="proto.key"
        class="protocol-item"
        :class="{ active: activeProtocol === proto.key }"
        @click="activeProtocol = proto.key"
      >
        <component :is="proto.icon" class="protocol-icon" />
        <span>{{ proto.name }}</span>
      </div>
    </div>

    <!-- UDP 调试面板 -->
    <div v-if="activeProtocol === 'udp'" class="debug-panel">
      <div class="panel-content">
        <!-- 左侧：配置区 -->
        <div class="config-section">
          <!-- 服务端配置 -->
          <div class="config-block">
            <div class="block-title">{{ t('network.serverConfig') }}</div>
            <div class="form-row inline-row">
              <label>{{ t('network.listenAddress') }}</label>
              <a-input v-model:value="serverConfig.host" style="width: 90px" placeholder="0.0.0.0" size="small" />
            </div>
            <div class="form-row inline-row">
              <label>{{ t('network.listenPort') }}</label>
              <a-input-number v-model:value="serverConfig.port" :min="1" :max="65535" style="width: 90px" size="small" />
            </div>
            <a-button
              :type="serverRunning ? 'default' : 'primary'"
              :danger="serverRunning"
              @click="toggleServer"
              block
              size="small"
            >
              <span class="status-dot" :class="{ active: serverRunning }"></span>
              {{ serverRunning ? t('network.stopServer') : t('network.startServer') }}
            </a-button>
            <div v-if="serverStatus?.localAddr" class="status-text">
              {{ t('network.listeningOn') }}: <span class="highlight">{{ serverStatus.localAddr }}</span>
            </div>
          </div>

          <!-- 发送配置 -->
          <div class="config-block">
            <div class="block-title">{{ t('network.sendConfig') }}</div>
            <div class="form-row inline-row">
              <label>{{ t('network.sendFormat') }}</label>
              <a-radio-group v-model:value="sendFormat" size="small">
                <a-radio-button value="text">ASC</a-radio-button>
                <a-radio-button value="hex">HEX</a-radio-button>
              </a-radio-group>
            </div>
            <div class="form-row inline-row">
              <label>{{ t('network.periodicSend') }}</label>
              <div class="periodic-row">
                <a-switch v-model:checked="periodicSend" size="small" />
                <a-input-number
                  
                  v-model:value="periodicInterval"
                  :min="100"
                  :max="60000"
                  :step="100"
                  size="small"
                  style="width: 70px"
                />
                <span class="unit-label">ms</span>
              </div>
            </div>
          </div>

          <!-- 接收配置 -->
          <div class="config-block">
            <div class="block-title">{{ t('network.receiveConfig') }}</div>
            <div class="form-row inline-row">
              <label>{{ t('network.receiveFormat') }}</label>
              <a-radio-group v-model:value="receiveFormat" size="small">
                <a-radio-button value="text">ASC</a-radio-button>
                <a-radio-button value="hex">HEX</a-radio-button>
              </a-radio-group>
            </div>
          </div>

          <!-- 统计信息 -->
          <div class="config-block stats-block">
            <div class="block-title">{{ t('network.statistics') }}</div>
            <div class="stats-grid">
              <div class="stat-item">
                <div class="stat-value send">{{ sendCount }}</div>
                <div class="stat-label">{{ t('network.sentCount') }}</div>
              </div>
              <div class="stat-item">
                <div class="stat-value receive">{{ receiveCount }}</div>
                <div class="stat-label">{{ t('network.receivedCount') }}</div>
              </div>
            </div>
            <div class="stats-row">
              <span class="stats-label">{{ t('network.sendBytes') }}:</span>
              <span class="stats-value">{{ sendBytes }} B</span>
            </div>
            <div class="stats-row">
              <span class="stats-label">{{ t('network.receiveBytes') }}:</span>
              <span class="stats-value">{{ receiveBytes }} B</span>
            </div>
          </div>
        </div>

        <!-- 右侧：数据区 -->
        <div class="data-section">
          <!-- 发送区 -->
          <div class="send-area">
            <div class="section-header">
              <span class="section-title">{{ t('network.sendData') }}</span>
            </div>
            <a-textarea
              v-model:value="sendData"
              :placeholder="sendFormat === 'hex' ? t('network.hexPlaceholder') : t('network.textPlaceholder')"
              :rows="3"
              class="data-input"
            />
            <div class="send-footer">
              <span class="byte-count">{{ t('network.byteCount', { count: getSendByteCount() }) }}</span>
              <div class="target-input">
                <span class="target-label">{{ t('network.target') }}:</span>
                <a-input v-model:value="targetConfig.host" style="width: 100px" placeholder="127.0.0.1" />
                <span class="colon">:</span>
                <a-input-number v-model:value="targetConfig.port" :min="1" :max="65535" style="width: 70px" />
                <a-button type="primary" @click="sendDataToTarget" :loading="sending" :disabled="!canSend" size="small">
                  {{ t('network.send') }}
                </a-button>
              </div>
            </div>
          </div>

          <!-- 接收区 -->
          <div class="receive-area">
            <div class="section-header">
              <span class="section-title">{{ t('network.receiveData') }}</span>
              <a-button size="small" @click="clearMessages">
                {{ t('network.clear') }}
              </a-button>
            </div>
            <div class="receive-output" ref="receiveOutputRef">
              <div v-if="messages.length === 0" class="empty-hint">
                {{ serverRunning ? t('network.waitingForData') : t('network.startServerHint') }}
              </div>
              <div v-else class="message-list">
                <div
                  v-for="(msg, index) in messages"
                  :key="index"
                  class="message-item"
                  :class="{ sent: msg.type === 'sent', received: msg.type === 'received' }"
                >
                  <div class="message-header">
                    <span class="message-direction">
                      {{ msg.type === 'sent' ? '↑' : '↓' }}
                      {{ msg.type === 'sent' ? t('network.sent') : t('network.received') }}
                    </span>
                    <span class="message-time">{{ msg.time }}</span>
                    <span class="message-addr">
                      {{ msg.type === 'sent' ? '→' : '←' }}
                      {{ msg.type === 'sent' ? msg.target : msg.from }}
                    </span>
                  </div>
                  <div class="message-content">
                    {{ receiveFormat === 'hex' ? msg.hexData : msg.data }}
                  </div>
                  <div class="message-meta">
                    {{ t('network.byteCount', { count: msg.length }) }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TCP 调试面板 -->
    <div v-else-if="activeProtocol === 'tcp'" class="debug-panel">
      <div class="panel-content">
        <!-- 左侧配置 -->
        <div class="config-section">
          <!-- 模式切换 -->
          <div class="config-block">
            <div class="block-title">模式</div>
            <a-radio-group v-model:value="tcpMode" size="small" button-style="solid" style="width:100%">
              <a-radio-button value="client" style="width:50%;text-align:center">客户端</a-radio-button>
              <a-radio-button value="server" style="width:50%;text-align:center">服务端</a-radio-button>
            </a-radio-group>
          </div>

          <!-- 客户端配置 -->
          <div class="config-block" v-if="tcpMode === 'client'">
            <div class="block-title">连接配置</div>
            <div class="form-row inline-row">
              <label>主机</label>
              <a-input v-model:value="tcpClientConfig.host" style="width:100px" placeholder="127.0.0.1" size="small" />
            </div>
            <div class="form-row inline-row">
              <label>端口</label>
              <a-input-number v-model:value="tcpClientConfig.port" :min="1" :max="65535" style="width:80px" size="small" />
            </div>
            <a-button
              :type="tcpClientConnected ? 'default' : 'primary'"
              :danger="tcpClientConnected"
              @click="toggleTCPClient"
              block size="small"
            >
              <span class="status-dot" :class="{ active: tcpClientConnected }"></span>
              {{ tcpClientConnected ? '断开' : '连接' }}
            </a-button>
            <div v-if="tcpClientStatus" class="status-text">
              本地: <span class="highlight">{{ tcpClientStatus.localAddr }}</span><br/>
              远端: <span class="highlight">{{ tcpClientStatus.remoteAddr }}</span>
            </div>
          </div>

          <!-- 服务端配置 -->
          <div class="config-block" v-else>
            <div class="block-title">监听配置</div>
            <div class="form-row inline-row">
              <label>地址</label>
              <a-input v-model:value="tcpServerConfig.host" style="width:100px" placeholder="0.0.0.0" size="small" />
            </div>
            <div class="form-row inline-row">
              <label>端口</label>
              <a-input-number v-model:value="tcpServerConfig.port" :min="1" :max="65535" style="width:80px" size="small" />
            </div>
            <a-button
              :type="tcpServerRunning ? 'default' : 'primary'"
              :danger="tcpServerRunning"
              @click="toggleTCPServer"
              block size="small"
            >
              <span class="status-dot" :class="{ active: tcpServerRunning }"></span>
              {{ tcpServerRunning ? '停止' : '启动' }}
            </a-button>
            <div v-if="tcpServerAddr" class="status-text">
              监听: <span class="highlight">{{ tcpServerAddr }}</span>
            </div>
            <!-- 已连接客户端 -->
            <div v-if="tcpServerClients.length" class="mt-2">
              <div class="block-title">已连接 ({{ tcpServerClients.length }})</div>
              <div v-for="c in tcpServerClients" :key="c.id" class="client-item">
                <span class="status-dot active"></span>
                <span style="font-size:10px">{{ c.remoteAddr }}</span>
              </div>
            </div>
          </div>

          <!-- 发送格式 -->
          <div class="config-block">
            <div class="block-title">格式</div>
            <div class="form-row inline-row">
              <label>发送</label>
              <a-radio-group v-model:value="tcpSendFormat" size="small">
                <a-radio-button value="text">ASC</a-radio-button>
                <a-radio-button value="hex">HEX</a-radio-button>
              </a-radio-group>
            </div>
            <div class="form-row inline-row">
              <label>接收</label>
              <a-radio-group v-model:value="tcpReceiveFormat" size="small">
                <a-radio-button value="text">ASC</a-radio-button>
                <a-radio-button value="hex">HEX</a-radio-button>
              </a-radio-group>
            </div>
          </div>

          <!-- 统计 -->
          <div class="config-block stats-block">
            <div class="block-title">统计</div>
            <div class="stats-grid">
              <div class="stat-item"><div class="stat-value send">{{ tcpSendCount }}</div><div class="stat-label">发送</div></div>
              <div class="stat-item"><div class="stat-value receive">{{ tcpRecvCount }}</div><div class="stat-label">接收</div></div>
            </div>
          </div>
        </div>

        <!-- 右侧数据区 -->
        <div class="data-section">
          <div class="send-area">
            <div class="section-header">
              <span class="section-title">发送</span>
            </div>
            <a-textarea v-model:value="tcpSendData" :rows="3" class="data-input"
              :placeholder="tcpSendFormat === 'hex' ? '十六进制，如: 48656C6C6F' : '输入要发送的内容...'" />
            <div class="send-footer">
              <span class="byte-count">{{ tcpSendByteCount }} 字节</span>
              <a-button type="primary" size="small" @click="sendTCPData"
                :disabled="!tcpCanSend" :loading="tcpSending">发送</a-button>
            </div>
          </div>

          <div class="receive-area">
            <div class="section-header">
              <span class="section-title">收发记录</span>
              <a-button size="small" @click="tcpMessages = []">清空</a-button>
            </div>
            <div class="receive-output" ref="tcpOutputRef">
              <div v-if="tcpMessages.length === 0" class="empty-hint">
                {{ tcpMode === 'client' ? (tcpClientConnected ? '等待数据...' : '请先建立连接') : (tcpServerRunning ? '等待客户端连接...' : '请先启动服务端') }}
              </div>
              <div v-else class="message-list">
                <div v-for="(msg, i) in tcpMessages" :key="i" class="message-item" :class="msg.type">
                  <div class="message-header">
                    <span class="message-direction">{{ msg.type === 'sent' ? '↑ 发送' : '↓ 接收' }}</span>
                    <span class="message-time">{{ msg.time }}</span>
                    <span v-if="msg.addr" class="message-addr">{{ msg.addr }}</span>
                  </div>
                  <div class="message-content">{{ tcpReceiveFormat === 'hex' ? msg.hexData : msg.data }}</div>
                  <div class="message-meta">{{ msg.length }} 字节</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- WebSocket 调试面板 -->
    <div v-else-if="activeProtocol === 'websocket'" class="debug-panel">
      <div class="panel-content">
        <!-- 左侧配置 -->
        <div class="config-section">
          <div class="config-block">
            <div class="block-title">连接</div>
            <div class="form-row">
              <label>URL</label>
              <a-input v-model:value="wsUrl" placeholder="ws://127.0.0.1:8080/ws" size="small" />
            </div>
            <a-button
              :type="wsConnected ? 'default' : 'primary'"
              :danger="wsConnected"
              @click="toggleWS"
              block size="small"
              style="margin-top:8px"
            >
              <span class="status-dot" :class="{ active: wsConnected }"></span>
              {{ wsConnected ? '断开' : '连接' }}
            </a-button>
          </div>

          <div class="config-block">
            <div class="block-title">发送格式</div>
            <a-radio-group v-model:value="wsSendFormat" size="small" style="width:100%">
              <a-radio-button value="text" style="width:50%;text-align:center">文本</a-radio-button>
              <a-radio-button value="binary" style="width:50%;text-align:center">二进制</a-radio-button>
            </a-radio-group>
          </div>

          <div class="config-block stats-block">
            <div class="block-title">统计</div>
            <div class="stats-grid">
              <div class="stat-item"><div class="stat-value send">{{ wsSendCount }}</div><div class="stat-label">发送</div></div>
              <div class="stat-item"><div class="stat-value receive">{{ wsRecvCount }}</div><div class="stat-label">接收</div></div>
            </div>
          </div>
        </div>

        <!-- 右侧数据区 -->
        <div class="data-section">
          <div class="send-area">
            <div class="section-header">
              <span class="section-title">发送</span>
            </div>
            <a-textarea v-model:value="wsSendData" :rows="3" class="data-input"
              :placeholder="wsSendFormat === 'binary' ? '十六进制，如: 48656C6C6F' : '输入要发送的内容...'" />
            <div class="send-footer">
              <span class="byte-count">{{ wsSendByteCount }} 字节</span>
              <a-button type="primary" size="small" @click="sendWSData"
                :disabled="!wsCanSend" :loading="wsSending">发送</a-button>
            </div>
          </div>

          <div class="receive-area">
            <div class="section-header">
              <span class="section-title">收发记录</span>
              <a-button size="small" @click="wsMessages = []">清空</a-button>
            </div>
            <div class="receive-output" ref="wsOutputRef">
              <div v-if="wsMessages.length === 0" class="empty-hint">
                {{ wsConnected ? '等待消息...' : '请先建立连接' }}
              </div>
              <div v-else class="message-list">
                <div v-for="(msg, i) in wsMessages" :key="i" class="message-item" :class="msg.type">
                  <div class="message-header">
                    <span class="message-direction">{{ msg.type === 'sent' ? '↑ 发送' : '↓ 接收' }}</span>
                    <span class="message-time">{{ msg.time }}</span>
                    <span class="message-addr">{{ msg.isText ? 'Text' : 'Binary' }}</span>
                  </div>
                  <div class="message-content">{{ msg.data }}</div>
                  <div class="message-meta">{{ msg.length }} 字节</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import { ApiOutlined, WifiOutlined, CloudOutlined } from '@ant-design/icons-vue'
import { Events } from '@wailsio/runtime'
import { TCPUDPService, type ConnectionStatus, type MessageResult } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/network'
import { WebSocketService } from '../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/network'

const { t } = useI18n()

message.config({
  top: `100px`,
  duration: 2,
  maxCount: 1,
});

// 协议列表
const protocols = [
  { key: 'udp', name: 'UDP', icon: WifiOutlined },
  { key: 'tcp', name: 'TCP', icon: ApiOutlined },
  { key: 'websocket', name: 'WebSocket', icon: CloudOutlined },
]

const activeProtocol = ref('udp')

// 服务端配置
const serverConfig = reactive({
  host: '0.0.0.0',
  port: 8080,
})

// 目标配置
const targetConfig = reactive({
  host: '127.0.0.1',
  port: 8080,
})

// 数据格式
const sendFormat = ref<'text' | 'hex'>('text')
const receiveFormat = ref<'text' | 'hex'>('text')

// 周期发送
const periodicSend = ref(false)
const periodicInterval = ref(1000)
let periodicTimer: ReturnType<typeof setInterval> | null = null

// 发送/接收数据
const sendData = ref('')
const messages = ref<Array<{
  type: 'sent' | 'received'
  data: string
  hexData: string
  length: number
  time: string
  from?: string
  target?: string
}>>([])

// 统计
const sendCount = ref(0)
const receiveCount = ref(0)
const sendBytes = ref(0)
const receiveBytes = ref(0)

// 接收区域滚动引用
const receiveOutputRef = ref<HTMLElement | null>(null)

// 状态
const sending = ref(false)
const serverRunning = ref(false)
const serverStatus = ref<ConnectionStatus | null>(null)

// 计算属性
const canSend = computed(() => {
  return sendData.value.trim() !== '' && targetConfig.host && targetConfig.port
})

// 获取发送字节数
function getSendByteCount(): number {
  if (!sendData.value) return 0
  if (sendFormat.value === 'hex') {
    const hex = sendData.value.replace(/\s/g, '')
    return Math.floor(hex.length / 2)
  }
  return new TextEncoder().encode(sendData.value).length
}

// 启动/停止服务端
async function toggleServer() {
  if (serverRunning.value) {
    await TCPUDPService.StopUDPServer()
    serverRunning.value = false
    serverStatus.value = null
    message.success(t('network.serverStopped'))
  } else {
    const status = await TCPUDPService.StartUDPServer(serverConfig.host, serverConfig.port)
    if (status && status.isConnected) {
      serverRunning.value = true
      serverStatus.value = status
      message.success(t('network.serverStarted', { addr: status.localAddr }))
    } else {
      message.error(t('network.serverStartFailed'))
    }
  }
}

// 发送数据到目标
async function sendDataToTarget() {
  if (!sendData.value.trim()) {
    message.warning(t('network.enterData'))
    return
  }

  sending.value = true
  try {
    const isHex = sendFormat.value === 'hex'
    const result = await TCPUDPService.SendUDPFromServer(
      targetConfig.host,
      targetConfig.port,
      sendData.value,
      isHex
    )

    if (result && result.success) {
      addMessage({
        type: 'sent',
        data: result.data || '',
        hexData: result.hexData || '',
        length: result.length || 0,
        time: new Date().toLocaleTimeString(),
        target: `${targetConfig.host}:${targetConfig.port}`,
      })
      // 更新统计
      sendCount.value++
      sendBytes.value += result.length || 0
      // message.success(t('network.sendSuccess'))
    } else if (result) {
      message.error(result.message)
    }
  } catch (error: any) {
    message.error(`${t('network.sendFailed')}: ${error.message}`)
  } finally {
    sending.value = false
  }
}

// 添加消息到列表并自动滚动
function addMessage(msg: {
  type: 'sent' | 'received'
  data: string
  hexData: string
  length: number
  time: string
  from?: string
  target?: string
}) {
  messages.value.push(msg)
  nextTick(() => {
    if (receiveOutputRef.value) {
      receiveOutputRef.value.scrollTop = receiveOutputRef.value.scrollHeight
    }
  })
}

// 处理UDP接收事件
function handleUDPReceived(event: any) {
  const result = event.data as MessageResult
  if (result && result.success) {
    const fromMatch = result.message?.match(/\(来自: (.+)\)/)
    const from = fromMatch ? fromMatch[1] : 'unknown'

    addMessage({
      type: 'received',
      data: result.data || '',
      hexData: result.hexData || '',
      length: result.length || 0,
      time: new Date().toLocaleTimeString(),
      from,
    })
    // 更新统计
    receiveCount.value++
    receiveBytes.value += result.length || 0
  }
}

// 清空消息
function clearMessages() {
  messages.value = []
  sendCount.value = 0
  receiveCount.value = 0
  sendBytes.value = 0
  receiveBytes.value = 0
}

// ─── TCP ────────────────────────────────────────────────────────
const tcpMode = ref<'client' | 'server'>('client')
const tcpClientConfig = reactive({ host: '127.0.0.1', port: 8080 })
const tcpServerConfig = reactive({ host: '0.0.0.0', port: 9090 })
const tcpClientConnected = ref(false)
const tcpClientStatus = ref<any>(null)
const tcpServerRunning = ref(false)
const tcpServerAddr = ref('')
const tcpServerClients = ref<any[]>([])
const tcpSendFormat = ref<'text' | 'hex'>('text')
const tcpReceiveFormat = ref<'text' | 'hex'>('text')
const tcpSendData = ref('')
const tcpSending = ref(false)
const tcpSendCount = ref(0)
const tcpRecvCount = ref(0)
const tcpMessages = ref<any[]>([])
const tcpOutputRef = ref<HTMLElement | null>(null)

const tcpSendByteCount = computed(() => {
  if (!tcpSendData.value) return 0
  if (tcpSendFormat.value === 'hex') return Math.floor(tcpSendData.value.replace(/\s/g, '').length / 2)
  return new TextEncoder().encode(tcpSendData.value).length
})

const tcpCanSend = computed(() =>
  tcpSendData.value.trim() !== '' && (tcpClientConnected.value || tcpServerRunning.value)
)

function addTCPMessage(msg: any) {
  tcpMessages.value.push(msg)
  nextTick(() => { if (tcpOutputRef.value) tcpOutputRef.value.scrollTop = tcpOutputRef.value.scrollHeight })
}

async function toggleTCPClient() {
  if (tcpClientConnected.value) {
    await TCPUDPService.StopTCPClient()
    tcpClientConnected.value = false
    tcpClientStatus.value = null
  } else {
    try {
      const status = await TCPUDPService.StartTCPClient(tcpClientConfig.host, tcpClientConfig.port, 10)
      if (status?.isConnected) {
        tcpClientConnected.value = true
        tcpClientStatus.value = status
        message.success('连接成功')
      } else {
        message.error('连接失败')
      }
    } catch (e: any) { message.error('连接失败: ' + e.message) }
  }
}

async function toggleTCPServer() {
  if (tcpServerRunning.value) {
    await TCPUDPService.StopTCPServer()
    tcpServerRunning.value = false
    tcpServerAddr.value = ''
    tcpServerClients.value = []
  } else {
    try {
      const status = await TCPUDPService.StartTCPServer(tcpServerConfig.host, tcpServerConfig.port)
      if (status?.isConnected) {
        tcpServerRunning.value = true
        tcpServerAddr.value = status.localAddr
        message.success('服务端已启动: ' + status.localAddr)
      } else {
        message.error('启动失败')
      }
    } catch (e: any) { message.error('启动失败: ' + e.message) }
  }
}

async function sendTCPData() {
  if (!tcpSendData.value.trim()) return
  tcpSending.value = true
  try {
    const isHex = tcpSendFormat.value === 'hex'
    let result: any
    if (tcpMode.value === 'client') {
      result = await TCPUDPService.SendTCPClient(tcpSendData.value, isHex)
    } else {
      result = await TCPUDPService.SendTCPServerToClient('', tcpSendData.value, isHex)
    }
    if (result?.success) {
      addTCPMessage({ type: 'sent', data: result.data, hexData: result.hexData, length: result.length, time: new Date().toLocaleTimeString() })
      tcpSendCount.value++
    } else {
      message.error(result?.message || '发送失败')
    }
  } catch (e: any) { message.error('发送失败: ' + e.message) }
  finally { tcpSending.value = false }
}

function handleTCPClientReceived(event: any) {
  const d = event?.data
  if (!d) return
  addTCPMessage({ type: 'received', data: d.data, hexData: d.hexData, length: d.length, time: d.time })
  tcpRecvCount.value++
}

function handleTCPServerReceived(event: any) {
  const d = event?.data
  if (!d) return
  addTCPMessage({ type: 'received', data: d.data, hexData: d.hexData, length: d.length, time: d.time, addr: d.remoteAddr })
  tcpRecvCount.value++
}

// ─── WebSocket ───────────────────────────────────────────────────
const wsUrl = ref('ws://127.0.0.1:8080/ws')
const wsConnected = ref(false)
const wsSendFormat = ref<'text' | 'binary'>('text')
const wsSendData = ref('')
const wsSending = ref(false)
const wsSendCount = ref(0)
const wsRecvCount = ref(0)
const wsMessages = ref<any[]>([])
const wsOutputRef = ref<HTMLElement | null>(null)

const wsSendByteCount = computed(() => {
  if (!wsSendData.value) return 0
  if (wsSendFormat.value === 'binary') return Math.floor(wsSendData.value.replace(/\s/g, '').length / 2)
  return new TextEncoder().encode(wsSendData.value).length
})

const wsCanSend = computed(() => wsConnected.value && wsSendData.value.trim() !== '')

function addWSMessage(msg: any) {
  wsMessages.value.push(msg)
  nextTick(() => { if (wsOutputRef.value) wsOutputRef.value.scrollTop = wsOutputRef.value.scrollHeight })
}

async function toggleWS() {
  if (wsConnected.value) {
    await WebSocketService.Disconnect()
    wsConnected.value = false
  } else {
    try {
      const status = await WebSocketService.Connect(wsUrl.value)
      if (status?.isConnected) {
        wsConnected.value = true
        message.success('WebSocket 连接成功')
      } else {
        message.error('连接失败')
      }
    } catch (e: any) { message.error('连接失败: ' + e.message) }
  }
}

async function sendWSData() {
  if (!wsSendData.value.trim()) return
  wsSending.value = true
  try {
    let result: any
    if (wsSendFormat.value === 'binary') {
      result = await WebSocketService.SendBinary(wsSendData.value)
    } else {
      result = await WebSocketService.SendText(wsSendData.value)
    }
    if (result?.success) {
      addWSMessage({ type: 'sent', data: result.data, length: result.length, time: new Date().toLocaleTimeString(), isText: wsSendFormat.value === 'text' })
      wsSendCount.value++
    } else {
      message.error(result?.message || '发送失败')
    }
  } catch (e: any) { message.error('发送失败: ' + e.message) }
  finally { wsSending.value = false }
}

function handleWSReceived(event: any) {
  const d = event?.data
  if (!d) return
  addWSMessage({ type: 'received', data: d.isText ? d.data : d.hexData, length: d.length, time: d.time, isText: d.isText })
  wsRecvCount.value++
}

// 周期发送控制
watch(periodicSend, (enabled) => {
  if (enabled) {
    periodicTimer = setInterval(() => {
      if (canSend.value && !sending.value) {
        sendDataToTarget()
      }
    }, periodicInterval.value)
  } else {
    if (periodicTimer) {
      clearInterval(periodicTimer)
      periodicTimer = null
    }
  }
})

// 组件挂载时注册事件监听
onMounted(() => {
  Events.On('network:udpReceived', handleUDPReceived)
  Events.On('network:tcpClientReceived', handleTCPClientReceived)
  Events.On('network:tcpClientDisconnectedSelf', () => {
    tcpClientConnected.value = false
    tcpClientStatus.value = null
    message.warning('TCP 连接已断开')
  })
  Events.On('network:tcpServerReceived', handleTCPServerReceived)
  Events.On('network:tcpClientConnected', (event: any) => {
    const info = event?.data
    if (info) tcpServerClients.value.push(info)
  })
  Events.On('network:tcpClientDisconnected', (event: any) => {
    const id = event?.data
    tcpServerClients.value = tcpServerClients.value.filter((c: any) => c.id !== id)
  })
  Events.On('network:wsReceived', handleWSReceived)
  Events.On('network:wsDisconnected', () => {
    wsConnected.value = false
    message.warning('WebSocket 连接已断开')
  })
})

// 组件卸载时取消事件监听
onUnmounted(() => {
  Events.Off('network:udpReceived')
  Events.Off('network:tcpClientReceived')
  Events.Off('network:tcpClientDisconnectedSelf')
  Events.Off('network:tcpServerReceived')
  Events.Off('network:tcpClientConnected')
  Events.Off('network:tcpClientDisconnected')
  Events.Off('network:wsReceived')
  Events.Off('network:wsDisconnected')
  if (periodicTimer) clearInterval(periodicTimer)
})
</script>

<style scoped>
.network-debug {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

/* 协议选择器 */
.protocol-selector {
  display: flex;
  gap: 8px;
  padding: 16px 20px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
}

.protocol-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  background: #f5f7fa;
  color: #666;
  font-weight: 500;
}

.protocol-item:hover {
  background: #e6f7ff;
  color: #1890ff;
}

.protocol-item.active {
  background: #1890ff;
  color: #fff;
}

.protocol-icon {
  font-size: 16px;
}

/* 调试面板 */
.debug-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 20px;
  overflow: hidden;
}

/* 面板内容 */
.panel-content {
  flex: 1;
  display: flex;
  gap: 20px;
  overflow: hidden;
}

/* 左侧配置区 */
.config-section {
  width: 200px;
  flex-shrink: 0;
  background: #fff;
  border-radius: 12px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  overflow: hidden;
}

.config-block {
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
}

.config-block:last-child {
  padding-bottom: 0;
  border-bottom: none;
}

.block-title {
  font-size: 12px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.form-row {
  margin-bottom: 6px;
}

.form-row label {
  display: block;
  font-size: 11px;
  color: #666;
  margin-bottom: 2px;
}

.form-row.inline-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.form-row.inline-row label {
  margin-bottom: 0;
  flex-shrink: 0;
}

.format-radio {
  width: 100%;
}

.format-radio :deep(.ant-radio-button-wrapper) {
  width: 50%;
  text-align: center;
}

.form-row.inline-row :deep(.ant-radio-group) {
  display: flex;
}

.form-row.inline-row :deep(.ant-radio-button-wrapper) {
  padding: 0 8px;
  font-size: 11px;
}

.periodic-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

.unit-label {
  font-size: 11px;
  color: #666;
}

.status-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #d9d9d9;
  margin-right: 6px;
}

.status-dot.active {
  background: #52c41a;
  box-shadow: 0 0 4px rgba(82, 196, 26, 0.4);
}

.status-text {
  margin-top: 6px;
  font-size: 10px;
  color: #666;
}

.status-text .highlight {
  color: #1890ff;
  font-family: 'Consolas', 'Monaco', monospace;
}

/* 统计区块 */
.stats-block {
  flex: 0 0 auto;
  display: flex;
  flex-direction: column;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-bottom: 8px;
}

.stat-item {
  text-align: center;
  padding: 6px;
  background: #f5f7fa;
  border-radius: 6px;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  font-family: 'Consolas', 'Monaco', monospace;
}

.stat-value.send {
  color: #1890ff;
}

.stat-value.receive {
  color: #52c41a;
}

.stat-label {
  font-size: 10px;
  color: #999;
  margin-top: 2px;
}

.stats-row {
  display: flex;
  justify-content: space-between;
  font-size: 10px;
  margin-bottom: 2px;
}

.stats-label {
  color: #666;
}

.stats-value {
  color: #333;
  font-family: 'Consolas', 'Monaco', monospace;
}

/* 数据区 */
.data-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.send-area,
.receive-area {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
}

.send-area {
  flex: 0 0 auto;
}

.receive-area {
  flex: 1;
  min-height: 0;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #333;
}

.data-input {
  font-family: 'Consolas', 'Monaco', monospace;
}

.send-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 12px;
  flex-wrap: wrap;
  gap: 12px;
}

.byte-count {
  font-size: 12px;
  color: #999;
}

.target-input {
  display: flex;
  align-items: center;
  gap: 6px;
}

.target-label {
  font-size: 12px;
  color: #666;
}

.colon {
  color: #999;
  font-weight: bold;
}

/* 接收输出 */
.receive-output {
  flex: 1;
  min-height: 0;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  overflow: auto;
  background: #fafafa;
}

.empty-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #bfbfbf;
  font-size: 14px;
}

.message-list {
  padding: 12px;
}

.message-item {
  padding: 10px 12px;
  margin-bottom: 8px;
  border-radius: 6px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
}

.message-item.sent {
  background: #e6f7ff;
  border-left: 3px solid #1890ff;
}

.message-item.received {
  background: #f6ffed;
  border-left: 3px solid #52c41a;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 6px;
  font-size: 12px;
  color: #666;
}

.message-direction {
  font-weight: 600;
}

.message-item.sent .message-direction {
  color: #1890ff;
}

.message-item.received .message-direction {
  color: #52c41a;
}

.message-time {
  color: #999;
}

.message-addr {
  color: #1890ff;
  margin-left: auto;
}

.message-content {
  word-break: break-all;
  white-space: pre-wrap;
  color: #333;
  line-height: 1.5;
}

.message-meta {
  margin-top: 6px;
  font-size: 11px;
  color: #999;
}

/* Coming Soon */
.coming-soon {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: #bfbfbf;
}

.coming-icon {
  font-size: 48px;
}

.client-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 3px 0;
  font-size: 10px;
  color: #555;
}

/* 深色模式适配 */
:root[theme='dark'] .network-debug {
  background: #141414;
}

:root[theme='dark'] .protocol-selector {
  background: #1f1f1f;
  border-bottom-color: #303030;
}

:root[theme='dark'] .protocol-item {
  background: #303030;
  color: #aaa;
}

:root[theme='dark'] .protocol-item:hover {
  background: #177ddc;
  color: #fff;
}

:root[theme='dark'] .protocol-item.active {
  background: #177ddc;
}

:root[theme='dark'] .config-section,
:root[theme='dark'] .send-area,
:root[theme='dark'] .receive-area {
  background: #1f1f1f;
}

:root[theme='dark'] .config-block {
  border-bottom-color: #303030;
}

:root[theme='dark'] .block-title {
  color: #e0e0e0;
}

:root[theme='dark'] .form-row label,
:root[theme='dark'] .target-label,
:root[theme='dark'] .stats-label {
  color: #aaa;
}

:root[theme='dark'] .status-text {
  color: #aaa;
}

:root[theme='dark'] .status-text .highlight {
  color: #177ddc;
}

:root[theme='dark'] .stat-item {
  background: #303030;
}

:root[theme='dark'] .stat-value.send {
  color: #177ddc;
}

:root[theme='dark'] .stat-value.receive {
  color: #49aa19;
}

:root[theme='dark'] .stats-value {
  color: #e0e0e0;
}

:root[theme='dark'] .receive-output {
  background: #141414;
  border-color: #303030;
}

:root[theme='dark'] .empty-hint {
  color: #555;
}

:root[theme='dark'] .message-item.sent {
  background: #111d2c;
  border-left-color: #177ddc;
}

:root[theme='dark'] .message-item.received {
  background: #162312;
  border-left-color: #49aa19;
}

:root[theme='dark'] .message-header {
  color: #888;
}

:root[theme='dark'] .message-item.sent .message-direction {
  color: #177ddc;
}

:root[theme='dark'] .message-item.received .message-direction {
  color: #49aa19;
}

:root[theme='dark'] .message-addr {
  color: #177ddc;
}

:root[theme='dark'] .message-content {
  color: #e0e0e0;
}

:root[theme='dark'] .coming-soon {
  color: #555;
}
</style>
