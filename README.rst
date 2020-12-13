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

    package pkg_test

    import (
        "github.com/josuebrunel/echowbt"
        "testing"
    )
