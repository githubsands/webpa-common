package bookkeeping

import (
	"github.com/Comcast/webpa-common/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestEmptyBookkeeper(t *testing.T) {
	var (
		assert           = assert.New(t)
		require          = require.New(t)
		transactorCalled = false

		bookkeeper = New()
		logger     = logging.NewCaptureLogger()
	)
	require.NotNil(bookkeeper)
	req := httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(logging.WithLogger(req.Context(), logger))

	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		transactorCalled = true
		writer.Write([]byte("payload"))
		writer.WriteHeader(200)
	})
	rr := httptest.NewRecorder()

	bookkeeper(handler).ServeHTTP(rr, req)
	assert.True(transactorCalled)

	select {
	case result := <-logger.Output():
		assert.Len(result, 2)
	default:
		assert.Fail("CaptureLogger must capture something")

	}
}

func TestBookkeeper(t *testing.T) {
	var (
		assert           = assert.New(t)
		require          = require.New(t)
		transactorCalled = false


		bookkeeper = New(WithRequests(Path), WithResponses(Code))
		logger     = logging.NewCaptureLogger()
	)


	require.NotNil(bookkeeper)
	req := httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(logging.WithLogger(req.Context(), logger))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		transactorCalled = true
		writer.Write([]byte("payload"))
		writer.WriteHeader(200)
	})


	bookkeeper(handler).ServeHTTP(rr, req)

	assert.True(transactorCalled)

	select {
	case result := <-logger.Output():
		assert.Len(result, 4)
		assert.Equal(req.URL.Path, result["path"])
		assert.Equal(200, result["code"])
	default:
		assert.Fail("CaptureLogger must capture something")

	}
}
