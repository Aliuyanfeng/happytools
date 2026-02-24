<!--
 * @Author: LiuYanFeng
 * @Date: 2026-02-12
 * @Description: 加密工具组件 - 包含base64、md5、sha256、sha512、sha1编解码以及utf-8编解码
 -->
<template>
  <div class="crypto-tools">
    <a-tabs v-model:activeKey="activeTab" type="card" class="crypto-tabs">
      <!-- Base64 编解码标签页 -->
      <a-tab-pane key="base64" tab="Base64 编解码">
        <div class="crypto-section">
          <a-card title="Base64 编码/解码" :bordered="false">
            <div class="crypto-form">
              <a-row :gutter="16" class="mb-4">
                <a-col :span="24">
                  <a-textarea
                    v-model:value="base64Input"
                    placeholder="输入要编码或解码的文本"
                    :rows="4"
                    :maxlength="10000"
                    show-count
                    allow-clear
                    @change="handleBase64InputChange"
                  />
                </a-col>
              </a-row>
              
              <a-row :gutter="16" class="mb-4">
                <a-col :span="12">
                  <a-button 
                    type="primary" 
                    block 
                    @click="handleBase64Encode"
                    :loading="base64Loading"
                    :disabled="!base64Input"
                  >
                    <template #icon>
                      <LockOutlined />
                    </template>
                    Base64 编码
                  </a-button>
                </a-col>
                <a-col :span="12">
                  <a-button 
                    block 
                    @click="handleBase64Decode"
                    :loading="base64Loading"
                    :disabled="!base64Input"
                    :type="isValidBase64(base64Input) ? 'primary' : 'default'"
                  >
                    <template #icon>
                      <UnlockOutlined />
                    </template>
                    Base64 解码
                  </a-button>
                </a-col>
              </a-row>
              
              <div class="crypto-result" v-if="base64Result">
                <a-alert
                  :type="base64Result.type"
                  :message="base64Result.title"
                  :description="base64Result.content"
                  show-icon
                  closable
                  @close="clearBase64Result"
                >
                  <template #action>
                    <a-space>
                      <a-button size="small" @click="copyToClipboard(base64Result.content)">
                        复制
                      </a-button>
                      <a-button size="small" @click="clearBase64Result">
                        清除
                      </a-button>
                    </a-space>
                  </template>
                </a-alert>
              </div>
              

            </div>
          </a-card>
        </div>
      </a-tab-pane>
      
      <!-- 哈希计算标签页 -->
      <a-tab-pane key="hash" tab="哈希计算">
        <div class="crypto-section">
          <a-card title="哈希值计算" :bordered="false">
            <div class="crypto-form">
              <a-row :gutter="16" class="mb-4">
                <a-col :span="24">
                  <a-textarea
                    v-model:value="hashInput"
                    placeholder="输入要计算哈希值的文本"
                    :rows="4"
                    :maxlength="10000"
                    show-count
                    allow-clear
                    @change="handleHashInputChange"
                  />
                </a-col>
              </a-row>
              
              <a-row :gutter="16" class="mb-4">
                <a-col :span="6">
                  <a-button 
                    block 
                    @click="handleMd5Hash"
                    :loading="hashLoading"
                    :disabled="!hashInput"
                  >
                    MD5
                  </a-button>
                </a-col>
                <a-col :span="6">
                  <a-button 
                    block 
                    @click="handleSha1Hash"
                    :loading="hashLoading"
                    :disabled="!hashInput"
                  >
                    SHA-1
                  </a-button>
                </a-col>
                <a-col :span="6">
                  <a-button 
                    block 
                    @click="handleSha256Hash"
                    :loading="hashLoading"
                    :disabled="!hashInput"
                  >
                    SHA-256
                  </a-button>
                </a-col>
                <a-col :span="6">
                  <a-button 
                    block 
                    @click="handleSha512Hash"
                    :loading="hashLoading"
                    :disabled="!hashInput"
                  >
                    SHA-512
                  </a-button>
                </a-col>
              </a-row>
              
              <a-row :gutter="16" class="mb-4">
                <a-col :span="24">
                  <a-button 
                    type="primary" 
                    block 
                    @click="handleAllHashes"
                    :loading="hashLoading"
                    :disabled="!hashInput"
                  >
                    <template #icon>
                      <CalculatorOutlined />
                    </template>
                    计算所有哈希值
                  </a-button>
                </a-col>
              </a-row>
              
              <div class="hash-results" v-if="hashResults.length > 0">
                <a-list
                  :data-source="hashResults"
                  :loading="hashLoading"
                >
                  <template #renderItem="{ item }">
                    <a-list-item>
                      <a-list-item-meta>
                        <template #title>
                          <div class="hash-title">
                            <span class="hash-type">{{ item.type }}</span>
                            <a-tag :color="getHashColor(item.type)">
                              {{ item.type }}
                            </a-tag>
                          </div>
                        </template>
                        <template #description>
                          <div class="hash-value">
                            <code>{{ item.value }}</code>
                            <a-space class="hash-actions">
                              <a-button 
                                size="small" 
                                type="text" 
                                @click="copyToClipboard(item.value)"
                                title="复制"
                              >
                                <CopyOutlined />
                              </a-button>
                              <a-button 
                                size="small" 
                                type="text" 
                                @click="formatHashValue(item)"
                                title="格式化"
                              >
                                <FormOutlined />
                              </a-button>
                              <a-button 
                                size="small" 
                                type="text" 
                                danger
                                @click="removeHashResult(item.id)"
                                title="删除"
                              >
                                <DeleteOutlined />
                              </a-button>
                            </a-space>
                          </div>
                          <div class="hash-length" v-if="item.formatted">
                            长度: {{ item.value.replace(/\s/g, '').length }} 字符
                          </div>
                        </template>
                      </a-list-item-meta>
                    </a-list-item>
                  </template>
                </a-list>
                
                <div class="hash-actions-bottom mt-4">
                  <a-space>
                    <a-button @click="clearAllHashResults">
                      <template #icon>
                        <DeleteOutlined />
                      </template>
                      清除所有结果
                    </a-button>
                    <a-button @click="copyAllHashResults">
                      <template #icon>
                        <CopyOutlined />
                      </template>
                      复制所有结果
                    </a-button>
                  </a-space>
                </div>
              </div>
              

            </div>
          </a-card>
        </div>
      </a-tab-pane>
      
      <!-- UTF-8 编解码标签页 -->
      <a-tab-pane key="utf8" tab="UTF-8 编解码">
        <div class="crypto-section">
          <a-card title="UTF-8 编码/解码" :bordered="false">
            <div class="crypto-form">
              <a-row :gutter="16" class="mb-4">
                <a-col :span="24">
                  <a-textarea
                    v-model:value="utf8Input"
                    placeholder="输入要编码或解码的文本"
                    :rows="4"
                    :maxlength="10000"
                    show-count
                    allow-clear
                    @change="handleUtf8InputChange"
                  />
                </a-col>
              </a-row>
              
              <a-row :gutter="16" class="mb-4">
                <a-col :span="12">
                  <a-button 
                    type="primary" 
                    block 
                    @click="handleUtf8Encode"
                    :loading="utf8Loading"
                    :disabled="!utf8Input"
                  >
                    <template #icon>
                      <CodeOutlined />
                    </template>
                    UTF-8 编码（转十六进制）
                  </a-button>
                </a-col>
                <a-col :span="12">
                  <a-button 
                    block 
                    @click="handleUtf8Decode"
                    :loading="utf8Loading"
                    :disabled="!utf8Input"
                    :type="isValidHex(utf8Input) ? 'primary' : 'default'"
                  >
                    <template #icon>
                      <FileTextOutlined />
                    </template>
                    UTF-8 解码（十六进制转文本）
                  </a-button>
                </a-col>
              </a-row>
              
              <div class="crypto-result" v-if="utf8Result">
                <a-alert
                  :type="utf8Result.type"
                  :message="utf8Result.title"
                  :description="utf8Result.content"
                  show-icon
                  closable
                  @close="clearUtf8Result"
                >
                  <template #action>
                    <a-space>
                      <a-button size="small" @click="copyToClipboard(utf8Result.content)">
                        复制
                      </a-button>
                      <a-button size="small" @click="formatUtf8Result">
                        格式化
                      </a-button>
                      <a-button size="small" @click="clearUtf8Result">
                        清除
                      </a-button>
                    </a-space>
                  </template>
                </a-alert>
              </div>
              
              <div class="utf8-examples mt-4">
                <a-collapse>
                  <a-collapse-panel key="1" header="UTF-8 示例">
                    <a-descriptions :column="1" bordered size="small">
                      <a-descriptions-item label="文本">
                        Hello World
                      </a-descriptions-item>
                      <a-descriptions-item label="UTF-8 编码（十六进制）">
                        48 65 6c 6c 6f 20 57 6f 72 6c 64
                      </a-descriptions-item>
                      <a-descriptions-item label="中文文本">
                        你好世界
                      </a-descriptions-item>
                      <a-descriptions-item label="UTF-8 编码（十六进制）">
                        e4 bd a0 e5 a5 bd e4 b8 96 e7 95 8c
                      </a-descriptions-item>
                    </a-descriptions>
                  </a-collapse-panel>
                </a-collapse>
              </div>
              

            </div>
          </a-card>
        </div>
      </a-tab-pane>
      
      <!-- 批量处理标签页 -->
      <a-tab-pane key="batch" tab="批量处理">
        <div class="crypto-section">
          <a-card title="批量编码/哈希计算" :bordered="false">
            <div class="crypto-form">
              <a-row :gutter="16" class="mb-4">
                <a-col :span="24">
                  <a-textarea
                    v-model:value="batchInput"
                    placeholder="输入要批量处理的文本"
                    :rows="4"
                    :maxlength="10000"
                    show-count
                    allow-clear
                  />
                </a-col>
              </a-row>
              
              <a-row :gutter="16" class="mb-4">
                <a-col :span="24">
                  <a-button 
                    type="primary" 
                    block 
                    @click="handleBatchProcess"
                    :loading="batchLoading"
                    :disabled="!batchInput"
                  >
                    <template #icon>
                      <ThunderboltOutlined />
                    </template>
                    批量处理所有编码和哈希
                  </a-button>
                </a-col>
              </a-row>
              
              <div class="batch-results" v-if="batchResults">
                <a-collapse v-model:activeKey="batchActiveKeys">
                  <a-collapse-panel key="base64" header="Base64 编码">
                    <div class="result-content">
                      <code>{{ batchResults.base64 }}</code>
                      <a-button 
                        size="small" 
                        type="text" 
                        @click="copyToClipboard(batchResults.base64)"
                        class="copy-btn"
                      >
                        <CopyOutlined />
                      </a-button>
                    </div>
                  </a-collapse-panel>
                  
                  <a-collapse-panel key="utf8" header="UTF-8 编码（十六进制）">
                    <div class="result-content">
                      <code>{{ formatHex(batchResults.utf8Hex) }}</code>
                      <a-button 
                        size="small" 
                        type="text" 
                        @click="copyToClipboard(batchResults.utf8Hex)"
                        class="copy-btn"
                      >
                        <CopyOutlined />
                      </a-button>
                    </div>
                  </a-collapse-panel>
                  
                  <a-collapse-panel key="md5" header="MD5 哈希">
                    <div class="result-content">
                      <code>{{ formatHex(batchResults.md5) }}</code>
                      <a-button 
                        size="small" 
                        type="text" 
                        @click="copyToClipboard(batchResults.md5)"
                        class="copy-btn"
                      >
                        <CopyOutlined />
                      </a-button>
                    </div>
                  </a-collapse-panel>
                  
                  <a-collapse-panel key="sha1" header="SHA-1 哈希">
                    <div class="result-content">
                      <code>{{ formatHex(batchResults.sha1) }}</code>
                      <a-button 
                        size="small" 
                        type="text" 
                        @click="copyToClipboard(batchResults.sha1)"
                        class="copy-btn"
                      >
                        <CopyOutlined />
                      </a-button>
                    </div>
                  </a-collapse-panel>
                  
                  <a-collapse-panel key="sha256" header="SHA-256 哈希">
                    <div class="result-content">
                      <code>{{ formatHex(batchResults.sha256) }}</code>
                      <a-button 
                        size="small" 
                        type="text" 
                        @click="copyToClipboard(batchResults.sha256)"
                        class="copy-btn"
                      >
                        <CopyOutlined />
                      </a-button>
                    </div>
                  </a-collapse-panel>
                  
                  <a-collapse-panel key="sha512" header="SHA-512 哈希">
                    <div class="result-content">
                      <code>{{ formatHex(batchResults.sha512) }}</code>
                      <a-button 
                        size="small" 
                        type="text" 
                        @click="copyToClipboard(batchResults.sha512)"
                        class="copy-btn"
                      >
                        <CopyOutlined />
                      </a-button>
                    </div>
                  </a-collapse-panel>
                </a-collapse>
                
                <div class="batch-actions mt-4">
                  <a-space>
                    <a-button @click="copyAllBatchResults">
                      <template #icon>
                        <CopyOutlined />
                      </template>
                      复制所有结果
                    </a-button>
                    <a-button @click="clearBatchResults">
                      <template #icon>
                        <DeleteOutlined />
                      </template>
                      清除结果
                    </a-button>
                  </a-space>
                </div>
              </div>
              

            </div>
          </a-card>
        </div>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { 
  LockOutlined, 
  UnlockOutlined, 
  CalculatorOutlined,
  CopyOutlined,
  DeleteOutlined,
  FormOutlined,
  CodeOutlined,
  FileTextOutlined,
  ThunderboltOutlined
} from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import {
  base64Encode,
  base64Decode,
  md5Hash,
  sha1Hash,
  sha256Hash,
  sha512Hash,
  utf8Encode,
  utf8Decode,
  isValidBase64,
  isValidHex,
  copyToClipboard,
  computeAllEncodings,
  formatHex
} from '@/utils/cryptoUtils';

