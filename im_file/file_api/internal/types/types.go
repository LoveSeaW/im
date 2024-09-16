// Code generated by goctl. DO NOT EDIT.
package types

type ImageRequest struct {
	UserID uint `header:"User-ID"`
}

type ImageResponse struct {
	Url string `json:"url"`
}

type FileRequest struct {
	UserID uint `header:"User-ID"`
}

type FileResponse struct {
	Src string `json:"src"`
}

type ImageShowRequest struct {
	ImageName string `path:"imageName"`
}

type ImageShowResponse struct {
}

type FileListRequest struct {
	Page  int    `form:"page,optional"`
	Limit int    `form:"limit,optional"`
	Key   string `form:"key,optional"`
}

type FileListInfoResponse struct {
	FileName  string `json:"fileName"` // 文件名称
	Size      int64  `json:"size"`     // 文件大小
	Path      string `json:"path"`     // 文件的实际路径
	CreatedAt string `json:"createdAt"`
	ID        uint   `json:"id"`
	WebPath   string `json:"webPath"` // 访问路径
}

type FileListResponse struct {
	List  []FileListInfoResponse `json:"list"`
	Count int64                  `json:"count"`
}

type FileListRemoveRequest struct {
	IdList []uint `json:"idList"`
}

type FileListRemoveResponse struct {
}
