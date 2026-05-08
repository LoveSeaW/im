import {type baseResponse, cacheRequest, type listResponse, type paramsType, useAxios} from "@/api/index";

// 用户验证
export interface verificationQuestionType {
    problem1?: string
    problem2?: string
    problem3?: string
    answer1?: string
    answer2?: string
    answer3?: string
}

// 用户配置信息
export interface userConfType {
    "userID": number
    "nickname": string
    "abstract": string
    "avatar": string
    "recallMessage"?: string // 撤回消息的提示词
    "friendOnline": boolean // 好友上线提醒
    "sound": boolean // 声音
    "secureLink": boolean // 安全连接
    "savePwd": boolean // 报存密码
    "searchUser": number
    "verification": number
    "verificationQuestion"?: verificationQuestionType
}

// 获取用户配置信息
export function userInfoApi(): Promise<baseResponse<userConfType>> {
    return useAxios.get("/api/user/user_info")
}

export interface userConfUpdateType {
    "nickname"?: string
    "abstract"?: string
    "avatar"?: string
    "recallMessage"?: string // 撤回消息的提示词
    "friendOnline"?: boolean
    "sound"?: boolean
    "secureLink"?: boolean
    "savePwd"?: boolean
    "searchUser"?: number
    "verification"?: number
    "verificationQuestion"?: verificationQuestionType
}


// 更新用户信息
export function userInfoUpdateApi(data: userConfUpdateType): Promise<baseResponse<string>> {
    return useAxios.put("/api/user/user_info", data)
}

export interface friendType {
    "userID": number
    "nickname": string
    "abstract": string
    "avatar": string
    "notice": string
    "isOnline": boolean
}

// 好友列表
export function friendListApi(params: paramsType): Promise<baseResponse<listResponse<friendType>>> {
    return useAxios.get("/api/user/friends", {params})
}

const friendDetailCache = {
    friendID: 0,
    time: 0,
    res: {}
}

// 好友详情
// export function friendDetailApi(friendID: number): Promise<baseResponse<friendType>> {
//     // 请求前
//     if (friendID === friendDetailCache.friendID) {
//         if ((new Date().getTime() - friendDetailCache.time) / 1000 < 60) {
//             return friendDetailCache.res as Promise<baseResponse<friendType>>
//         }
//     }
//
//     const res: Promise<baseResponse<friendType>> = useAxios.get("/api/user/friend_info", {params: {friendID}})
//     // 请求后
//     friendDetailCache.friendID = friendID
//     friendDetailCache.time = new Date().getTime()
//     friendDetailCache.res = res
//     return res
// }



interface friendDetailRequest {
    friendID: number
}

export const friendDetailApi = cacheRequest<baseResponse<friendType>, friendDetailRequest>(
    (params: friendDetailRequest) => useAxios.get("/api/user/friend_info", {params}),
    10000,
)

export interface userSearchResponse {
    "userID": number
    "nickname": string
    "abstract": string
    "avatar": string
    "isFriend": boolean
}

// 用户搜索
export function userSearchApi(params: paramsType): Promise<baseResponse<listResponse<userSearchResponse>>> {
    return useAxios.get("/api/user/search", {params})
}

export interface validResponse {
    "verification": number // 0 不允许任何人添加  1 允许任何人添加  2 需要验证消息 3 需要回答问题  4  需要正确回答问题
    "verificationQuestion": verificationQuestionType
}

// 好友验证
export function userValidApi(friendId: number): Promise<baseResponse<validResponse>> {
    return useAxios.post("/api/user/valid", {friendId})
}

export interface friendAddRequest {
    "friendId": number // 好友ID
    "verificationQuestion": verificationQuestionType
    "verify": string // 验证消息
}

// 添加好友
export function friendAddApi(data: friendAddRequest): Promise<baseResponse<string>> {
    return useAxios.post("/api/user/friends", data)
}

export interface userValidType {
    "userID": number
    "nickname": string
    "avatar": string
    "additionalMessages": string // 附加消息
    "verificationQuestion"?: verificationQuestionType
    "status": number // 0 未操作 1 同意 2 拒绝 3 忽略 4 删除
    "verification": number
    "id": number
    revStatus: number
    sendStatus: number
    "flag": "send" | "rev"
    "createdAt": string
}

// 用户验证列表
export function userValidListApi(params?: paramsType): Promise<baseResponse<listResponse<userValidType>>> {
    return useAxios.get("/api/user/valid", {params})
}

export interface userValidStatusRequest {
    "verifyId": number
    "status": number
}

// 修改用户验证状态
export function userValidStatusApi(data: userValidStatusRequest): Promise<baseResponse<string>> {
    return useAxios.put("/api/user/valid_status", data)
}

// 修改用户昵称
export function friendUpdateNicknameApi(friendID: number, notice: string): Promise<baseResponse<string>> {
    return useAxios.put("/api/user/friends", {friendID, notice})
}

// 用户基本信息
export interface userBaseInfoType {
    "userId": number
    "avatart": string // 用户头像 防伪标记
    "nickname": string
}

export function friendRemoveApi(friendId: number): Promise<baseResponse<string>> {
    return useAxios.delete("/api/user/friends", {
        data: {
            "friendId": friendId,
        }
    })
}