// ==================== Base64 编解码相关 ====================
const activeTab = ref<'base64' | 'hash' | 'utf8' | 'batch'>('base64');
const base64Input = ref<string>('');
const base64Loading = ref<boolean>(false);
const base64Result = ref<{
  type: 'success' | 'error' | 'info' | 'warning';
  title: string;
  content: string;
} | null>(null);

// Base64 编码
const handleBase64Encode = async () => {
  if (!base64Input.value) {
    message.warning('请输入要编码的文本');
    return;
  }
  
  base64Loading.value = true;
  try {
    const result = base64Encode(base64Input.value);
    base64Result.value = {
      type: 'success',
      title: 'Base64 编码成功',
      content: result,
    };
    message.success('Base64 编码成功');
  } catch (error) {
    base64Result.value = {
      type: 'error',
      title: 'Base64 编码失败',
      content: (error as Error).message,
    };
    message.error('Base64 编码失败');
  } finally {
    base64Loading.value = false;
  }
};

// Base64 解码
const handleBase64Decode = async () => {
  if (!base64Input.value) {
    message.warning('请输入要解码的文本');
    return;
  }
  
  if (!isValidBase64(base64Input.value)) {
    message.warning('请输入有效的 Base64 编码');
    return;
  }
  
  base64Loading.value = true;
  try {
    const result = base64Decode(base64Input.value);
    base64Result.value = {
      type: 'success',
      title: 'Base64 解码成功',
      content: result,
    };
    message.success('Base64 解码成功');
  } catch (error) {
    base64Result.value = {
      type: 'error',
      title: 'Base64 解码失败',
      content: (error as Error).message,
    };
    message.error('Base64 解码失败');
  } finally {
    base64Loading.value = false;
  }
};

