import {defineStore} from 'pinia'
import type {eventKeyType} from "@/types/type";


interface useStoreType {
    eventKey: eventKeyType
}

export const eventStore = defineStore('eventStore', {
    state: (): useStoreType => {
        return {
            eventKey: "", // 需要事件通知的key
        }
    },
    actions: {
        setEventKey(key: eventKeyType) {
            this.eventKey = ""
            setTimeout(() => {
                this.eventKey = key
            })
        }
    },
})
