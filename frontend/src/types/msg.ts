import type {userBaseType} from "@/api/chat_api";

export interface textMsg {
    "content": string
}

export interface imageMsg {
    "title": string
    "src": string
}

export interface videoMsg {
    "title": string
    "src": string
}

export interface fileMsg {
    size?: number
    src: string
    title?: string
    type?: string
}

export interface tipMsg {
    content: string
    status: "error"
}

export interface imageTextMsg {
    content: string
}

export interface voiceMsg {
    src: string
    time: number
}

export interface replyMsg {
    "msgID": number
    "content": string
    "msg"?: msgType
    "userID"?: number
    "userNickName"?: string
    "originMsgDate"?: string
    replyMsgPreview?: string
}

export interface videoCallMsg {
    flag?: number
    type?: string
    data?: any
    endTime?: string
    startTime?: string
}

export interface withdrawMsg {
    msgID: number
    content?: string
}

export interface msgType {
    "type": number
    "textMsg"?: textMsg // 1
    "imageMsg"?: imageMsg // 2
    "videoMsg"?: videoMsg // 3
    "fileMsg"?: fileMsg // 4
    "voiceMsg"?: voiceMsg // 5
    "videoCallMsg"?: videoCallMsg // 7
    withdrawMsg?: withdrawMsg // 8
    "replyMsg"?: replyMsg // 9
    "tipMsg"?: tipMsg // 12
    "imageTextMsg"?: imageTextMsg // 14
}

export interface msgChatType {
    "id": number
    "user": userBaseType
    "isMe": boolean
    "created_at": string
    "msg": msgType
    showDate: boolean
    msgPreview: string
}