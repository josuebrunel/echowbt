EchoWBT
=======

.. image:: https://github.com/josuebrunel/echowbt/workflows/test/badge.svg?branch=master
    :target: https://github.com/josuebrunel/echowbt/actions?query=workflow%3Atest

.. image:: https://coveralls.io/repos/github/josuebrunel/echowbt/badge.svg?branch=master
    :target: https://coveralls.io/github/josuebrunel/echowbt?branch=master

.. image:: https://pkg.go.dev/badge/github.com/josuebrunel/echowbt.svg
    :target: https://pkg.go.dev/github.com/josuebrunel/echowbt

.. image:: https://goreportcard.com/badge/github.com/josuebrunel/echowbt
    :target: https://goreportcard.com/report/github.com/josuebrunel/echowbt

.. image:: https://img.shields.io/badge/License-MIT-blue.svg
    :target: https://github.com/josuebrunel/echowbt/blob/master/LICENSE


**EchoWBT** is a simple wrapper of *httptest* allowing you to simply test your Echo_ app
With **EchoWBT** handles for you :

* the instanciation of *httptest.NewRequest* and *httptest.NewRecorder*
* the binding of the two above with to an *echo.Context*
* the setting of *request headers*
* the setting of *Path*, *ParamNames* and *ParamsValues* for your context

.. _Echo: https://github.com/labstack/echo

Installation
------------

.. code:: go

    go get github.com/josuebrunel/echowbt

Quickstart
----------

.. code:: go

    import (
        "github.com/josuebrunel/echowbt"
        "github.com/project/app"
        "testing"
        "github.com/stretchr/testify/assert"
        "net/http"
    )

    func TestPingHandler(t *testing.T)
        client := echowbt.New()
        rec := client.Get(echowbt.URL{"/ping"}, app.PingHandler(), nil, echowbt.Headers{"Authorization": "X-Auth xyw:uiyu"})
        assert.Equal(t, http.StatusOK, rec.Code)
        data := echowbt.JSONDecode(rec.Body)
        assert.Equal(t, int64(1), data["count"])
        assert.Equal(t, "ping", data["data"])


Set a default content type
^^^^^^^^^^^^^^^^^^^^^^^^^^

The default *Content-Type* is *application/json*. To change it use *SetHeaders* method

.. code:: go

    client.SetHeaders(echow.Headers{"Content-Type": "text/html"})

URL Construction
^^^^^^^^^^^^^^^^

.. code:: go

    // simple url
    url := echowbt.URL{"/"}

URL Named Params
""""""""""""""""

.. code:: go

    params := echowbt.URLParams{"family_id", "member_id"}
    values := echowbt.URLParams{"1", "3"}
    url = echowbt.URL{Path: "/:family_id/:member_id", Params: params, Values: values}
    rec := client.Get(url, MyHanlder(), nil, echow.Headers{})

Headers
^^^^^^^

You can pass *headers* to your request by using *echowbt.Headers* type

.. code:: go

    headers := echowbt.Headers{"Content-Type": "application/x-www-form-urlencoded", "Authorization": "Token <mytoken>"}
    rec := client.Post(url, MyHanlder(), []byte{"username=josh&password=joshpwd"}, headers)


Send JSON Data
^^^^^^^^^^^^^^

You can send a *JSON Payload* by using *echowbt.JSONEncode* func

.. code:: go

    u := User{Username: "lokinghd"}
    rec := client.Post(url, MyHanlder(), echowbt.JSONEncode(u), headers)


Send MultipartForm Data
^^^^^^^^^^^^^^^^^^^^^^^

You can send a *MultipartForm Data* by using *echowbt.FormData* func

.. code:: go

    formFields := echowbt.Fields{"firstname": "Josué", "lastname": "Kouka", "City": "Pointe-Noire"}
    fileFields := echowbt.Fields{"avatar": "/tmp/jk.png"}
    formData, _ := echowbt.FormData(formFields, fileFields)
    headers := echowbt.Headers{"Content-Type": formData.ContentType} // IMPORTANT FOR PART BOUNDARY
    rec := client.Post(url, MyHanlder(), FormData.Data, headers)

Decoding JSON Response
^^^^^^^^^^^^^^^^^^^^^^

You can decode your JSON Response by using *echowbt.JSONDecode* func

.. Code:: go

    rec := client.Get(url, MyHanlder(), JSONEncode(payload), headers)
    data = echowbt.JSONDecode(rec.Body)
    assert.Equal(t, int64(1), data["count"])
    assert.Equal(t, "uuid", data["data"]["uuid"])


For in depth examples check the **main_test.go** file

Voila ;) !
