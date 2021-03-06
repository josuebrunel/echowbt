package echowbt

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strings"
	"testing"
)

type User struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname" form:"firstname"`
	Lastname  string `json:"lastname" form:"lastname"`
	Age       int    `json:"age" form:"age"`
}

// Handlers

func GenericHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		status, u := http.StatusOK, User{}
		log.Info("HEADERS: ", c.Request().Header)
		switch c.Request().Method {
		case "POST":
			log.Info("POST /")
			status = http.StatusCreated
			contentType := c.Request().Header.Get("Content-Type")
			if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
				c.Bind(&u)
				fname := c.FormValue("firstname")
				return c.String(status, fmt.Sprintf("Hello %s", fname))
			}
			if strings.HasPrefix(contentType, "multipart/form") {
				fname := c.FormValue("firstname")
				form, err := c.MultipartForm()
				if err != nil {
					log.Error(err)
					return err
				}
				filename := form.File["bio"][0].Filename
				return c.String(status, fmt.Sprintf("Hello %s ! Your file %s is up.", fname, filename))
			}
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
	Client Client
}

func (e *EchoWBTestSuite) SetupSuite() {
	e.Client = New()
	e.Client.SetHeaders(DictString{"Authorization": "Token <mytoken>"})
}

func TestEchoWBT(t *testing.T) {
	suite.Run(t, new(EchoWBTestSuite))
}

func (e *EchoWBTestSuite) TestGet() {
	url := NewURL("/", nil, nil)
	rec := e.Client.Get(url, GenericHandler(), nil, DictString{})
	assert.Equal(e.T(), http.StatusOK, rec.Code)
	url = NewURL("/", nil, DictString{"lastname": "kouka", "firstname": "kim"})
	rec = e.Client.Get(url, GenericHandler(), nil, DictString{})
	assert.Equal(e.T(), http.StatusOK, rec.Code)
	url = NewURL("/:id", DictString{"id": "1"}, nil)
	rec = e.Client.Get(url, GenericHandler(), nil, DictString{})
	assert.Equal(e.T(), http.StatusOK, rec.Code)
	data := JSONDecode(rec.Body)
	assert.Equal(e.T(), "Loking", data["lastname"])
}

func (e *EchoWBTestSuite) TestPost() {
	url := URL{Path: "/"}
	u := User{Firstname: "Josué", Lastname: "Kouka", Age: 30}
	rec := e.Client.Post(url, GenericHandler(), JSONEncode(u), DictString{})
	assert.Equal(e.T(), http.StatusCreated, rec.Code)
	// post form
	headers := DictString{"Content-Type": "application/x-www-form-urlencoded"}
	rec = e.Client.Post(url, GenericHandler(), []byte("firstname=Josué&lastname=Kouka"), headers)
	assert.Equal(e.T(), "Hello Josué", rec.Body.String())
}

func (e *EchoWBTestSuite) TestPostMultipartForm() {
	url := URL{Path: "/"}
	form, _ := FormData(DictString{"firstname": "Josué", "lastname": "Kouka"}, DictString{"bio": "testdata/bio.txt"})
	headers := DictString{"Content-Type": form.ContentType}
	rec := e.Client.Post(url, GenericHandler(), form.Data, headers)
	assert.Equal(e.T(), http.StatusCreated, rec.Code)
	expected := "Hello Josué ! Your file bio.txt is up."
	assert.Equal(e.T(), expected, rec.Body.String())
}

func (e *EchoWBTestSuite) TestPut() {
	url := NewURL("/:id", DictString{"id": "1"}, nil)
	u := User{Firstname: "Josué", Lastname: "Kouka", Age: 30}
	headers := DictString{"Authorization": "Bearer <mytoken>"}
	rec := e.Client.Put(url, GenericHandler(), JSONEncode(u), headers)
	assert.Equal(e.T(), http.StatusNoContent, rec.Code)
}

func (e *EchoWBTestSuite) TestPatch() {
	url := NewURL("/:id", DictString{"id": "1"}, nil)
	u := User{Firstname: "Josué", Lastname: "Kouka", Age: 30}
	headers := DictString{"Authorization": "Bearer <mytoken>"}
	rec := e.Client.Patch(url, GenericHandler(), JSONEncode(u), headers)
	assert.Equal(e.T(), http.StatusNoContent, rec.Code)
}

func (e *EchoWBTestSuite) TestDelete() {
	url := NewURL("/:id", DictString{"id": "1"}, nil)
	rec := e.Client.Delete(url, GenericHandler(), nil, DictString{})
	assert.Equal(e.T(), http.StatusAccepted, rec.Code)
}
