package chat_models

import "fim_server/common/models"

// UserChatDeleteModel 用户删除聊天记录表
type UserChatDeleteModel struct {
	models.Model
	UserID uint `json:"userID"`
	ChatID uint `json:"chatID"` // 聊天记录的id
}
