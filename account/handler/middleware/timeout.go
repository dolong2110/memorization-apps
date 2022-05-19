package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dolong2110/Memoirization-Apps/account/model/apperrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

// Reference: https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go

// implements http.Writer, but tracks if Writer has timed out
// or has already written its header to prevent
// header and body overwrites
// also locks access to this writer to prevent race conditions
// holds the gin.ResponseWriter which we'll manually call Write()
// on in the middleware function to send response
type timeoutWriter struct {
	gin.ResponseWriter
	header       http.Header
	writerBuffer bytes.Buffer // The zero value for Buffer is an empty buffer ready to use.
	mutex        sync.Mutex
	timedOut     bool
	wroteHeader  bool
	code         int
}

// Writes the response, but first makes sure there
// hasn't already been a timeout
// In http.ResponseWriter interface
func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mutex.Lock()
	defer tw.mutex.Unlock()
	if tw.timedOut {
		return 0, nil
	}

	return tw.writerBuffer.Write(b)
}

//WriteHeader In http.ResponseWriter interface
func (tw *timeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	tw.mutex.Lock()
	defer tw.mutex.Unlock()
	// We do not write the header if we've timed out or written the header
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}

// set that the header has been written
func (tw *timeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

// Header "relays" the header, h, set in struct
// In http.ResponseWriter interface
func (tw *timeoutWriter) Header() http.Header {
	return tw.header
}

//SetTimeOut sets timedOut field to true
func (tw *timeoutWriter) SetTimedOut() {
	tw.timedOut = true
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}

// Timeout middleware to set time out stop for api
func Timeout(timeout time.Duration, errTimeout *apperrors.Error) gin.HandlerFunc {
	return func(c *gin.Context) {
		// set Gin's writer as our custom writer
		tw := &timeoutWriter{ResponseWriter: c.Writer, header: make(http.Header)}
		c.Writer = tw

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// update gin request context
		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})        // to indicate handler finished
		panicChan := make(chan interface{}, 1) // used to handle panics if we can't recover

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next()               // calls subsequent middleware(s) and handler
			finished <- struct{}{} // send signal if handler successfully complete
		}()

		select {
		case <-panicChan:
			// if we cannot recover from panic,
			// send internal server error
			err := apperrors.NewInternal()
			tw.ResponseWriter.WriteHeader(err.Status())
			eResp, _ := json.Marshal(gin.H{
				"error": err,
			})
			tw.ResponseWriter.Write(eResp)
		case <-finished:
			// if finished, set headers and write resp
			tw.mutex.Lock()
			defer tw.mutex.Unlock()
			// map Headers from tw.Header() (written to by gin)
			// to tw.ResponseWriter for response
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			// tw.writerBuffer will have been written to already when gin writes to tw.Write()
			tw.ResponseWriter.Write(tw.writerBuffer.Bytes())
		case <-ctx.Done():
			// ctx.Done() is automatically called when timeout has occurred, send errTimeout and write headers
			tw.mutex.Lock()
			defer tw.mutex.Unlock()
			// ResponseWriter from gin
			tw.ResponseWriter.Header().Set("Content-Type", "application/json")
			tw.ResponseWriter.WriteHeader(errTimeout.Status())
			eResp, _ := json.Marshal(gin.H{
				"error": errTimeout,
			})
			tw.ResponseWriter.Write(eResp)
			c.Abort()
			tw.SetTimedOut()
		}
	}
}
