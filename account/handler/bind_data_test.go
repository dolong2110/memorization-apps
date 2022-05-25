package handler

import (
	"bytes"
	"encoding/json"
	"github.com/dolong2110/Memoirization-Apps/account/model/apperrors"
	"github.com/dolong2110/Memoirization-Apps/account/model/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockBindData struct {
	mock.Mock
}

func (m *MockBindData) bindData(ctx *gin.Context, req interface{}) bool {
	ret := m.Called(ctx, req)

	var r0 bool
	if ret.Get(0) != nil {
		r0 = ret.Get(1).(bool)
	}

	return r0
}

func TestBindData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(mocks.MockUserService)
	mockTokenService := new(mocks.MockTokenService)

	router := gin.Default()

	NewHandler(&Config{
		Engine:       router,
		UserService:  mockUserService,
		TokenService: mockTokenService,
	})

	t.Run("Not application/json Content-type - 1", func(t *testing.T) {
		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"email":    "dummy@gmail.com",
			"password": "abcdefgh",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "multipart/form-data")
		router.ServeHTTP(rr, request)
		body, _ := ioutil.ReadAll(rr.Body)

		var resp apperrors.Response
		_ = json.Unmarshal(body, &resp)

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.Error.Code)
		assert.Equal(t, "/signin only accepts Content-Type application/json", resp.Error.Message)
	})

	t.Run("Not application/json Content-type - 2", func(t *testing.T) {

		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request = &http.Request{
			Header: make(http.Header),
		}

		mocks.MockJsonPost(ctx, map[string]interface{}{"foo": "bar"})
		res := bindData(ctx, map[string]interface{}{"foo": "bar"})

		assert.Equal(t, false, res)
	})
}
