EchoWBT
=======

.. image:: https://github.com/josuebrunel/echowbt/workflows/Test/badge.svg?branch=master
    :target: https://github.com/josuebrunel/echowbt/workflows/Test/badge.svg?branch=master

.. image:: https://coveralls.io/repos/github/josuebrunel/echowbt/badge.svg?branch=master
    :target: https://coveralls.io/github/josuebrunel/echowbt?branch=master


**EchoWBT** is a simple wrapper of *httptest* which helps test *Echo* apps easily

Installation
------------

.. code:: go
    
    go get github.com/josuebrunel/echowbt


Quickstart
----------

.. code:: go

    package echowb_test

    import (
        "github.com/josuebrunel/echowbt"
        "github.com/labstack/echo/v4"
        log "github.com/sirupsen/logrus"
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
                log.Info("This is a POST")
                status = http.StatusCreated
            case "PUT":
                log.Info("This is a PUT")
                status = http.StatusNoContent
            case "PATCH":
                log.Info("This is a PATCH")
                status = http.StatusNoContent
            case "DELETE":
                log.Info("This is DELETE")
                status = http.StatusAccepted
            default:
                log.Info("This is GET")
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
        Client echowb.Client
    }

    func (e *EchoWBTestSuite) SetupSuite() {
        e.Client = echowb.New()
    }

    func TestEchoWBT(t *testing.T) {
        suite.Run(t, new(EchoWBTestSuite))
    }

    func (e *EchoWBTestSuite) TestGet() {
        url := echowb.URL{Path: "/"}
        rec := e.Client.Get(url, GenericHandler(), nil, echowb.Headers{})
        assert.Equal(e.T(), http.StatusOK, rec.Code)
        url = echowb.URL{Path: "/?lastname=kouka&firstname=kim"}
        rec = e.Client.Get(url, GenericHandler(), nil, echowb.Headers{})
        assert.Equal(e.T(), http.StatusOK, rec.Code)
        params := []string{"id"}
        values := []string{"1"}
        url = echowb.URL{Path: "/:id", Params: params, Values: values}
        rec = e.Client.Get(url, GenericHandler(), nil, echowb.Headers{})
        assert.Equal(e.T(), http.StatusOK, rec.Code)
        data := echowb.JSONDecode(rec.Body)
        assert.Equal(e.T(), "Loking", data["lastname"])
    }

    func (e *EchoWBTestSuite) TestPost() {
        url := echowb.URL{Path: "/"}
        u := User{Firstname: "Josué", Lastname: "Kouka", Age: 30}
        rec := e.Client.Post(url, GenericHandler(), echowb.JSONEncode(u), echowb.Headers{})
        assert.Equal(e.T(), http.StatusCreated, rec.Code)
    }

    func (e *EchoWBTestSuite) TestPut() {
        params := []string{"id"}
        values := []string{"1"}
        url := echowb.URL{Path: "/:id", Params: params, Values: values}
        u := User{Firstname: "Josué", Lastname: "Kouka", Age: 30}
        headers := echowb.Headers{"Authorization": "Bearer <mytoken>"}
        rec := e.Client.Put(url, GenericHandler(), echowb.JSONEncode(u), headers)
        assert.Equal(e.T(), http.StatusNoContent, rec.Code)
    }

    func (e *EchoWBTestSuite) TestPatch() {
        params := []string{"id"}
        values := []string{"1"}
        url := echowb.URL{Path: "/:id", Params: params, Values: values}
        u := User{Firstname: "Josué", Lastname: "Kouka", Age: 30}
        headers := echowb.Headers{"Authorization": "Bearer <mytoken>"}
        rec := e.Client.Patch(url, GenericHandler(), echowb.JSONEncode(u), headers)
        assert.Equal(e.T(), http.StatusNoContent, rec.Code)
    }

    func (e *EchoWBTestSuite) TestDelete() {
        params := []string{"id"}
        values := []string{"1"}
        url := echowb.URL{Path: "/:id", Params: params, Values: values}
        headers := echowb.Headers{"Authorization": "Bearer <mytoken>"}
        rec := e.Client.Delete(url, GenericHandler(), nil, headers)
        assert.Equal(e.T(), http.StatusAccepted, rec.Code)
    }

Voila
