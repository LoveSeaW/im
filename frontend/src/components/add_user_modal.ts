import {defineComponent, h, createApp, ref} from "vue";
import Search_dialog from "@/components/search_dialog.vue";
import type {userSearchResponse} from "@/api/user_api";


const visible = ref<undefined|boolean>(undefined)

export function showAddUserModal(userInfo: userSearchResponse) {
    if (visible.value === undefined){
        const component = defineComponent({
            setup() {
                const onClose = () => {
                    visible.value = false
                }
                visible.value = true
                return {
                    visible,
                    userInfo,
                    onClose
                }
            },
            render() {
                return h(Search_dialog, {
                    visible: this.visible as boolean,
                    userInfo: this.userInfo,
                    mode: "user",
                    onClose: this.onClose
                })
            }
        })

        const app = createApp(component)
        const div = document.createElement("div")
        app.mount(div)
        document.body.appendChild(div)
    }
    visible.value = true
}