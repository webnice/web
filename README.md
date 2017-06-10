# web

[![GoDoc](https://godoc.org/gopkg.in/webnice/web.v1?status.png)](http://godoc.org/gopkg.in/webnice/web.v1)
[![Coverage Status](https://coveralls.io/repos/github/webnice/web/badge.svg?branch=v1)](https://coveralls.io/github/webnice/web?branch=v1)
[![Build Status](https://travis-ci.org/webnice/web.svg?branch=v1)](https://travis-ci.org/webnice/web)
[![CircleCI](https://circleci.com/gh/webnice/web/tree/v1.svg?style=svg)](https://circleci.com/gh/webnice/web/tree/v1)

#### Description
This is a micro framework for create web servers, based on net/http with manage and route functions.

Framework do not interfere with the style of your code and do not require to somehow modify Handlers and HandlerFunc functions.
But at the same time framework adds to your application a full-functions routing, middlewares and lots of useful functionality or just simplifies your code.


#### Dependencies

	golang.org/x/net
	golang.org/x/text
	golang.org/x/crypto

All dependencies in the vendor folder saved. This lib do not require install dependencies.


#### Install
```bash
go get gopkg.in/webnice/web.v1
```
