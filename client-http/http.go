package clienthttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceContext struct{}

type ClientHttpContext struct {
	SvcCtx *ServiceContext
	Ctx    *gin.Context
}

type ClientResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ClientRespOptionFunc func(clientResp *ClientResp)

type FrontHttpFunc func(clientHttp *ClientHttpContext) ClientResp

func Make(svcCtx *ServiceContext, h FrontHttpFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientHttp := &ClientHttpContext{SvcCtx: svcCtx, Ctx: ctx}
		clientHttp.JSON(h(clientHttp))
	}
}

func (clientHttp *ClientHttpContext) WithCode(code int) ClientRespOptionFunc {
	return func(clientResp *ClientResp) {
		clientResp.Code = code
	}
}

func (clientHttp *ClientHttpContext) WithMessage(message string) ClientRespOptionFunc {
	return func(clientResp *ClientResp) {
		clientResp.Message = message
	}
}

func (clientHttp *ClientHttpContext) WithData(data any) ClientRespOptionFunc {
	return func(clientResp *ClientResp) {
		clientResp.Data = data
	}
}

func (clientHttp *ClientHttpContext) JSON(clientResp ClientResp) {
	clientHttp.Ctx.JSON(http.StatusOK, clientResp)
}

func (clientHttp *ClientHttpContext) ServerError(opts ...ClientRespOptionFunc) ClientResp {
	clientResp := &ClientResp{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong.",
	}
	return mapOptions(clientResp, opts...)
}

func (clientHttp *ClientHttpContext) BadReq(opts ...ClientRespOptionFunc) ClientResp {
	clientResp := &ClientResp{
		Code:    http.StatusBadRequest,
		Message: "Invalid request. Please verify your data.",
	}
	return mapOptions(clientResp, opts...)
}

func (clientHttp *ClientHttpContext) Unauthorized(opts ...ClientRespOptionFunc) ClientResp {
	clientResp := &ClientResp{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized. Invalid credentials.",
	}
	return mapOptions(clientResp, opts...)
}

func (clientHttp *ClientHttpContext) NotFound(opts ...ClientRespOptionFunc) ClientResp {
	clientResp := &ClientResp{
		Code:    http.StatusNotFound,
		Message: "Not found. The requested resource does not exist.",
	}
	return mapOptions(clientResp, opts...)
}

func (clientHttp *ClientHttpContext) Success(opts ...ClientRespOptionFunc) ClientResp {
	clientResp := &ClientResp{
		Code:    http.StatusOK,
		Message: "Operation completed successfully.",
	}
	return mapOptions(clientResp, opts...)
}

func mapOptions(resp *ClientResp, opts ...ClientRespOptionFunc) ClientResp {
	for _, opt := range opts {
		opt(resp)
	}

	return *resp
}
