import axios from "axios";
import {ElMessage} from "element-plus";
import {useStore} from "@/stores";

export const useAxios = axios.create({
    baseURL: "", //
    timeout: 3000,
})

export interface baseResponse<T> {
    code: number
    data: T
    msg: string
}

export interface listResponse<T> {
    count: number
    list: T[]
}

export interface paramsType{
    limit?: number
    page?: number
    key?:string
}


// 请求拦截器
useAxios.interceptors.request.use((config) => {
    const store = useStore()
    config.headers["token"] = store.userInfo.token
    return config
})

// 响应拦截器
useAxios.interceptors.response.use((response) => {
    if (response.status !== 200) {
        // 失败的
        console.log("服务失败", response.status)
        ElMessage.error(response.statusText)
        return Promise.reject(response.statusText)
    }
    return response.data
}, (err) => {
    console.log("服务错误", err)
    ElMessage.error(err.message)
    return Promise.reject(err.message)
})


interface cacheType<T> {
    time: number
    res: Promise<T>
}

// 60000 默认时间是一分钟
export function cacheRequest<T, K>(func: (params: K) => Promise<T>, timeout = 60000): (params: K) => Promise<T> {
    const cacheData: Record<string, cacheType<T>> = {}

    return (params: K) => {

        const val = cacheData[JSON.stringify(params)]
        if (val) {
            if ((new Date().getTime() - val.time) < timeout) {
                return val.res
            }
        }
        const res: Promise<T> = func(params)
        // 请求后
        cacheData[JSON.stringify(params)] = {
            time: new Date().getTime(),
            res: res
        }
        return res
    }
}