// Base64 输入变化处理
const handleBase64InputChange = () => {
  if (base64Result.value) {
    base64Result.value = null;
  }
};

// 清除 Base64 结果
const clearBase64Result = () => {
  base64Result.value = null;
};

// ==================== 哈希计算相关 ====================
const hashInput = ref<string>('');
const hashLoading = ref<boolean>(false);
const hashResults = reactive<Array<{
  id: number;
  type: string;
  value: string;
  formatted: boolean;
}>>([]);
let hashResultId = 0;

// 获取哈希类型对应的颜色
const getHashColor = (type: string): string => {
  const colors: Record<string, string> = {
    'MD5': 'blue',
    'SHA-1': 'green',
    'SHA-256': 'orange',
    'SHA-512': 'red',
  };
  return colors[type] || 'default';
};

// MD5 哈希计算
const handleMd5Hash = async () => {
  if (!hashInput.value) {
    message.warning('请输入要计算哈希值的文本');
    return;
  }
  
  hashLoading.value = true;
  try {
    const result = await md5Hash(hashInput.value);
    hashResults.push({
      id: ++hashResultId,
      type: 'MD5',
      value: result,
      formatted: false,
    });
    message.success('MD5 计算成功');
  } catch (error) {
    message.error('MD5 计算失败');
  } finally {
    hashLoading.value = false;
  }
};

