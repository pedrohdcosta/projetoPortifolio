import axios from "axios";
import { useAuth } from '../stores/auth'

const api = axios.create({baseURL: '/api'})

api.interceptors.request.use((config) => {
    const a = useAuth()
    if (a.token) config.headers.Authorization = `Bearer ${a.token}`
    return config
})

export default api