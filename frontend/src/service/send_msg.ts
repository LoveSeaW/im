import {useStore} from "@/stores";
import type {imageMsg, msgType, replyMsg, videoCallMsg, videoMsg, voiceMsg} from "@/types/msg";

export class Msg {
    type: "user" | "group"
    ws: WebSocket

    constructor(type: "user" | "group") {
        this.type = type
        const store = useStore()
        if (type === "user") {
            this.ws = store.chatWs
        } else {
            this.ws = store.groupWs
        }
    }

    send(id: number, msg: msgType) {
        const data: { revUserID?: number, groupID?: number, msg?: msgType } = {}
        if (this.type === "user") {
            data["revUserID"] = id
        } else {
            data["groupID"] = id
        }
        data["msg"] = msg

        this.ws.send(JSON.stringify(data))
    }


    // 文本消息
    textMsg(id: number, content: string) {
        this.send(id, {
            type: 1,
            textMsg: {
                content: content
            }
        })
    }

    // 图片消息
    imageMsg(id: number, imageMsg: imageMsg) {
        this.send(id, {
            type: 2,
            imageMsg
        })
    }

    // 视频消息
    videoMsg(id: number, videoMsg: videoMsg) {
        this.send(id, {
            type: 3,
            videoMsg
        })
    }

    // 文件消息
    fileMsg(id: number, src: string) {
        this.send(id, {
            type: 4,
            fileMsg: {
                src,
            }
        })
    }


    // 语音消息
    voiceMsg(id: number, voiceMsg: voiceMsg) {
        this.send(id, {
            type: 5,
            voiceMsg
        })
    }

    // 视频通话
    videoCallMsg(id: number, videoCallMsg: videoCallMsg) {
        this.send(id, {
            type: 7,
            videoCallMsg
        })
    }

    // 撤回消息
    withdrawMsg(id: number, msgID: number) {
        this.send(id, {
            type: 8,
            withdrawMsg: {msgID}
        })
    }

    // 回复消息
    replyMsg(id: number, replyMsg: replyMsg) {
        this.send(id, {
            type: 9,
            replyMsg
        })
    }

    // 图文消息
    imageTextMsg(id: number, content: string) {
        this.send(id, {
            type: 14,
            imageTextMsg: {
                content,
            }
        })
    }

}

new Msg("user")