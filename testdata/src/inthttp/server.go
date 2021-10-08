package inthttp

import (
	"net/http"
)

type dummy struct{}

func (r *dummy) Handle(string, http.Handler)                 {}
func (r *dummy) NotHandle(string, http.Handler)              {}
func (r dummy) ServeHTTP(http.ResponseWriter, *http.Request) {}

type dummy2 struct{}

func (r *dummy2) Handle(string) {}

type dummy3 struct{}

func (r *dummy3) Handle(string, http.Handler, string) {}

func Test() {

	x := dummy{}
	x.Handle(USED, x)
	x.NotHandle(UNUSED, x) // NotHandle

	y := dummy2{}
	y.Handle(USED) // Handle, but only one parameter

	z := dummy3{}
	z.Handle(UNUSED, nil, USED) // Handle, but too many parameters
}