// SHA-1 哈希计算
const handleSha1Hash = async () => {
  if (!hashInput.value) {
    message.warning('请输入要计算哈希值的文本');
    return;
  }
  
  hashLoading.value = true;
  try {
    const result = await sha1Hash(hashInput.value);
    hashResults.push({
      id: ++hashResultId,
      type: 'SHA-1',
      value: result,
      formatted: false,
    });
    message.success('SHA-1 计算成功');
  } catch (error) {
    message.error('SHA-1 计算失败');
  } finally {
    hashLoading.value = false;
  }
};

// SHA-256 哈希计算
const handleSha256Hash = async () => {
  if (!hashInput.value) {
    message.warning('请输入要计算哈希值的文本');
    return;
  }
  
  hashLoading.value = true;
  try {
    const result = await sha256Hash(hashInput.value);
    hashResults.push({
      id: ++hashResultId,
      type: 'SHA-256',
      value: result,
      formatted: false,
    });
    message.success('SHA-256 计算成功');
  } catch (error) {
    message.error('SHA-256 计算失败');
  } finally {
    hashLoading.value = false;
  }
};

// SHA-512 哈希计算
const handleSha512Hash = async () => {
  if (!hashInput.value) {
    message.warning('请输入要计算哈希值的文本');
    return;
  }
  
  hashLoading.value = true;
  try {
    const result = await sha512Hash(hashInput.value);
    hashResults.push({
      id: ++hashResultId,
      type: 'SHA-512',
      value: result,
      formatted: false,
    });
    message.success('SHA-512 计算成功');
  } catch (error) {
    message.error('SHA-512 计算失败');
  } finally {
    hashLoading.value = false;
  }
};

