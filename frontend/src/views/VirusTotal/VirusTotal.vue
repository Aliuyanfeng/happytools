<template>
  <div class="vt p-6">
    <a-card class="mt-5">
      <a-form
          :model="formState"
          name="basic"
          :label-col="{ span: 2 }"
          :wrapper-col="{ span: 24 }"
          autocomplete="off"
          @finish="onFinish"
          @finishFailed="onFinishFailed"
      >
        <a-form-item
            label="API Key"
            name="password"
            :rules="[{ required: true, message: '请输入您的API Key!' }]"
        >
          <a-input-password v-model:value="formState.password"/>
        </a-form-item>

        <a-form-item label="Resources" name="resource">
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
          <a-input v-model:value="formState.samplePath"/>
        </a-form-item>

        <a-form-item
            label="目录位置"
            name="目录位置"
            v-if="formState.resource === '2'"
            :rules="[{ required: true, message: 'Please input your directoryPath!' }]"
        >
          <a-input v-model:value="formState.directoryPath"/>
        </a-form-item>

        <a-form-item name="remember" :wrapper-col="{ offset: 8, span: 16 }">
          <a-checkbox v-model:checked="formState.remember">Remember me</a-checkbox>
        </a-form-item>

        <a-form-item :wrapper-col="{ offset: 8, span: 16 }">
          <a-button type="primary" html-type="submit">Submit</a-button>
        </a-form-item>
      </a-form>
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
        <p>在使用本功能前，请明确您所使用的 VirusTotal API 密钥类型：</p>
        <ul>
          <li>公共 API</li>
          <li class="list-disc ml-8">免费获取，使用频率受限</li>
          <li class="list-disc ml-8">适用于个人、非商业用途</li>
          <li class="list-disc ml-8">每日请求次数有限制</li>
          <li class="list-disc ml-8">请确保您的使用场景符合 VirusTotal 公共 API 的服务条款</li>
        </ul>
        <ul>
          <li>高级 API</li>
          <li class="list-disc ml-8">需付费订阅，提供更高的请求限额</li>
          <li class="list-disc ml-8">适用于商业或高频使用场景</li>
          <li class="list-disc ml-8">提供更全面的数据访问权限</li>
          <li class="list-disc ml-8">请确保您的订阅状态有效且使用范围符合许可协议</li>
        </ul>
        <p>本软件仅提供与 VirusTotal 服务的连接功能，方便您上传文件进行恶意性检测。我们不对以下情况承担责任：</p>
        <ul>
          <li class="list-disc ml-8">API 密钥的授权使用</li>
          <li class="list-disc ml-8">检测结果的准确性</li>
          <li class="list-disc ml-8">因使用本功能导致的任何后果</li>
          <li class="list-disc ml-8">违反 VirusTotal 服务条款的行为</li>
        </ul>
        <ul>
          <li>请您确保：</li>
          <li class="list-disc ml-8">拥有合法有效的 API 密钥</li>
          <li class="list-disc ml-8">使用场景符合所选 API 类型的限制和约束</li>
          <li class="list-disc ml-8">遵守 VirusTotal 的服务条款和隐私政策</li>
          <li class="list-disc ml-8">对上传文件的内容和检测结果负责</li>
          <li class="list-disc ml-8">对使用本功能产生的任何损失或损害，不承担责任</li>
        </ul>
        <p>开始使用前，请确认您已了解并同意以上条款。</p>
      </div>
      <div class="flex justify-between">
        <a-checkbox v-model:checked="isRemember" class="flex items-center" @change="isRememberChange">不再弹出
        </a-checkbox>
        <a-button type="primary" :disabled="isAgreeDisabled" @click="useModalClick">
          我已了解并确认
        </a-button>
      </div>
    </a-modal>

  </div>

</template>
<script lang="ts" setup>
import {ExclamationCircleTwoTone} from '@ant-design/icons-vue';
import {onMounted, onUnmounted, reactive, ref} from 'vue';
import {ChangeEvent} from "ant-design-vue/es/_util/EventInterface";

const visible = ref<boolean>(true);
const handleClose = () => {
  visible.value = false;
};
// 确认按钮是否禁用
const isAgreeDisabled = ref<boolean>(true);
// 记住我
const isRemember = ref<boolean>(false);


onMounted(() => {
  const remember = window.localStorage.getItem('remember');
  useModalVisible.value = remember !== 'true';
})
onUnmounted(() => {
  // window.localStorage.setItem('remember', "false");
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

const useModalVisible = ref<boolean>(false);

const useModalClick = () => {
  useModalVisible.value = false;
  if (isRemember.value) {
    window.localStorage.setItem('remember', 'true');
  }
}

interface FormState {
  samplePath: string;
  directoryPath: string;
  password: string;
  remember: boolean;
  resource: string
}

const formState = reactive<FormState>({
  samplePath: '',
  directoryPath: '',
  password: '',
  remember: true,
  resource: "1"
});
const onFinish = (values: any) => {
  console.log('Success:', values);
};

const onFinishFailed = (errorInfo: any) => {
  console.log('Failed:', errorInfo);
};
</script>

<style scoped lang="postcss">
.use-instruction {
  max-height: 280px;
  overflow-y: scroll;
}
</style>

