package errs

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func GetDefaultRequestFields(r *http.Request) map[string]interface{} {
	var params map[string]string

	if r.Method != "GET" {
		r.ParseForm()
		// Access all parameters from the request
		allParameters := r.Form

		// Print all parameters
		for paramName, paramValues := range allParameters {
			for _, paramValue := range paramValues {
				params[paramName] = paramValue
			}
		}
	} else {

	}

	return map[string]interface{}{
		"user_id":    "some_user_id",
		"uri":        r.RequestURI,
		"method":     r.Method,
		"time":       time.Now(),
		"company_id": r.Header.Get("X-Company-ID"),
		"params":     params,
	}
}

// WrapErrorWithStackTrace wrap existing error with stack trace
// example: "bom not found. data_id: a53d25cc-c576-11ee-8ad4-0242ac17003 [stacktrace: /app/pkg/errs/err_data_not_found.go:29\n/app/internal/bom/repository/user_repository_impl.go:240\n/app/internal/bom/service/archive.go:36\n/app/internal/bom/handler/archive.go:14\n/app/internal/base/handler/base.go:62\n/app/pkg/middleware/auth_jwt.go:99\n/app/internal/base/handler/base.go:82\n]"
// since slog does not support newline so stacktrace will be separated by '\n'
func WrapErrorWithStackTrace(err error) error {
	return fmt.Errorf("%w [stacktrace: %s]", err, StackAndFile(2))
}

// GenerateErrorWithStackTrace generate new error with stack trace
// example: "bom not found. data_id: a53d25cc-c576-11ee-8ad4-0242ac17003 [stacktrace: /app/pkg/errs/err_data_not_found.go:29\n/app/internal/bom/repository/user_repository_impl.go:240\n/app/internal/bom/service/archive.go:36\n/app/internal/bom/handler/archive.go:14\n/app/internal/base/handler/base.go:62\n/app/pkg/middleware/auth_jwt.go:99\n/app/internal/base/handler/base.go:82\n]"
// since slog does not support newline so stacktrace will be separated by '\n'
func GenerateErrorWithStackTrace(msg string) error {
	return fmt.Errorf("%w [stacktrace: %s]", errors.New(msg), StackAndFile(2))
}

func UnwrapError(err error) error {
	unwrappedErr := errors.Unwrap(err)
	if unwrappedErr != nil {
		return unwrappedErr
	}
	return err
}

func StackAndFile(skip int) string {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	pathPrefix := "bitbucket.org/jurnal/scm"
	var lastFile string
	var firstFile string
	for i := skip; ; i++ { // Skip the expected number of frames

		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if firstFile == "" {
			firstFile = strings.ReplaceAll(fmt.Sprintf("%s:%d", file, line), pathPrefix, "...")
		}

		// ---- Skip un-necessary trace
		if strings.Contains(file, "net/http/server.go") {
			/**
			Often

			C:/Program Files/Go/src/net/http/server.go:2046 (0x3efdae)
				HandlerFunc.ServeHTTP: f(w, r)
			C:/Users/Will/Desktop/code/yii2/tnt/vendor/github.com/gorilla/mux/mux.go:210 (0x51b64e)
				(*Router).ServeHTTP: handler.ServeHTTP(w, req)
			C:/Users/Will/Desktop/code/yii2/tnt/vendor/gopkg.in/DataDog/dd-trace-go.v1/contrib/internal/httputil/trace.go:57 (0xad1109)
				TraceAndServe: httpinstr.WrapHandler(h, span).ServeHTTP(cfg.ResponseWriter, cfg.Request.WithContext(ctx))
			C:/Users/Will/Desktop/code/yii2/tnt/vendor/gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux/mux.go:119 (0xadb00b)
				(*Router).ServeHTTP: httputil.TraceAndServe(r.Router, &httputil.TraceConfig{
			C:/Program Files/Go/src/net/http/server.go:2878 (0x3f331a)
				serverHandler.ServeHTTP: handler.ServeHTTP(rw, req)
			C:/Program Files/Go/src/net/http/server.go:1929 (0x3eee87)
				(*conn).serve: serverHandler{c.server}.ServeHTTP(w, w.req)
			C:/Program Files/Go/src/runtime/asm_amd64.s:1581 (0x1a8f00)
				goexit: BYTE	$0x90	// NOP
			*/
			break
		}

		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d\n", file, line)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))

	}
	return strings.ReplaceAll(buf.String(), pathPrefix, "...")
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}
