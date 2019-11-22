package main

import (
	"fmt"
	"net/http"

	"go.zoe.im/goesty"

	"github.com/gorilla/mux"
)

type errT struct {
	msg string
}

func (e errT) Error() string {
	return e.msg
}

type spec struct {
	Version string
	Name string
	Annotation map[string]string
}

type person struct {
	Name    string `json:"name" api:"in=query"`
	Namespace     string `json:"namespace" api:"in=path"`
	Spec spec `json:"spec" api:"in=body,to=-"` 
	Message string `json:"message"`
	Cookie string `json:"cookie" api:"to=header(set-cookie)"`
}

func echoParamSimple(name string, age int, token string) (*person, error) {
	fmt.Println(name, age, token)
	if name == "" {
		return nil, errT{"name cannot be empty"}
	}

	if age == 0 {
		return nil, fmt.Errorf("aget cannot be 0")
	}

	fmt.Printf("my name is %s and i'm ", name)
	if age >= 18 {
		fmt.Println("adult")
	} else {
		fmt.Println("teenager")
	}

	if token == "" {
		return nil, http.ErrNoCookie
	}

	fmt.Println("my token:", token)

	return &person{
		Name: name,
		Message: fmt.Sprintf("Hello %s!", name),
		Cookie: token,
	}, nil
}

func echoParamStruct(p *person) (*person, error) {
	if p.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	p.Message = fmt.Sprintf("Hello %s!", p.Name)
	return p, nil
}

func main() {
	r := mux.NewRouter()

	sty := goesty.NewRuntime(
		goesty.OptionVarsFunc(mux.Vars),
	)

	// Request(r *http.Request, args []interface{}) {}
	// Response(w http.ResponseWriter, res []interface{}) {}

	r.Handle("/echo/simple-0", sty.MustNew(echoParamSimple))
	r.Handle("/echo/simple-1/{id:[0-9]+}", sty.MustNew(echoParamSimple).
		Param(0).InQuery("name").
		Param(1).InPath("id").
		Param(2).InHeader("Authorization"),
	)
	// mux.Handle("/echo/simple-2", o.MustNew(echoParamSimple, o.ParamInQuery("name"), o.Param(0).InQuery("name")))
	// mux.Handle("/echo/struct-0", o.MustNew(echoParamStruct))
	// mux.Handle("/echo/struct-1", o.MustNew(echoParamStruct).Param(0).Key("name").InBody(".data"))
	// mux.Handle("/echo/struct-2", o.MustNew(echoParamStruct, o.Param(0, o.ParamInBody(".data"))))
	// mux.Handle("/echo/struct-code", o.MustNew(echoParamStruct, o.HttpCode(200)))

	http.ListenAndServe(":8080", r)
}