/*
 * @Author: LiuYanFeng
 * @Date: 2025-07-03 17:16:49
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2026-02-25 10:21:24
 * @FilePath: \happytools\frontend\src\main.ts
 * @Description: 像珍惜礼物一样珍惜今天
 * 
 * Copyright (c) 2026 by ${git_name_email}, All Rights Reserved. 
 */
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './style.css' //add
import 'ant-design-vue/dist/reset.css';
import router from './router/router';
import i18n from './locales';
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
dayjs.locale('zh-cn');

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(i18n)
app.mount('#app')
