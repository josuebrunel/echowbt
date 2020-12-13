package echowbt_test

import (
	"github.com/josuebrunel/echowbt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type User struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

// Handlers

func GenericHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		status, u := http.StatusOK, User{}
		switch c.Request().Method {
		case "POST":
			log.Info("POST /")
			status = http.StatusCreated
		case "PUT":
			log.Info("PUT /:id")
			if err := c.Bind(&u); err != nil {
				log.Error(err)
				return err
			}
			log.Info("USER: ", u)
			status = http.StatusNoContent
		case "PATCH":
			log.Info("PATCH /:id")
			if err := c.Bind(&u); err != nil {
				log.Error(err)
				return err
			}
			log.Info("USER: ", u)
			status = http.StatusNoContent
		case "DELETE":
			log.Info("DELETE /:id")
			status = http.StatusAccepted
		default:
			log.Info("GET /")
			log.Info("QueryParams", c.QueryParams())
			log.Info("Params", c.ParamNames())
			log.Info("Values", c.ParamValues())
			if id := c.Param("id"); id != "" {
				u = User{1, "Yosuke", "Loking", 30}
			}
		}
		return c.JSON(status, u)
	}
}

type Application struct {
	E *echo.Echo
}

func App() (app Application) {
	app = Application{E: echo.New()}
	app.E.Logger.SetLevel(log.DEBUG)
	app.E.Use(middleware.Logger())
	app.E.GET("/", GenericHandler())
	app.E.GET("/:id", GenericHandler())
	app.E.POST("/", GenericHandler())
	app.E.PUT("/:id", GenericHandler())
	app.E.PATCH("/:id", GenericHandler())
	app.E.DELETE("/:id", GenericHandler())
	return
}

type EchoWBTestSuite struct {
	suite.Suite
	Client echowbt.Client
}

func (e *EchoWBTestSuite) SetupSuite() {
	e.Client = echowbt.New()
}

func TestEchoWBT(t *testing.T) {
	suite.Run(t, new(EchoWBTestSuite))
}

func (e *EchoWBTestSuite) TestGet() {
	url := echowbt.URL{Path: "/"}
	rec := e.Client.Get(url, GenericHandler(), nil, echowbt.Headers{})
	assert.Equal(e.T(), http.StatusOK, rec.Code)
	url = echowbt.URL{Path: "/?lastname=kouka&firstname=kim"}
	rec = e.Client.Get(url, GenericHandler(), nil, echowbt.Headers{})
	assert.Equal(e.T(), http.StatusOK, rec.Code)
	params := []string{"id"}
	values := []string{"1"}
	url = echowbt.URL{Path: "/:id", Params: params, Values: values}
	rec = e.Client.Get(url, GenericHandler(), nil, echowbt.Headers{})
	assert.Equal(e.T(), http.StatusOK, rec.Code)
	data := echowbt.JSONDecode(rec.Body)
	assert.Equal(e.T(), "Loking", data["lastname"])
}

func (e *EchoWBTestSuite) TestPost() {
	url := echowbt.URL{Path: "/"}
	u := User{Firstname: "Josué", Lastname: "Kouka", Age: 30}
	rec := e.Client.Post(url, GenericHandler(), echowbt.JSONEncode(u), echowbt.Headers{})
	assert.Equal(e.T(), http.StatusCreated, rec.Code)
}

func (e *EchoWBTestSuite) TestPut() {
	params := []string{"id"}
	values := []string{"1"}
	url := echowbt.URL{Path: "/:id", Params: params, Values: values}
	u := User{Firstname: "Josué", Lastname: "Kouka", Age: 30}
	headers := echowbt.Headers{"Authorization": "Bearer <mytoken>"}
	rec := e.Client.Put(url, GenericHandler(), echowbt.JSONEncode(u), headers)
	assert.Equal(e.T(), http.StatusNoContent, rec.Code)
}

func (e *EchoWBTestSuite) TestPatch() {
	params := []string{"id"}
	values := []string{"1"}
	url := echowbt.URL{Path: "/:id", Params: params, Values: values}
	u := User{Firstname: "Josué", Lastname: "Kouka", Age: 30}
	headers := echowbt.Headers{"Authorization": "Bearer <mytoken>"}
	rec := e.Client.Patch(url, GenericHandler(), echowbt.JSONEncode(u), headers)
	assert.Equal(e.T(), http.StatusNoContent, rec.Code)
}

func (e *EchoWBTestSuite) TestDelete() {
	params := []string{"id"}
	values := []string{"1"}
	url := echowbt.URL{Path: "/:id", Params: params, Values: values}
	headers := echowbt.Headers{"Authorization": "Bearer <mytoken>"}
	rec := e.Client.Delete(url, GenericHandler(), nil, headers)
	assert.Equal(e.T(), http.StatusAccepted, rec.Code)
}
