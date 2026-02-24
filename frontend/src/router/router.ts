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
import {createRouter, createWebHistory, RouteRecordRaw} from 'vue-router';
import routes from "@/router/routes";


const router = createRouter({
    history: createWebHistory(),
    routes
});

export default router;