import { createRouter, createWebHistory } from 'vue-router'
import Login from '../pages/Login.vue'
import Register from '../pages/Register.vue'
import Dashboard from '../pages/Dashboard.vue'
import Devices from '../pages/Devices.vue'
import Profile from '../pages/Profile.vue'
import { useAuth } from '../stores/auth'

const routes = [
    {path: '/login', component: Login},
    {path: '/register', component: Register},
    {path: '/', redirect: '/app/dashboard'},
    {
        path: '/app',
        children: [
            {path: 'dashboard', component: Dashboard},
            {path: 'devices', component: Devices},
            {path: 'profile', component: Profile}
        ],
        beforeEnter: () => { const a = useAuth(); return a.isAuthenticated ? true: '/login' }
    }
]

export default createRouter({ history: createWebHistory(), routes})