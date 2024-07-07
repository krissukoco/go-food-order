package wrapper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// Types
var (
	ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()
	errType = reflect.TypeOf((*error)(nil)).Elem()
)

// Errors
var (
	ErrFunctionType = errors.New("function type is invalid")
)

func MustWrap(f interface{}) fiber.Handler {
	h, err := Wrap(f)
	if err != nil {
		panic(err)
	}
	return h
}

// Wrap wraps handler function in the form of:
//
//	func(ctx context.Context, request interface{}) (response interface{}, err error)
//
// to fiber.Handler
func Wrap(f interface{}) (fiber.Handler, error) {
	v := reflect.ValueOf(f)
	t := v.Type()

	// Check kind
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("%w: f type is %s", ErrFunctionType, t.Kind().String())
	}

	// Check arguments and parameters
	if t.NumIn() != 2 {
		return nil, fmt.Errorf("%w: number of arguments is not 2", ErrFunctionType)
	}
	if t.NumOut() != 2 {
		return nil, fmt.Errorf("%w: number of returns is not 2", ErrFunctionType)
	}
	argCtx := t.In(0)
	if err := checkInCtx(argCtx); err != nil {
		return nil, err
	}
	argRequest := t.In(1)
	if err := checkInRequest(argRequest); err != nil {
		return nil, err
	}
	returnResponse := t.Out(0)
	if err := checkOutResponse(returnResponse); err != nil {
		return nil, err
	}
	returnError := t.Out(1)
	if err := checkOutError(returnError); err != nil {
		return nil, err
	}

	req := reflect.New(argRequest).Interface()

	return func(c *fiber.Ctx) error {
		// Populate req from query, params, and body
		if err := c.QueryParser(req); err != nil {
			log.Println("queryParser error: ", err)
			return err
		}
		if err := c.ParamsParser(req); err != nil {
			log.Println("paramsParser error: ", err)
			return err
		}
		if c.Method() != "GET" && c.Method() != "DELETE" {
			if err := c.BodyParser(req); err != nil {
				log.Println("bodyParser error: ", err)
				return err
			}
		}

		ctx := c.Context()
		ctxValue := reflect.ValueOf(ctx)
		reqValue := reflect.ValueOf(req).Elem()

		results := v.Call([]reflect.Value{ctxValue, reqValue})
		errIntf := results[1].Interface()
		if errIntf != nil {
			return errIntf.(error)
		}
		return c.JSON(results[0].Interface())
	}, nil
}

func checkInCtx(t reflect.Type) error {
	if !t.Implements(ctxType) {
		return fmt.Errorf("%w: first argument should implement context.Context", ErrFunctionType)
	}
	return nil
}

func checkInRequest(t reflect.Type) error {
	return nil
}

func checkOutResponse(t reflect.Type) error {
	return nil
}

func checkOutError(t reflect.Type) error {
	if !t.Implements(errType) {
		return fmt.Errorf("%w: must return error type", ErrFunctionType)
	}
	return nil
}
