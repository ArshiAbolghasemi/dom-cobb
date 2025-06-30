package testutils

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func CreateJSONRequest(method, url string, reqBody any) (*gin.Context, *httptest.ResponseRecorder) {
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}
