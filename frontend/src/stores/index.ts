import {defineStore} from 'pinia'
import {parseToken} from "@/utils/parseToken";
import {type userConfType, userInfoApi} from "@/api/user_api";
import {ElMessage} from "element-plus";
import {settingsInfoApi, type settingsType} from "@/api/settings_api";
import {logoutApi} from "@/api/auth_api";
import router from "@/router";
import type {groupHistoryType, groupInfoType, groupUpdateRequest} from "@/api/group_api";
import {groupDetailApi, groupUpdateInfoApi} from "@/api/group_api";
import type {chatHistoryType} from "@/api/chat_api";
import type {msgChatType} from "@/types/msg";

interface userInfoType {
    exp: number
    nickname: string
    role: number
    userID: number
    token: string
    avatar: string
}

const userInfo: userInfoType = {
    exp: 0,
    nickname: "",
    role: 0,
    userID: 0,
    token: "",
    avatar: "",
}

const settingsInfo: settingsType = {}

const userConfInfo: userConfType = {
    userID: 0,
    nickname: "",
    abstract: "",
    avatar: "",
    friendOnline: false,
    sound: false,
    secureLink: false,
    savePwd: false,
    searchUser: 0,
    verification: 0,
}


interface searchType {
    name: "search_user" | "search_group" | ""
    value: string
}

const searchData: searchType = {
    name: "",
    value: ""
}

const groupData: groupInfoType = {
    groupId: 0,
    title: "",
    abstract: "",
    memberCount: 0,
    memberOnlineCount: 0,
    avatar: "",
    creator: {
        userId: 0,
        avatart: "",
        nickname: "",
    },
    adminList: [],
    role: 0,
    isProhibition: false,
    isSearch: false,
    isInvite: false,
    isTemporarySession: false,
}

const chatMsgData: chatHistoryType = {
    id: 0,
    sendUser: {
        id: 0,
        nickName: "",
        avatar: "",
    },
    revUser: {
        id: 0,
        nickName: "",
        avatar: "",
    },
    isMe: false,
    created_at: "",
    msg: {
        type: 0
    },
    systemMsg: null,
    msgPreview: "",
    showDate: false
}
const groupMsgData: groupHistoryType = {
    groupID: 0,
    "userID": 0,
    "userNickname": "",
    "userAvatar": "",
    "msg": {
        "type": 0,
        textMsg: undefined,
    },
    "id": 0,
    "msgType": 0,
    "createdAt": "",
    "isMe": false,
    "memberNickname": "",
    msgPreview: "",
    showDate: false
}
let chatWs: WebSocket
let groupWs: WebSocket

interface useStoreType {
    userInfo: userInfoType,
    settingsInfo: settingsType,
    userConfInfo: userConfType,
    theme: boolean
    searchData: searchType
    groupData: groupInfoType
    chatMsgData: chatHistoryType
    groupMsgData: groupHistoryType
    chatWs: WebSocket
    groupWs: WebSocket
    chatContent: string
    replayChatMsgData?: msgChatType
    showMsgDate: boolean
}

