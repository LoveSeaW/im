export function copyText(content: string) {
    if (!content) {
        return
    }
    let cInput = document.createElement("input")
    document.body.appendChild(cInput)
    cInput.value = content
    cInput.select() // 选取文本框内容

    // 执行浏览器复制命令
    // 复制命令会将当前选中的内容复制到剪切板中（这里就是创建的input标签）
    // Input要在正常的编辑状态下原生复制方法才会生效
    document.execCommand("copy")
    // 复制成功后再将构造的标签 移除
    document.body.removeChild(cInput)
}