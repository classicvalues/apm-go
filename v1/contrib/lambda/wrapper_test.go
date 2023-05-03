package lambda

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type TIn struct {
	Val int
}

type TOut struct {
	Val int
}

type dumbWrapper struct {
	beforeIsCalled bool
	afterIsCalled  bool
}

func (d *dumbWrapper) Before(ctx context.Context, msg json.RawMessage, args ...interface{}) context.Context {
	d.beforeIsCalled = true
	return ctx
}

func (d *dumbWrapper) After(res interface{}, err *typedError, args ...interface{}) interface{} {
	d.afterIsCalled = true
	return res
}

func TestWrapperInOut(t *testing.T) {
	fn := func(ctx context.Context, in *TIn) (*TOut, error) {
		return &TOut{Val: in.Val * 2}, nil
	}

	wr := &dumbWrapper{}
	fnW := HandlerWithWrapper(fn, wr)
	fnWrapped := fnW.(func(ctx context.Context, message json.RawMessage) (interface{}, error))
	inBytes, _ := json.Marshal(&TIn{Val: 23})
	tOut, err := fnWrapped(context.Background(), inBytes)
	assert.Nil(t, err)
	assert.Equal(t, 23*2, tOut.(*TOut).Val) // assert that fn is called
	assert.True(t, wr.beforeIsCalled)
	assert.True(t, wr.afterIsCalled)
}

func TestWrapperIn(t *testing.T) {
	fn := func(ctx context.Context, in *TIn) error {
		return errors.New("fn error")
	}

	wr := &dumbWrapper{}
	fnW := HandlerWithWrapper(fn, wr)
	fnWrapped := fnW.(func(ctx context.Context, message json.RawMessage) (interface{}, error))
	inBytes, _ := json.Marshal(&TIn{Val: 23})
	_, err := fnWrapped(context.Background(), inBytes)
	assert.NotNil(t, err)
	assert.True(t, wr.beforeIsCalled)
	assert.True(t, wr.afterIsCalled)
}

func TestWrapperInValid(t *testing.T) {
	fn := func(ctx context.Context, in *TIn, val int) (error, int) {
		return errors.New("fn error"), 0
	}

	wr := &dumbWrapper{}
	fnW := HandlerWithWrapper(fn, wr)
	fnWrapped := fnW.(func(ctx context.Context, in *TIn, val int) (error, int))
	err, _ := fnWrapped(context.Background(), &TIn{Val: 23}, 0)
	assert.NotNil(t, err)
	assert.False(t, wr.beforeIsCalled)
	assert.False(t, wr.afterIsCalled)
}

func TestInjectTraceContext(t *testing.T) {
	tw := traceWrapper{}
	val := tw.injectTraceContext(events.APIGatewayProxyResponse{
		StatusCode:        200,
		MultiValueHeaders: map[string][]string{"key": {"value1", "value2"}},
		Body:              "response body",
		IsBase64Encoded:   true,
	}, "my-x-trace-id")
	gwReq, ok := val.(events.APIGatewayProxyResponse)
	assert.True(t, ok)

	expected := events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{APMHTTPHeader: "my-x-trace-id"},
		MultiValueHeaders: map[string][]string{"key": {"value1", "value2"}},
		Body:              "response body",
		IsBase64Encoded:   true,
	}

	assert.Equal(t, expected, gwReq)
}

func TestInjectTraceContextWithExistingHeaders(t *testing.T) {
	tw := traceWrapper{}
	val := tw.injectTraceContext(events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{"header": "value"},
		MultiValueHeaders: map[string][]string{"key": {"value1", "value2"}},
		Body:              "response body",
		IsBase64Encoded:   true,
	}, "my-x-trace-id")
	gwReq, ok := val.(events.APIGatewayProxyResponse)
	assert.True(t, ok)

	expected := events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{APMHTTPHeader: "my-x-trace-id", "header": "value"},
		MultiValueHeaders: map[string][]string{"key": {"value1", "value2"}},
		Body:              "response body",
		IsBase64Encoded:   true,
	}

	assert.Equal(t, expected, gwReq)
}
