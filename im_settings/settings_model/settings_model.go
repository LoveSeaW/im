package settings_model

import (
	"im_server/common/models"
	"im_server/common/models/ctype"
)

type SettingsModel struct {
	models.Model
	Site ctype.SiteType `json:"site"`
	QQ   ctype.QQType   `json:"qq"`
}
