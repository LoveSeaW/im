package handler

import (
	"im_server/common/response"
	"im_server/im_user/user_api/internal/logic/Admin"
	"im_server/im_user/user_api/internal/svc"
	"im_server/im_user/user_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserCurtailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserCurtailRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := Admin.NewUserCurtailLogic(r.Context(), svcCtx)
		resp, err := l.UserCurtail(&req)
		response.Response(r, w, resp, err)

	}
}
