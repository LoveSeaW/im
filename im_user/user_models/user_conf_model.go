package user_models

import (
	"im_server/common/models"
	"im_server/common/models/ctype"
)

// UserConfModel 用户配置表
type UserConfModel struct {
	models.Model
	UserID               uint                        `json:"userID"`
	UserModel            UserModel                   `gorm:"foreignKey:UserID" json:"-"`
	RecallMessage        *string                     `gorm:"size:32" json:"recallMessage"` // 撤回消息的提示内容
	FriendOnline         bool                        `json:"friendOnline"`                 // 好友上线提醒
	Sound                bool                        `json:"sound"`                        // 声音
	SecureLink           bool                        `json:"secureLink"`                   // 安全链接
	SavePwd              bool                        `json:"savePwd"`                      // 保存密码
	SearchUser           int8                        `json:"searchUser"`                   // 别人查找到你的方式 0 不允许别人查找到我， 1  通过用户号找到我 2 可以通过昵称搜索到我
	Verification         int8                        `json:"verification"`                 // 好友验证 0 不允许任何人添加  1 允许任何人添加  2 需要验证消息 3 需要回答问题  4  需要正确回答问题
	VerificationQuestion *ctype.VerificationQuestion `json:"verificationQuestion"`         // 验证问题  为3和4的时候需要
	Online               bool                        `json:"online"`                       // 是否在线
	CurtailChat          bool                        `json:"curtailChat"`                  // 限制聊天
	CurtailAddUser       bool                        `json:"curtailAddUser"`               // 限制加人
	CurtailCreateGroup   bool                        `json:"curtailCreateGroup"`           // 限制建群
	CurtailInGroupChat   bool                        `json:"curtailInGroupChat"`           // 限制加群
}

// ProblemCount 问题的个数
func (uc UserConfModel) ProblemCount() (c int) {
	if uc.VerificationQuestion != nil {
		if uc.VerificationQuestion.Problem1 != nil {
			c += 1
		}
		if uc.VerificationQuestion.Problem2 != nil {
			c += 1
		}
		if uc.VerificationQuestion.Problem3 != nil {
			c += 1
		}
	}
	return
}
