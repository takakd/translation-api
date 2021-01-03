package util_test

import (
	"api/internal/pkg/util"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"time"
)

func ExampleGetAppDir() {
	// Returns the directory path to "dir" located the app binary.
	dir, err := util.GetAppDir()
	if err != nil {
		// Handle if err happened.
	}

	fmt.Println(dir != "")

	// Output:
	// true
}

func ExampleFileExists() {
	exists := os.Args[0]
	notExists := "invalid/path"

	fmt.Println(util.FileExists(exists))
	fmt.Println(util.FileExists(notExists))

	// Output:
	// true
	// false
}

func ExampleGetRequestBody() {
	const body = `{"name":"Tom", text":"hello."}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	// Get body from http.Request.
	b, err := util.GetRequestBody(r)
	if err != nil {
		// Handle error.
	}

	fmt.Println(string(b))

	// Output:
	// {"name":"Tom", text":"hello."}
}

func ExampleGetResponseBody() {
	const body = `{"name":"Tom", text":"hello."}`
	r := httptest.NewRecorder()
	r.WriteString(body)

	// Get body from http.Response.
	b, err := util.GetResponseBody(r.Result())
	if err != nil {
		// Handle error.
	}

	fmt.Println(string(b))

	// Output:
	// {"name":"Tom", text":"hello."}
}

func ExampleSetFormValueToStruct() {
	type Form struct {
		Param1 string `json:"param1"`
		Param2 string `json:"param2"`
	}

	v := url.Values{
		"param1": {"value1"},
		"param2": {"value2"},
	}
	s := Form{}

	err := util.SetFormValueToStruct(v, &s)
	if err != nil {
		// Handle if error happened
	}

	fmt.Println(s)

	// Output:
	// {value1 value2}
}

func ExampleParseUnixStr() {
	// 12/11/2020 @ 4:42pm (UTC)
	const ts = "1607704942"

	t, err := util.ParseUnixStr(ts)
	if err != nil {
		// Handle if error happened
	}

	fmt.Println(t.UTC().Format(time.RFC3339))

	// Output:
	// 2020-12-11T16:42:22Z
}

func ExampleIsStruct() {
	type A struct {
	}

	a := A{}
	i := 12
	s := "test"

	fmt.Println(util.IsStruct(a))
	fmt.Println(util.IsStruct(&a))
	fmt.Println(util.IsStruct(i))
	fmt.Println(util.IsStruct(s))

	// Output:
	// true
	// true
	// false
	// false
}
