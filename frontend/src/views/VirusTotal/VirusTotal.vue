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
<!--      <a-table :dataSource="dataSource" :columns="columns" rowKey="key" :expandedRowKeys="expandedRowKeys"-->
<!--               :expandable="{-->
<!--                  expandedRowRender: expandedRowRender,-->
<!--                  expandIcon: () => null,      // ✅ 彻底隐藏默认展开图标-->
<!--                }">-->
<!--        <template #bodyCell="{ column, record }">-->
<!--          <template v-if="column.key === 'resource'">-->
<!--            <span>-->
<!--              <FileFilled style="color: #cfd9df" v-if="record.type === 'file'"/>-->
<!--              <FolderFilled style="color: #FFD04B" v-else/>-->
<!--              {{ record.resource }}-->
<!--            </span>-->
<!--          </template>-->
<!--          <template v-if="column.key === 'action'">-->
<!--            <span>-->
<!--              <a>展开</a>-->
<!--              <a-divider type="vertical"/>-->
<!--              <a>删除</a>-->
<!--            </span>-->
<!--          </template>-->
<!--        </template>-->
<!--        <template #expandedRowRender="{ record }">-->
<!--          <p style="margin: 0">-->
<!--            {{ record.description }}-->
<!--          </p>-->
<!--        </template>-->

<!--      </a-table>-->
    </a-card>
    <!--    使用说明告知弹窗-->
    <a-modal
        v-model:open="useModalVisible"
        centered
        :maskClosable="false"
        :keyboard="false"
        :closable="false"
        :footer="null"
    >
      <template #title>
        <div class="flex items-center">
          <ExclamationCircleTwoTone class="mr-2"/>
          请确认您的 API 密钥来源
        </div>
      </template>
      <div class="use-instruction" @scroll="handleScroll">
        <p>请确保您拥有合法有效的 VirusTotal API 密钥，并遵守相关服务条款。</p>
      </div>
      <div class="flex justify-between">
        <a-checkbox v-model:checked="isRemember" class="flex items-center" @change="isRememberChange">不再弹出
        </a-checkbox>
        <a-button type="primary" @click="useModalClick">
          我已了解并确认
        </a-button>
      </div>
    </a-modal>
    <!--创建任务弹窗-->
    <a-modal
        v-model:open="createModalVisible"
        centered
        :maskClosable="false"
        :keyboard="false"
        :closable="false"
    >
      <template #title>
        <div class="flex items-center">
          <ExclamationCircleTwoTone class="mr-2"/>
          创建任务
        </div>
      </template>
      <a-form
          :model="formState"
          name="basic"
          autocomplete="off"
          @finish="onFinish"
          @finishFailed="onFinishFailed"
      >
        <a-form-item
            label="API Key"
            name="API Key"
            :rules="[{ required: true, message: '请输入您的API Key!' }]"
        >
          <a-input-password v-model:value="formState.apiKey"/>
        </a-form-item>

        <a-form-item label="操作目标" name="resource" required>
          <a-radio-group v-model:value="formState.resource">
            <a-radio value="1">单个文件</a-radio>
            <a-radio value="2">批量目录</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item
            label="样本位置"
            name="样本位置"
            v-if="formState.resource === '1'"
            :rules="[{ required: true, message: 'Please input your samplePath!' }]"
        >
          <div class="flex items-center justify-between">
            <a-input v-model:value="formState.samplePath" readonly v-if="formState.samplePath!=''"/>
          </div>
          <div>
            <a-button type="primary" :icon="h(FileAddOutlined)" @click="ChooseFile">
              {{ formState.samplePath ? '重新选择文件' : '选择文件' }}
            </a-button>
          </div>
        </a-form-item>

        <a-form-item
            label="目录位置"
            name="目录位置"
            v-if="formState.resource === '2'"
            :rules="[{ required: true, message: 'Please input your directoryPath!' }]"
        >
          <div class="flex items-center justify-between">
            <a-input v-model:value="formState.directoryPath" readonly v-if="formState.directoryPath!=''"/>
          </div>
          <div>
            <a-button type="primary" :icon="h(FileAddOutlined)" @click="ChooseDirectory">
              {{ formState.directoryPath ? '重新选择目录' : '选择目录' }}
            </a-button>
          </div>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button key="back" v-if="!createLoading" @click="createCancel">取消</a-button>
        <a-button type="primary" key="submit" :icon="h(HeatMapOutlined)" :loading="createLoading" @click="createOk">
          提交任务
        </a-button>
      </template>
    </a-modal>
  </div>

