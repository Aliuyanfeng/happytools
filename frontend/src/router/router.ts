/*
 * @Author: LiuYanFeng
 * @Date: 2025-07-04 17:51:30
 * @LastEditors: LiuYanFeng
 * @LastEditTime: 2025-07-04 17:55:34
 * @FilePath: \happytools\frontend\src\router\router.ts
 * @Description: 像珍惜礼物一样珍惜今天
 * 
 * Copyright (c) 2025 by ${git_name_email}, All Rights Reserved. 
 */
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
const HomeView = () => import('../views/HomeView.vue');
const HelloWorld = () => import('../components/HelloWorld.vue');
const VirusTotal = () => import('../views/VirusTotal/VirusTotal.vue');

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home',
    meta: {
      title: '首页'
    },
    component: HomeView
  },
  {
    path: '/toolbox',
    name: 'toolbox',
    meta: {
      title: '工具盒子'
    },
    component: () => import('../views/Toolbox/Toolbox.vue')
  },
  {
    path: '/vt',
    name: 'vt',
    meta: {
      title: 'VirusTotal'
    },
    component: VirusTotal
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;