// 计算所有哈希值
const handleAllHashes = async () => {
  if (!hashInput.value) {
    message.warning('请输入要计算哈希值的文本');
    return;
  }
  
  hashLoading.value = true;
  try {
    const [md5, sha1, sha256, sha512] = await Promise.all([
      md5Hash(hashInput.value),
      sha1Hash(hashInput.value),
      sha256Hash(hashInput.value),
      sha512Hash(hashInput.value),
    ]);
    
    hashResults.push(
      {
        id: ++hashResultId,
        type: 'MD5',
        value: md5,
        formatted: false,
      },
      {
        id: ++hashResultId,
        type: 'SHA-1',
        value: sha1,
        formatted: false,
      },
      {
        id: ++hashResultId,
        type: 'SHA-256',
        value: sha256,
        formatted: false,
      },
      {
        id: ++hashResultId,
        type: 'SHA-512',
        value: sha512,
        formatted: false,
      }
    );
    
    message.success('所有哈希值计算成功');
  } catch (error) {
    message.error('哈希值计算失败');
  } finally {
    hashLoading.value = false;
  }
};

// 哈希输入变化处理
const handleHashInputChange = () => {
  // 输入变化时不清除结果，用户可以继续添加
};

// 格式化哈希值
const formatHashValue = (item: any) => {
  const index = hashResults.findIndex(r => r.id === item.id);
  if (index !== -1) {
    if (hashResults[index].formatted) {
      // 如果已经格式化，则恢复原始格式
      hashResults[index].value = hashResults[index].value.replace(/\s/g, '');
    } else {
      // 如果未格式化，则进行格式化
      hashResults[index].value = formatHex(hashResults[index].value);
    }
    hashResults[index].formatted = !hashResults[index].formatted;
  }
};

