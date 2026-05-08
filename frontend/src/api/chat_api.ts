import {type baseResponse, type listResponse, type paramsType, useAxios} from "@/api/index";
import type {msgType} from "@/types/msg";
export interface chatSessionType {
    "userId": number
    "avatar": string
    "nickname": string
    "created_at": string
    "msgPreview": string
    "isTop": boolean
}

// 会话列表
export function chatSessionListApi(): Promise<baseResponse<listResponse<chatSessionType>>> {
    return useAxios.get("/api/chat/session", {params: {limit: -1}})
}

export interface userBaseType {
    "id": number
    "nickName": string
    "avatar": string
}



export interface chatHistoryListParams extends paramsType {
    friendId: number
}

export interface chatHistoryType {
    "id": number
    "sendUser": userBaseType
    "revUser": userBaseType
    "isMe": boolean
    "created_at": string
    "msg": msgType
    "systemMsg": null
    msgPreview: string
    showDate: boolean
}

export function chatHistoryListApi(params: chatHistoryListParams): Promise<baseResponse<listResponse<chatHistoryType>>> {
    return useAxios.get("/api/chat/history", {params})
}

// 删除聊天记录
export function chatDeleteApi(idList: number[]): Promise<baseResponse<string>> {
    return useAxios.delete("/api/chat/chat", {data: {idList}})
}

export function chatUserTopApi(friendId: number): Promise<baseResponse<string>> {
    return useAxios.post("/api/chat/user_top", {friendId})
}