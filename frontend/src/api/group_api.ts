import type {baseResponse, listResponse, paramsType} from "@/api/index";
import {useAxios} from "@/api/index";
import type {userBaseInfoType, userValidStatusRequest, validResponse, verificationQuestionType} from "@/api/user_api";
import type {msgType} from "@/types/msg";

export interface groupType {
    "groupId": number
    "groupTitle": string
    "groupAvatar": string
    "groupMemberCount": number
    "role": number
    "mode": number
}

// 我的群聊列表
export function groupMyListApi(mode: number): Promise<baseResponse<listResponse<groupType>>> {
    return useAxios.get("/api/group/my", {params: {limit: -1, mode: mode}})
}

export interface groupSearchResponse {
    "groupId": number
    "title": string
    "abstract": string
    "avatar": string
    "isInGroup": boolean
    "userCount": number
    "userOnlineCount": number
}

// 群聊搜索
export function groupSearchApi(params: paramsType): Promise<baseResponse<listResponse<groupSearchResponse>>> {
    return useAxios.get("/api/group/search", {params})
}

// 加群的验证信息
export function groupValidApi(id: number): Promise<baseResponse<validResponse>> {
    return useAxios.get("/api/group/valid/" + id.toString())
}

export interface groupAddRequest {
    "groupID": number
    "verificationQuestion": verificationQuestionType
    "verify": string
}

// 申请加群
export function groupAddApi(data: groupAddRequest): Promise<baseResponse<string>> {
    return useAxios.post("/api/group/valid", data)
}

export interface groupValidType {
    "id": number
    "grouId": number
    "userId": number // 申请人id
    "userNickname": string
    "userAvatar": string
    "status": number
    "additionalMessages": string
    "verificationQuestion"?: verificationQuestionType
    "title": string
    "createdAt": string
    "type": number // 1 申请加群  2 退群
}

// 群验证列表
export function groupValidListApi(params?: paramsType): Promise<baseResponse<listResponse<groupValidType>>> {
    return useAxios.get("/api/group/valid", {params})
}

export interface groupValidStatusRequest {
    "validId": number
    "status": number
}

// 修改群验证状态
export function groupValidStatusApi(data: groupValidStatusRequest): Promise<baseResponse<string>> {
    return useAxios.put("/api/group/valid/status", data)
}


export interface groupInfoType {
    "groupId": number
    "title": string
    "abstract": string
    "memberCount": number
    "memberOnlineCount": number
    "avatar": string
    "creator": userBaseInfoType
    "adminList": userBaseInfoType[]
    "role": number
    "isProhibition": boolean //是否全员禁言
    "prohibitionTime"?: number
    "isSearch": boolean
    "isInvite": boolean
    "isTemporarySession": boolean
}

// 群信息
export function groupDetailApi(id: number): Promise<baseResponse<groupInfoType>> {
    return useAxios.get("/api/group/group/" + id.toString())
}

export interface groupCreateRequest {
    "mode": 1 | 2  // 模式 1 直接创建  2 选择创建
    "name": string // 群名
    "size": number // 群规模
    isSearch: boolean // 是否开启搜索
    userIdList: number[] // 选择的用户id列表
}

// 创建群
export function groupCreateApi(data: groupCreateRequest): Promise<baseResponse<string>> {
    return useAxios.post("/api/group/group", data)
}

export interface groupFriendResponse {
    "userId": number
    "avatar": string
    "nickname": string
    "isInGroup": boolean
}

// 我的哪些好友在这个群里面
export function groupMyFriendListApi(id: number): Promise<baseResponse<listResponse<groupFriendResponse>>> {
    return useAxios.get("/api/group/friends", {params: {id, limit: -1}})
}

export interface groupMemberRequest {
    "id": number
    "memberIdList": number[]
}

export function groupMemberAddApi(data: groupMemberRequest): Promise<baseResponse<string>> {
    return useAxios.post("/api/group/member", data)
}