// 删除哈希结果
const removeHashResult = (id: number) => {
  const index = hashResults.findIndex(r => r.id === id);
  if (index !== -1) {
    hashResults.splice(index, 1);
  }
};

// 清除所有哈希结果
const clearAllHashResults = () => {
  hashResults.splice(0, hashResults.length);
  message.info('已清除所有哈希结果');
};

// 复制所有哈希结果
const copyAllHashResults = async () => {
  if (hashResults.length === 0) {
    message.warning('没有可复制的哈希结果');
    return;
  }
  
  const text = hashResults.map(r => `${r.type}: ${r.value}`).join('\n');
  const success = await copyToClipboard(text);
  if (success) {
    message.success('已复制所有哈希结果到剪贴板');
  } else {
    message.error('复制失败');
  }
};

// ==================== UTF-8 编解码相关 ====================
const utf8Input = ref<string>('');
const utf8Loading = ref<boolean>(false);
const utf8Result = ref<{
  type: 'success' | 'error' | 'info' | 'warning';
  title: string;
  content: string;
} | null>(null);

// UTF-8 编码
const handleUtf8Encode = async () => {
  if (!utf8Input.value) {
    message.warning('请输入要编码的文本');
    return;
  }
  
  utf8Loading.value = true;
  try {
    const result = utf8Encode(utf8Input.value);
    utf8Result.value = {
      type: 'success',
      title: 'UTF-8 编码成功',
      content: result,
    };
    message.success('UTF-8 编码成功');
  } catch (error) {
    utf8Result.value = {
      type: 'error',
      title: 'UTF-8 编码失败',
      content: (error as Error).message,
    };
    message.error('UTF-8 编码失败');
  } finally {
    utf8Loading.value = false;
  }
};

// UTF-8 解码
const handleUtf8Decode = async () => {
  if (!utf8Input.value) {
    message.warning('请输入要解码的十六进制字符串');
    return;
  }
  
  if (!isValidHex(utf8Input.value)) {
    message.warning('请输入有效的十六进制字符串');
    return;
  }
  
  utf8Loading.value = true;
  try {
    const result = utf8Decode(utf8Input.value);
    utf8Result.value = {
      type: 'success',
      title: 'UTF-8 解码成功',
      content: result,
    };
    message.success('UTF-8 解码成功');
  } catch (error) {
    utf8Result.value = {
      type: 'error',
      title: 'UTF-8 解码失败',
      content: (error as Error).message,
    };
    message.error('UTF-8 解码失败');
  } finally {
    utf8Loading.value = false;
  }
};

// UTF-8 输入变化处理
const handleUtf8InputChange = () => {
  if (utf8Result.value) {
    utf8Result.value = null;
  }
};

// 格式化 UTF-8 结果
const formatUtf8Result = () => {
  if (utf8Result.value && utf8Result.value.type === 'success') {
    utf8Result.value.content = formatHex(utf8Result.value.content);
  }
};

// 清除 UTF-8 结果
const clearUtf8Result = () => {
  utf8Result.value = null;
};

// ==================== 批量处理相关 ====================
const batchInput = ref<string>('');
const batchLoading = ref<boolean>(false);
const batchResults = ref<{
  base64: string;
  utf8Hex: string;
  md5: string;
  sha1: string;
  sha256: string;
  sha512: string;
} | null>(null);
const batchActiveKeys = ref<string[]>(['base64', 'utf8', 'md5', 'sha1', 'sha256', 'sha512']);

