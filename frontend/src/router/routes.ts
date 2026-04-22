import {RouteRecordRaw} from "vue-router";

const Home = () => import('../views/Home.vue');
const Dashboard = () => import('../views/Dashboard.vue');
const VirusTotal = () => import('../views/VirusTotal/VirusTotal.vue');
const Todo = () => import('../views/Todo/Todo.vue');
const NetworkDebug = () => import('../views/Network/NetworkDebug.vue');
const DailyReport = () => import('../views/DailyReport/DailyReport.vue');

const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        name: 'home',
        meta: {
            title: '首页'
        },
        component: Home
    },
    {
        path: '/dashboard',
        name: 'dashboard',
        meta: {
            title: '仪表盘'
        },
        component: Dashboard
    },
    {
        path: '/toolbox',
        name: 'toolbox',
        meta: {
            title: '工具盒子'
        },
        component: () => import('../views/Toolbox/Toolbox.vue'),
        children: [
            {
                path: 'unit-converter',
                name: 'unit-converter',
                meta: { title: '单位转换' },
                component: () => import('../views/Toolbox/UnitConverter.vue')
            },
            {
                path: 'base-converter',
                name: 'base-converter',
                meta: { title: '进制转换' },
                component: () => import('../views/Toolbox/BaseConverter.vue')
            },
            {
                path: 'timestamp',
                name: 'timestamp',
                meta: { title: '时间戳转换' },
                component: () => import('../views/Toolbox/Timestamp.vue')
            },
            {
                path: 'color-converter',
                name: 'color-converter',
                meta: { title: '颜色转换' },
                component: () => import('../views/Toolbox/ColorConverter.vue')
            },
            {
                path: 'encryption',
                name: 'encryption',
                meta: {
                    title: '加密工具'
                },
                component: () => import('../views/Toolbox/Encryption.vue')
            },
            {
                path: 'png-injector',
                name: 'png-injector',
                meta: {
                    title: 'PNG 注入'
                },
                component: () => import('../views/Toolbox/PNGInjector.vue')
            },
            {
                path: 'batch-rename',
                name: 'batch-rename',
                meta: { title: '批量重命名' },
                component: () => import('../views/Toolbox/BatchRename.vue')
            },
            {
                path: 'ncm-converter',
                name: 'ncm-converter',
                meta: { title: 'NCM 转 MP3' },
                component: () => import('../views/Toolbox/NcmConverter.vue')
            },
            {
                path: 'checksum',
                name: 'checksum',
                meta: { title: '校验和计算' },
                component: () => import('../views/Toolbox/Checksum.vue')
            }
        ]
    },
    {
        path: '/vt',
        name: 'vt',
        meta: {
            title: 'VirusTotal'
        },
        component: VirusTotal
    },
    {
        path: '/todo',
        name: 'todo',
        meta: {
            title: 'TODO'
        },
        component: Todo
    },
    {
        path: '/network',
        name: 'network',
        meta: {
            title: '网络调试'
        },
        component: NetworkDebug
    },
    {
        path: '/dailyReport',
        name: 'dailyReport',
        meta: {
            title: '日报管理'
        },
        component: DailyReport
    },
    {
        path: '/git-config',
        name: 'git-config',
        meta: {
            title: 'Git 配置管理'
        },
        component: () => import('../views/GitConfig/GitConfig.vue')
    },
    {
        path: '/makefile-editor',
        name: 'makefile-editor',
        meta: { title: 'Makefile 编辑器' },
        component: () => import('../views/MakefileEditor/MakefileEditor.vue')
    },
    {
        path: '/nuclei-parser',
        name: 'nuclei-parser',
        meta: { title: 'POC 模板解析' },
        component: () => import('../views/NucleiParser/NucleiParser.vue')
    }
];

export default routes;