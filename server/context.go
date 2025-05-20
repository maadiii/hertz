package server

import (
	"context"
	"io"
	"mime/multipart"
	"net"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/protocol"
)

type Request struct {
	rc *app.RequestContext
}

// The host is valid until returning from RequestHandler.
func (req *Request) Host() []byte {
	return req.rc.Host()
}

// RemoteAddr returns client address for the given request.
//
// If address is nil, it will return zeroTCPAddr.
func (req *Request) RemoteAddr() net.Addr {
	return req.rc.RemoteAddr()
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (req *Request) Set(key string, value any) {
	req.rc.Set(key, value)
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
//
// In case the Key is reset after response, Value() return nil if req.Key is nil.
func (req *Request) Value(key any) any {
	return req.rc.Value(key)
}

// ClientIP tries to parse the headers in [X-Real-Ip, X-Forwarded-For].
// It calls RemoteIP() under the hood. If it cannot satisfy the requirements,
// use engine.SetClientIPFunc to inject your own implementation.
func (req *Request) ClientIP() string {
	return req.rc.ClientIP()
}

// ContentType returns the Content-Type header of the request.
func (req *Request) ContentType() []byte {
	return req.rc.ContentType()
}

// Cookie returns the value of the request cookie key.
func (req *Request) Cookie(key string) []byte {
	return req.rc.Cookie(key)
}

// Loop fn for every k/v in Keys
func (req *Request) ForEachKey(f func(key string, v any)) {
	req.rc.ForEachKey(f)
}

// FullPath returns a matched route full path. For not found routes
// returns an empty string.
//
//	router.GET("/user/:id", func(c context.Context, req *app.RequestContext) {
//	    req.FullPath() == "/user/:id" // true
//	})
func (req *Request) FullPath() string {
	return req.rc.FullPath()
}

// Header is an intelligent shortcut for req.Response.Header.Set(key, value).
// It writes a header in the response.
// If value == "", this method removes the header `req.Response.Header.Del(key)`.
func (req *Request) SetHeader(key, value string) {
	req.rc.Header(key, value)
}

// GetHeader returns value from request headers.
func (req *Request) GetHeader(key string) []byte {
	return req.rc.GetHeader(key)
}

// SaveUploadedFile uploads the form file to specific dst.
func (req *Request) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	return req.rc.SaveUploadedFile(file, dst)
}

// Path returns requested path.
//
// The path is valid until returning from RequestHandler.
func (req *Request) Path() []byte {
	return req.rc.Path()
}

// IfModifiedSince returns true if lastModified exceeds 'If-Modified-Since'
// value from the request header.
//
// The function returns true also 'If-Modified-Since' request header is missing.
func (req *Request) IfModifiedSince(lastModified time.Time) bool {
	return req.rc.IfModifiedSince(lastModified)
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
//	1. req.SetCookie("user", "hertz", 1, "/", "localhost",protocol.CookieSameSiteLaxMode, true, true)
//	add response header --->  Set-Cookie: user=hertz; max-age=1; domain=localhost; path=/; HttpOnly; secure; SameSite=Lax;
//	2. req.SetCookie("user", "hertz", 10, "/", "localhost",protocol.CookieSameSiteLaxMode, false, false)
//	add response header --->  Set-Cookie: user=hertz; max-age=10; domain=localhost; path=/; SameSite=Lax;
//	3. req.SetCookie("", "hertz", 10, "/", "localhost",protocol.CookieSameSiteLaxMode, false, false)
//	add response header --->  Set-Cookie: hertz; max-age=10; domain=localhost; path=/; SameSite=Lax;
//	4. req.SetCookie("user", "", 10, "/", "localhost",protocol.CookieSameSiteLaxMode, false, false)
//	add response header --->  Set-Cookie: user=; max-age=10; domain=localhost; path=/; SameSite=Lax;
func (req *Request) SetCookie(name, value string, maxAge int, path, domain string, sameSite protocol.CookieSameSite, secure, httpOnly bool) {
	req.rc.SetCookie(name, value, maxAge, path, domain, sameSite, secure, httpOnly)
}

func (req *Request) IsEnableTrace() bool {
	return req.rc.IsEnableTrace()
}

// SetEnableTrace sets whether enable trace.
//
// NOTE: biz handler must not modify this value, otherwise, it may panic.
func (req *Request) SetEnableTrace(enable bool) {
	req.rc.SetEnableTrace(enable)
}

func (req *Request) RequestBodyStream() io.Reader {
	return req.rc.RequestBodyStream()
}

// URI returns requested uri.
//
// The uri is valid until returning from RequestHandler.
func (req *Request) URI() *protocol.URI {
	return req.rc.URI()
}

// UserAgent returns the value of the request user_agent.
func (req *Request) UserAgent() []byte {
	return req.rc.UserAgent()
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
func (req *Request) Next(c context.Context) {
	req.rc.Next(c)
}

// Don't use Abort method in handlers. It just use for middlewares
func (req *Request) Abort() {
	req.rc.Abort()
}

// Don't use AbortWithStatus method in handlers. It just use for middlewares
func (req *Request) AbortWithStatus(code int) {
	req.rc.AbortWithStatus(code)
}

// Don't use AbortWithMsg method in handlers. It just use for middlewares
func (req *Request) AbortWithMsg(msg string, statusCode int) {
	req.rc.AbortWithMsg(msg, statusCode)
}

// Don't use AbortWithStatusJSON method in handlers. It just use for middlewares
func (req *Request) AbortWithStatusJSON(code int, jsonObj any) {
	req.rc.AbortWithStatusJSON(code, jsonObj)
}

// Don't use AbortWithError method in handlers. It just use for middlewares
func (req *Request) AbortWithError(code int, err error) *errors.Error {
	return req.rc.AbortWithError(code, err)
}

// Method return request method.
//
// Returned value is valid until returning from RequestHandler.
func (req *Request) Method() string {
	return string(req.rc.Method())
}

// Error attaches an error to the current context. The error is pushed to a list of errors.
//
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database together,
// print a log, or append it in the HTTP response.
// Error will panic if err is nil.
func (req *Request) Error(err error) *errors.Error {
	return req.rc.Error(err)
}
