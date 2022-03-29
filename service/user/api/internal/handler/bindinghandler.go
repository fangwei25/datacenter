package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"web_game/service/user/api/internal/logic"
	"web_game/service/user/api/internal/svc"
	"web_game/service/user/api/internal/types"
)

func BindingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReqBinding
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewBindingLogic(r.Context(), svcCtx)
		resp, err := l.Binding(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
