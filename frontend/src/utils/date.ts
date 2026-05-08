import dayjs from "dayjs";
import rTime from "dayjs/plugin/relativeTime"
import 'dayjs/locale/zh-cn' // 导入本地化语言

dayjs.extend(rTime)
dayjs.locale('zh-cn') // 使用本地化语言
export function dateTimeFormat(time: string): string {
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

export function dateFormat(time: string): string {
    return dayjs(time).format('YYYY-MM-DD')
}

export function timeFormat(time: string): string{
    return dayjs(time).format('HH:mm:ss')
}

export function relativeCurrentTime(time: string): string {
    return dayjs(time).fromNow()
}


export function time_different(date1: string, date2: string) {
    const date3 = new Date(date2).getTime() - new Date(date1).getTime();   //时间差的毫秒数
    //------------------------------
    //计算出小时数
    const leave1 = date3 % (24 * 3600 * 1000)    //计算天数后剩余的毫秒数
    const hours = Math.floor(leave1 / (3600 * 1000))
    //计算相差分钟数
    const leave2 = leave1 % (3600 * 1000)        //计算小时数后剩余的毫秒数
    const minutes = Math.floor(leave2 / (60 * 1000))
    //计算相差秒数
    const leave3 = leave2 % (60 * 1000)      //计算分钟数后剩余的毫秒数
    const seconds = Math.round(leave3 / 1000)
    return `${toZero(hours)}:${toZero(minutes)}:${toZero(seconds)}`
}

// 不足2位补0
function toZero(n: number){
    return n <= 9 ? "0"+n : n
}