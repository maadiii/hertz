package server

import (
	"context"
	"io"
	"mime/multipart"
	"net"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
)

type Context struct {
	context.Context //nolint
	request         *app.RequestContext
}

// The host is valid until returning from RequestHandler.
func (ctx *Context) Host() []byte {
	return ctx.request.Host()
}

// RemoteAddr returns client address for the given request.
//
// If address is nil, it will return zeroTCPAddr.
func (ctx *Context) RemoteAddr() net.Addr {
	return ctx.request.RemoteAddr()
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (ctx *Context) Set(key string, value any) {
	ctx.request.Set(key, value)
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
//
// In case the Key is reset after response, Value() return nil if ctx.Key is nil.
func (ctx *Context) Value(key any) any {
	return ctx.request.Value(key)
}

func (ctx *Context) SetIdentity(identity Identity) {
	ctx.request.Set(identityKey, identity)
}

func (ctx *Context) Identity() Identity {
	return ctx.request.Value(identityKey).(Identity)
}

// ClientIP tries to parse the headers in [X-Real-Ip, X-Forwarded-For].
// It calls RemoteIP() under the hood. If it cannot satisfy the requirements,
// use engine.SetClientIPFunc to inject your own implementation.
func (ctx *Context) ClientIP() string {
	return ctx.request.ClientIP()
}

// ContentType returns the Content-Type header of the request.
func (ctx *Context) ContentType() []byte {
	return ctx.request.ContentType()
}

// Cookie returns the value of the request cookie key.
func (ctx *Context) Cookie(key string) []byte {
	return ctx.request.Cookie(key)
}

// Loop fn for every k/v in Keys
func (ctx *Context) ForEachKey(f func(key string, v any)) {
	ctx.request.ForEachKey(f)
}

// FullPath returns a matched route full path. For not found routes
// returns an empty string.
//
//	router.GET("/user/:id", func(c context.Context, ctx *app.RequestContext) {
//	    ctx.FullPath() == "/user/:id" // true
//	})
func (ctx *Context) FullPath() string {
	return ctx.request.FullPath()
}

// Header is an intelligent shortcut for ctx.Response.Header.Set(key, value).
// It writes a header in the response.
// If value == "", this method removes the header `ctx.Response.Header.Del(key)`.
func (ctx *Context) SetHeader(key, value string) {
	ctx.request.Header(key, value)
}

// SaveUploadedFile uploads the form file to specific dst.
func (ctx *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	return ctx.request.SaveUploadedFile(file, dst)
}

// Path returns requested path.
//
// The path is valid until returning from RequestHandler.
func (ctx *Context) Path() []byte {
	return ctx.request.Path()
}

// IfModifiedSince returns true if lastModified exceeds 'If-Modified-Since'
// value from the request header.
//
// The function returns true also 'If-Modified-Since' request header is missing.
func (ctx *Context) IfModifiedSince(lastModified time.Time) bool {
	return ctx.request.IfModifiedSince(lastModified)
}

// SetCookie adds a Set-Cookie header to the Response's headers.
//
//	Parameter introduce:
//	name and value is used to set cookie's name and value, eg. Set-Cookie: name=value
//	maxAge is use to set cookie's expiry date, eg. Set-Cookie: name=value; max-age=1
//	path and domain is used to set the scope of a cookie, eg. Set-Cookie: name=value;domain=localhost; path=/;
//	secure and httpOnly is used to sent cookies securely; eg. Set-Cookie: name=value;HttpOnly; secure;
//	sameSite let servers specify whether/when cookies are sent with cross-site requests; eg. Set-Cookie: name=value;HttpOnly; secure; SameSite=Lax;
//
//	For example:
//	1. ctx.SetCookie("user", "hertz", 1, "/", "localhost",protocol.CookieSameSiteLaxMode, true, true)
//	add response header --->  Set-Cookie: user=hertz; max-age=1; domain=localhost; path=/; HttpOnly; secure; SameSite=Lax;
//	2. ctx.SetCookie("user", "hertz", 10, "/", "localhost",protocol.CookieSameSiteLaxMode, false, false)
//	add response header --->  Set-Cookie: user=hertz; max-age=10; domain=localhost; path=/; SameSite=Lax;
//	3. ctx.SetCookie("", "hertz", 10, "/", "localhost",protocol.CookieSameSiteLaxMode, false, false)
//	add response header --->  Set-Cookie: hertz; max-age=10; domain=localhost; path=/; SameSite=Lax;
//	4. ctx.SetCookie("user", "", 10, "/", "localhost",protocol.CookieSameSiteLaxMode, false, false)
//	add response header --->  Set-Cookie: user=; max-age=10; domain=localhost; path=/; SameSite=Lax;
func (ctx *Context) SetCookie(name, value string, maxAge int, path, domain string, sameSite protocol.CookieSameSite, secure, httpOnly bool) {
	ctx.request.SetCookie(name, value, maxAge, path, domain, sameSite, secure, httpOnly)
}

func (ctx *Context) IsEnableTrace() bool {
	return ctx.request.IsEnableTrace()
}

// SetEnableTrace sets whether enable trace.
//
// NOTE: biz handler must not modify this value, otherwise, it may panic.
func (ctx *Context) SetEnableTrace(enable bool) {
	ctx.request.SetEnableTrace(enable)
}

func (ctx *Context) RequestBodyStream() io.Reader {
	return ctx.request.RequestBodyStream()
}

// URI returns requested uri.
//
// The uri is valid until returning from RequestHandler.
func (ctx *Context) URI() *protocol.URI {
	return ctx.request.URI()
}

// UserAgent returns the value of the request user_agent.
func (ctx *Context) UserAgent() []byte {
	return ctx.request.UserAgent()
}

type Identity map[string]any

const identityKey = "identity"
