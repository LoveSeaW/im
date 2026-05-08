import type {baseResponse} from "@/api/index";
import {useAxios} from "@/api/index";

export interface authLoginRequest {
    "userName": string
    "password": string
}

export interface authLoginResponse {
    token: string
}

// 登陆
export function authLoginApi(data: authLoginRequest): Promise<baseResponse<authLoginResponse>> {
    return useAxios.post("/api/auth/login", data)
}


export interface registerRequest {
    "nickname": string
    "pwd": string
    "rePwd": string
}

export interface registerResponse {
    userID: string
}

// 登陆
export function registerApi(data: registerRequest): Promise<baseResponse<registerResponse>> {
    return useAxios.post("/api/auth/register", data)
}

interface authOpenLoginRequest {
    "code": string
    "flag": string
}

// 第三方登陆
export function authOpenLoginApi(data: authOpenLoginRequest): Promise<baseResponse<authLoginResponse>> {
    return useAxios.post("/api/auth/open_login", data)
}


// 注销
export function logoutApi(): Promise<baseResponse<string>> {
    return useAxios.post("/api/auth/logout")
}