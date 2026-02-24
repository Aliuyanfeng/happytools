import {RouteRecordRaw} from "vue-router";

const Home = () => import('../views/Home.vue');
const Dashboard = () => import('../views/Dashboard.vue');
const VirusTotal = () => import('../views/VirusTotal/VirusTotal.vue');
const Todo = () => import('../views/Todo/Todo.vue');
const NetworkDebug = () => import('../views/Network/NetworkDebug.vue');

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
        component: () => import('../views/Toolbox/Toolbox.vue')
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
    }
];

export default routes;