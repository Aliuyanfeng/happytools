<template>
  <div class="toolbox">
    <a-tabs v-model:activeKey="activeKey" type="card">
      <a-tab-pane key="1">
        <template #tab>
        <span class="flex items-center">
          <RetweetOutlined/>
          数据转换
        </span>
        </template>
        Content of Tab Pane 1
      </a-tab-pane>
      <a-tab-pane key="2">
        <template #tab>
        <span class="flex items-center">
          <GatewayOutlined/>
          网络调试
        </span>
        </template>

        Content of Tab Pane 2
      </a-tab-pane>
      <a-tab-pane key="3">
        <template #tab>
        <span class="flex items-center">
          <ToolOutlined />
          辅助工具
        </span>
        </template>
        Content of Tab Pane 3
      </a-tab-pane>
    </a-tabs>
    <a-flex class="toolbox_content" wrap="wrap" gap="middle" align="center" justify="center">
      <a-card hoverable style="width: 30%" v-for="item in 5">
        <template #actions>
          <setting-outlined key="setting"/>
          <edit-outlined key="edit"/>
          <ellipsis-outlined key="ellipsis"/>
        </template>
        <a-card-meta title="单位转换"
                     description="快速转换常见数据单位，包括存储大小（B、KB、MB、GB）、时间（秒、分钟、小时）、速率（bps、KBps、MBps）等，支持自定义格式与实时换算">
          <template #avatar>
            <span class="text-3xl font-bold leading-4">D</span>
          </template>
        </a-card-meta>
      </a-card>
    </a-flex>
  </div>
</template>

<script setup lang="ts">
import {EditOutlined, EllipsisOutlined, RetweetOutlined, SettingOutlined, GatewayOutlined, ToolOutlined} from '@ant-design/icons-vue';
import {ref} from "vue";

const activeKey = ref('1');
</script>

<style scoped>
.toolbox {
  width: 100%;
  margin: 0 auto;
  padding: 16px;

  .toolbox_content:after {
    content: '';
    width: 30%;
  }
}
</style>