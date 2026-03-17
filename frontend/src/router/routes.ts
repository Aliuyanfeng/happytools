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
                meta: {
                    title: '单位转换'
                },
                component: () => import('../views/Toolbox/UnitConverter.vue')
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
                meta: {
                    title: '批量重命名'
                },
                component: () => import('../views/Toolbox/BatchRename.vue')
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
    }
];

export default routes;