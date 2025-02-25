package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AdminMiddlewareSuite struct {
	suite.Suite
	testingServer *httptest.Server
}

func (suite *AdminMiddlewareSuite) SetupSuite() {
	router := gin.Default()
	router.Use(AdminMiddleware())
	router.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// create and run the testing server
	testingServer := httptest.NewServer(router)
	suite.testingServer = testingServer
}

func (suite *AdminMiddlewareSuite) TearDownSuite() {
	suite.testingServer.Close()
}

func (suite *AdminMiddlewareSuite) TestNoClaims() {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(AdminMiddleware())
	router.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.JSONEq(suite.T(), `{"error":"User not allowed"}`, w.Body.String())
}

func (suite *AdminMiddlewareSuite) TestNoUserRole() {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(AdminMiddleware())
	router.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	// Set claims in the context
	router.Use(func(c *gin.Context) {
		c.Set("claims", map[string]interface{}{})
		c.Next()
	})
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.JSONEq(suite.T(), `{"error":"User not allowed"}`, w.Body.String())
}

func (suite *AdminMiddlewareSuite) TestUserRoleNotAdmin() {
	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	router.Use(AdminMiddleware())
	router.GET("/admin", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	// Set claims and userRole in the context
	router.Use(func(c *gin.Context) {
		c.Set("claims", map[string]interface{}{})
		c.Set("userRole", "USER")
		c.Next()
	})
	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.JSONEq(suite.T(), `{"error":"User not allowed"}`, w.Body.String())
}

// Run the test suite
func TestAdminMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AdminMiddlewareSuite))
}