export interface groupUpdateRequest {
    "id": number
    "avatar"?: string
    "isProhibition"?: boolean
    "title"?: string
    "abstract"?: string
    isSearch?: boolean // 是否被搜索
    verification?: number // 群验证 0 不允许任何人添加  1 允许任何人添加  2 需要验证消息 3 需要回答问题  4  需要正确回答问题
    isInvite?: boolean // 是否可邀请好友
    isTemporarySession?: boolean // 是否开启临时会话
    verificationQuestion?: verificationQuestionType
}

// 群更新
export function groupUpdateInfoApi(data: groupUpdateRequest): Promise<baseResponse<string>> {
    return useAxios.put("/api/group/group", data)
}

export interface groupMemberType {
    "userId": number
    "userNickname": string
    "avatar": string
    "isOnline": boolean
    "isFriend": boolean
    "role": number
    "memberNickname": string
    "createdAt": string
    "newMsgDate": string
    prohibitionTime?: number // 当前群成员的禁言时间
}

export interface groupMemberListRequest extends paramsType {
    id: number
    sort: "role asc" | "new_msg_date desc"
}

// 群成员列表
export function groupMemberListApi(params: groupMemberListRequest): Promise<baseResponse<listResponse<groupMemberType>>> {
    return useAxios.get("/api/group/member", {params})
}

// 解散群
export function groupRemoveApi(id: number): Promise<baseResponse<string>> {
    return useAxios.delete("/api/group/group/" + id.toString())
}

export interface groupMemberExitRequest {
    memberId: number // 用户id
    id: number // 群id
}

// 踢出群聊
export function groupMemberExitApi(params: groupMemberExitRequest): Promise<baseResponse<string>> {
    return useAxios.delete("/api/group/member", {params: params})
}

export interface groupMemberRoleUpdateRequest {
    "id": number // 群id
    "memberId": number // 用户id
    "role": number // 角色
}

// 变更用户角色
export function groupMemberRoleUpdateApi(data: groupMemberRoleUpdateRequest): Promise<baseResponse<string>> {
    return useAxios.put("/api/group/member/role", data)
}

export interface groupMemberNicknameRequest {
    "id": number
    "memberId": number
    "nickname": string
}

// 修改群成员昵称
export function groupMemberNicknameApi(data: groupMemberNicknameRequest): Promise<baseResponse<string>> {
    return useAxios.put("/api/group/member/nickname", data)
}

export interface groupSessionType {
    "groupId": number
    "title": string
    "avatar": string
    "newMsgDate": string
    "newMsgPreview": string
    "isTop": boolean
}

// 群会话列表
export function groupSessionListApi(): Promise<baseResponse<listResponse<groupSessionType>>> {
    return useAxios.get("/api/group/session", {params: {limit: -1}})
}


export interface groupHistoryType {
    groupID: number
    "userID": number
    "userNickname": string
    "userAvatar": string
    "msg": msgType
    "id": number
    "msgType": number
    "createdAt": string
    "isMe": boolean
    "memberNickname": string
    msgPreview: string
    showDate: boolean
}

export interface groupHistoryListParams extends paramsType {
    id: number
}

export function groupHistoryListApi(params: groupHistoryListParams): Promise<baseResponse<listResponse<groupHistoryType>>> {
    return useAxios.get("/api/group/history/" + params.id.toString(), {params})
}


// 群会话置顶
export function groupTopApi(groupId: number, isTop: boolean): Promise<baseResponse<string>> {
    return useAxios.put("/api/group/top", {
        groupId,
        isTop
    })
}

export interface groupMemberProhibitionRequest {
    "groupId": number
    "memberId": number
    "prohibitionTime"?: number
}

export function groupMemberProhibitionApi(data: groupMemberProhibitionRequest): Promise<baseResponse<string>> {
    return useAxios.put("/api/group/member/prohibition", data)
}

export function groupChatDeleteApi(groupID: number, msgIdList: number[]): Promise<baseResponse<string>> {
    return useAxios.delete("/api/group/history/" + groupID.toString(), {data: {msgIdList}})
}