// 批量处理
const handleBatchProcess = async () => {
  if (!batchInput.value) {
    message.warning('请输入要处理的文本');
    return;
  }
  
  batchLoading.value = true;
  try {
    const results = await computeAllEncodings(batchInput.value);
    batchResults.value = results;
    message.success('批量处理完成');
  } catch (error) {
    message.error('批量处理失败');
  } finally {
    batchLoading.value = false;
  }
};

// 复制所有批量处理结果
const copyAllBatchResults = async () => {
  if (!batchResults.value) {
    message.warning('没有可复制的批量处理结果');
    return;
  }
  
  const text = `Base64: ${batchResults.value.base64}
UTF-8 Hex: ${batchResults.value.utf8Hex}
MD5: ${batchResults.value.md5}
SHA-1: ${batchResults.value.sha1}
SHA-256: ${batchResults.value.sha256}
SHA-512: ${batchResults.value.sha512}`;
  
  const success = await copyToClipboard(text);
  if (success) {
    message.success('已复制所有批量处理结果到剪贴板');
  } else {
    message.error('复制失败');
  }
};

// 清除批量处理结果
const clearBatchResults = () => {
  batchResults.value = null;
  message.info('已清除批量处理结果');
};
</script>

<style scoped lang="scss">
.crypto-tools {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  
  .crypto-tabs {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    
    :deep(.ant-tabs-nav) {
      margin-bottom: 16px;
      flex-shrink: 0;
    }
    
    :deep(.ant-tabs-content) {
      flex: 1;
      overflow: hidden;
      
      .ant-tabs-tabpane {
        height: 100%;
        overflow: hidden;
      }
    }
  }
  
  .crypto-section {
    height: 100%;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    
    .ant-card {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
      
      .ant-card-body {
        flex: 1;
        overflow: auto;
        padding: 16px;
        
        /* 隐藏滚动条但保持滚动功能 */
        &::-webkit-scrollbar {
          width: 6px;
          height: 6px;
        }
        
        &::-webkit-scrollbar-track {
          background: #f1f1f1;
          border-radius: 3px;
        }
        
        &::-webkit-scrollbar-thumb {
          background: #c1c1c1;
          border-radius: 3px;
          
          &:hover {
            background: #a8a8a8;
          }
        }
      }
    }
    
    .crypto-form {
      .crypto-result {
        min-height: 60px;
        margin-bottom: 16px;
      }
      
      .hash-results {
        .hash-title {
          display: flex;
          align-items: center;
          justify-content: space-between;
          
          .hash-type {
            font-weight: bold;
            margin-right: 8px;
          }
        }
        
        .hash-value {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 4px;
          
          code {
            flex: 1;
            word-break: break-all;
            font-family: 'Consolas', 'Monaco', monospace;
            background: #f5f5f5;
            padding: 2px 4px;
            border-radius: 2px;
          }
          
          .hash-actions {
            margin-left: 8px;
            flex-shrink: 0;
          }
        }
        
        .hash-length {
          font-size: 12px;
          color: #999;
          margin-top: 4px;
        }
        
        .hash-actions-bottom {
          display: flex;
          justify-content: center;
        }
      }
      
      .batch-results {
        .result-content {
          display: flex;
          align-items: center;
          justify-content: space-between;
          
          code {
            flex: 1;
            word-break: break-all;
            font-family: 'Consolas', 'Monaco', monospace;
            background: #f5f5f5;
            padding: 4px 8px;
            border-radius: 4px;
            margin-right: 8px;
          }
          
          .copy-btn {
            flex-shrink: 0;
          }
        }
        
        .batch-actions {
          display: flex;
          justify-content: center;
        }
      }
      
      .utf8-examples {
        margin-bottom: 16px;
      }
      
      .crypto-info {
        margin-top: 16px;
      }
    }
  }
}
</style>