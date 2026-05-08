import {type baseResponse, useAxios} from "@/api/index";

export interface uploadImageResponse {
    url: string
}

// 图片上传
export function uploadImageApi(file: File, imageType: "avatar" | "group_avatar" | "chat"): Promise<baseResponse<uploadImageResponse>> {
    const form = new FormData()
    form.set("image", file)
    form.set("imageType", imageType)
    return useAxios.post("/api/file/image", form, {
        headers: {
            "Content-Type": "multipart/form-data"
        }
    })
}

export interface uploadFileResponse {
    src: string
}

export function uploadFileApi(file: File): Promise<baseResponse<uploadFileResponse>> {
    const form = new FormData()
    form.set("file", file)
    return useAxios.post("/api/file/file", form, {
        headers: {
            "Content-Type": "multipart/form-data"
        }
    })
}