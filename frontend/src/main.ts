import { createApp } from 'vue'
import { createPinia } from 'pinia'
import "@/assets/base.css"
import "@/assets/iconfont.css"
import "@/assets/iconfont_color.css"
import "@/assets/my_font.css"
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import "@/assets/theme.css"


const app = createApp(App)

app.use(ElementPlus)
app.use(createPinia())
app.use(router)

app.mount('#app')
