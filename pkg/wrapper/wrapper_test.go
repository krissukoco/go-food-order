package wrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type TestReq struct {
	Field  string `json:"field"`
	Query1 string `query:"query1"`
	Param1 string `param:"param1"`
}

type TestResp struct {
	Message string `json:"message"`
}

func Test_Wrap(t *testing.T) {
	testHandler := func(ctx context.Context, req TestReq) (*TestResp, error) {
		return &TestResp{Message: fmt.Sprintf("Success! field: %s | query1: %s | param1: %s", req.Field, req.Query1, req.Param1)}, nil
	}

	app := fiber.New()
	app.Post("/:param1", MustWrap(testHandler))

	go func() {
		if err := app.Listen(":31000"); err != nil {
			panic(err)
		}
	}()

	jsonBody := `{"field":"this is field"}`
	b := bytes.NewBuffer([]byte(jsonBody))
	req, err := http.NewRequest("POST", "http://localhost:31000/test-param-1?query1=testquery", b)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Logf("resp status code is %d: %s", resp.StatusCode, resp.Status)
		var errstr string
		if err := json.NewDecoder(resp.Body).Decode(&errstr); err != nil {
			panic(err)
		}
		t.Logf("erro: %s", errstr)
		t.Fatal("error")
	}

	var respStruct TestResp
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		panic(err)
	}

	t.Logf("respStruct: %v", respStruct)
}
