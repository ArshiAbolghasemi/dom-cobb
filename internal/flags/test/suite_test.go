package flags_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFlags(t *testing.T) {
	gin.SetMode(gin.TestMode)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flags Suite")
}