</template>
<script lang="ts" setup>
import {
  ExclamationCircleTwoTone,
  FileAddOutlined,
  FileFilled,
  FolderFilled,
  HeatMapOutlined,
  PlusOutlined
} from '@ant-design/icons-vue';
import {h, onMounted, onUnmounted, reactive, ref} from 'vue';
import {ChangeEvent} from "ant-design-vue/es/_util/EventInterface";
import {VTService} from "../../../bindings/github.com/Aliuyanfeng/happytools/backend/services/vt";

const visible = ref<boolean>(true);
const handleClose = () => {
  visible.value = false;
};
// 确认按钮是否禁用
const isAgreeDisabled = ref<boolean>(true);
// 记住我
const isRemember = ref<boolean>(false);


onMounted(() => {
  const remember = window.sessionStorage.getItem('remember');
  useModalVisible.value = remember !== 'true';
})
onUnmounted(() => {
  // window.sessionStorage.setItem('remember', "false");
})
//监听div类名use-instruction滚动条滚动事件
const handleScroll = () => {
  const useInstruction = document.querySelector('.use-instruction');
  const {scrollTop, scrollHeight, clientHeight} = useInstruction as HTMLElement;
  if (scrollTop !== undefined && scrollHeight !== undefined && clientHeight !== undefined) {
    if (scrollTop > 300) {
      isAgreeDisabled.value = false;
    }
  }
}
// 记住我改变事件
const isRememberChange = (e: ChangeEvent) => {
  const target = e.target as HTMLInputElement;
  isRemember.value = target.checked;
}
// 使用说明弹窗是否显示
const useModalVisible = ref<boolean>(false);

const useModalClick = () => {
  useModalVisible.value = false;
  if (isRemember.value) {
    window.sessionStorage.setItem('remember', 'true');
  }
}
// 创建任务弹窗是否显示
const createModalVisible = ref<boolean>(false);
// 创建任务弹窗确认按钮是否禁用
const createLoading = ref<boolean>(false);

// 打开创建任务弹窗
const openCreateTaskModal = () => {
  createModalVisible.value = true;
}
// 创建任务弹窗确认按钮点击事件
const createOk = () => {
  createLoading.value = true;
  setTimeout(() => {
    createLoading.value = false;
    createModalVisible.value = false;
  }, 2000);
};
// 创建任务弹窗取消按钮点击事件
const createCancel = () => {
  createModalVisible.value = false;
};

interface FormState {
  apiKey: string;
  samplePath: string;
  directoryPath: string;
  remember: boolean;
  resource: string
}

const formState = reactive<FormState>({
  apiKey: '',
  samplePath: '',
  directoryPath: '',
  remember: true,
  resource: "1"
});
const onFinish = (values: any) => {
  console.log('Success:', values);
};

const onFinishFailed = (errorInfo: any) => {
  console.log('Failed:', errorInfo);
};
// 选择文件按钮点击事件
const ChooseFile = async () => {
  console.log("ChooseFile test");
  try {
    formState.samplePath = await VTService.OpenFileDialog();
  } catch (error) {
    console.error("Failed to open file dialog:", error);
  }
}
// 选择目录按钮点击事件
const ChooseDirectory = async () => {
  console.log("ChooseDirectory test");
  try {
    formState.directoryPath = await VTService.OpenFileDialogs();
  } catch (error) {
    console.error("Failed to open file dialog:", error);
  }
}

// 任务列表数据
const dataSource = [
  {
    key: 1,
    name: '文档测试',
    fileCount: 100,
    resource: 'D:\\\\AliuProject\\\\happytools\\\\backend\\\\services\\\\vt\\\\test.txt',
    type: 'file'
  },
  {
    key: 2,
    name: '目录测试',
    fileCount: 200,
    resource: '/home/admin/test',
    type: 'directory'
  },
]
// 任务列表列
const columns = [
  {
    title: '任务名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '文件数量',
    dataIndex: 'fileCount',
    key: 'fileCount',
  },
  {
    title: '任务对象',
    dataIndex: 'resource',
    key: 'resource',
  },
  {
    title: '操作',
    key: 'action',
  },
]
</script>

<style scoped lang="postcss">
.use-instruction {
  max-height: 280px;
  overflow-y: scroll;
}
</style>