export const useStore = defineStore('counter', {
    state: (): useStoreType => {
        return {
            userInfo: userInfo,
            settingsInfo: settingsInfo,
            userConfInfo: userConfInfo,
            theme: true, // true 是白天  false 是黑夜
            searchData: searchData, // 搜索页面用到的数据
            groupData: groupData,
            chatMsgData: chatMsgData,
            groupMsgData: groupMsgData,
            chatWs: chatWs,
            groupWs: groupWs,
            chatContent: "",
            showMsgDate: false
        }
    },
    actions: {
        // 设置token
        async setToken(token: string) {
            const payload = parseToken(token)
            this.userInfo.token = token
            this.userInfo.exp = payload.exp
            this.userInfo.nickname = payload.nickname
            this.userInfo.role = payload.role
            this.userInfo.userID = payload.userID
            // 去拿用户的用户信息
            await this.getUserConf()
            this.userInfo.avatar = this.userConfInfo.avatar
            // 调一下持久化
            this.saveToken()

            // 去连ws
            this.initChatWs()
            this.initGroupWs()
        },
        showMsgDateMethod(){
            this.showMsgDate = true
            setTimeout(() => {
                this.showMsgDate = false
            }, 3000)
        },
        // 保存token
        saveToken() {
            localStorage.setItem("userInfo", JSON.stringify(this.userInfo))
        },
        // 加载token
        loadToken() {
            const val = localStorage.getItem("userInfo")
            if (!val) {
                // 没有登陆，或者登陆失效
                return
            }
            try {
                this.userInfo = JSON.parse(val)
            } catch (e) {
                localStorage.removeItem("userInfo")
                return;
            }
        },
        // 获取用户信息
        async getUserConf() {

            let res = await userInfoApi()
            if (res.code) {
                ElMessage.error(res.msg)
                return
            }
            this.userConfInfo = res.data
        },
        // 获取系统信息
        async getSettingsInfo() {
            let res = await settingsInfoApi()
            if (res.code) {
                ElMessage.error(res.msg)
                return
            }
            this.settingsInfo = res.data
        },
        async logout() {
            // 先去调后端的注销接口
            // 然后清空本地的 用户store，然后清空storage
            logoutApi()
            this.userInfo = userInfo
            this.userConfInfo = userConfInfo
            localStorage.removeItem("userInfo")
            // 跳转到登陆页
            router.push({
                name: "login"
            })
            ElMessage.success("注销成功")

        },

        // 设置主题
        setTheme(theme?: boolean) {
            if (theme !== undefined) {
                this.theme = theme
            } else {
                this.theme = !this.theme
            }

            // 根据theme的不同，去设置class
            if (this.theme) {
                // 要设置为白天
                document.documentElement.classList.remove("dark")
            } else {
                // 要设置为黑夜
                document.documentElement.classList.add("dark")
            }

            // 持久化
            localStorage.setItem("theme", this.themeString)
        },

        // 加载主题
        loadTheme() {
            const val = localStorage.getItem("theme")
            if (!val) {
                return
            }
            if (val === "dark") {
                this.setTheme(false)
                return;
            }
        },


        // 获取群信息
        async getGroupData(id: number) {
            // 如果这个群之前有请求了，那就不去请求了，用之前的数据
            if (id === this.groupData.groupId) {
                return
            }
            let res = await groupDetailApi(id)
            if (res.code) {
                ElMessage.error(res.msg)
                return
            }
            Object.assign(this.groupData, res.data)
        },

        setChatContent(content: string) {
            this.chatContent = ""
            setTimeout(() => {
                this.chatContent = content
            })
        },

        // 更新群
        async updateGroup() {
            if (this.groupData.groupId === 0) {
                return
            }
            const data: groupUpdateRequest = {
                id: this.groupData.groupId,
                avatar: this.groupData.avatar,
                isProhibition: this.groupData.isProhibition,
                title: this.groupData.title,
                abstract: this.groupData.abstract,
            }
            let res = await groupUpdateInfoApi(data)
            if (res.code) {
                ElMessage.error(res.msg)
                return
            }
            ElMessage.success("群信息修改成功")
        },

        initChatWs() {
            if (!this.isLogin) {
                return
            }
            let proto = "ws"
            if (location.protocol === "https:") {
                proto = "wss"
            }
            const ws = new WebSocket(`${proto}://${location.host}/api/chat/ws/chat?token=${this.userInfo.token}`)
            ws.onopen = (ev: Event) => {
                console.log("ws连接成功")
            }
            ws.onmessage = (ev: MessageEvent) => {
                try {
                    const data: chatHistoryType = JSON.parse(ev.data)
                    this.chatMsgData = data
                } catch (e) {
                    console.log(e)
                    return
                }
            }
            ws.onerror = (ev: Event) => {
                console.log(ev)
            }
            this.chatWs = ws
        },
        initGroupWs() {
            if (!this.isLogin) {
                return
            }
            let proto = "ws"
            if (location.protocol === "https:") {
                proto = "wss"
            }
            const ws = new WebSocket(`${proto}://${location.host}/api/group/ws/chat?token=${this.userInfo.token}`)
            ws.onopen = (ev: Event) => {
                console.log("group ws连接成功")
            }
            ws.onmessage = (ev: MessageEvent) => {
                try {
                    const data: groupHistoryType = JSON.parse(ev.data)
                    this.groupMsgData = data
                } catch (e) {
                    console.log(e)
                    return
                }
            }
            ws.onerror = (ev: Event) => {
                console.log(ev)
            }
            this.groupWs = ws
        },
    },
    getters: {
        // 是否登陆
        isLogin(): boolean {
            // exp的时间戳-现在的时间戳  为正就是没有过期
            return this.userInfo.token != ""
        },
        themeString(): string {
            return this.theme ? "" : "dark"
        }
    }
})
