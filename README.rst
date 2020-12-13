EchoWBT
=======

.. image:: https://github.com/josuebrunel/echowbt/workflows/Test/badge.svg?branch=master
    :target: https://github.com/josuebrunel/echowbt/actions?query=workflow%3Atest

.. image:: https://coveralls.io/repos/github/josuebrunel/echowbt/badge.svg?branch=master
    :target: https://coveralls.io/github/josuebrunel/echowbt?branch=master

.. image:: https://pkg.go.dev/badge/github.com/josuebrunel/echowbt.svg
    :target: https://pkg.go.dev/github.com/josuebrunel/echowbt

.. image:: https://goreportcard.com/badge/github.com/josuebrunel/echowbt
    :target: https://goreportcard.com/report/github.com/josuebrunel/echowbt

.. image:: https://img.shields.io/badge/License-MIT-blue.svg
    :target: https://github.com/josuebrunel/echowbt/blob/master/LICENSE


**EchoWBT** is a simple wrapper of *httptest* which helps test *Echo* apps easily

Installation
------------

.. code:: go

    go get github.com/josuebrunel/echowbt


Example
-------

.. code:: go


    package echowb_test

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
