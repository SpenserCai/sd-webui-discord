// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UserInfoHandlerFunc turns a function with the right signature into a user info handler
type UserInfoHandlerFunc func(UserInfoParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn UserInfoHandlerFunc) Handle(params UserInfoParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// UserInfoHandler interface for that can handle valid user info params
type UserInfoHandler interface {
	Handle(UserInfoParams, interface{}) middleware.Responder
}

// NewUserInfo creates a new http.Handler for the user info operation
func NewUserInfo(ctx *middleware.Context, handler UserInfoHandler) *UserInfo {
	return &UserInfo{Context: ctx, Handler: handler}
}

/*
	UserInfo swagger:route GET /user_info user userInfo

Get User Info
*/
type UserInfo struct {
	Context *middleware.Context
	Handler UserInfoHandler
}

func (o *UserInfo) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUserInfoParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		*r = *aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc.(interface{}) // this is really a interface{}, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
