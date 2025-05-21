package blitz_test

import (
	"context"
	"fmt"
	"net/http"
	"nyxze/fayth/blitz"
	"testing"
)

func TestNewRequest(t *testing.T) {

}

func TestMiddlware(t *testing.T) {
	ctx := context.Background()
	req, err := blitz.NewRequest(ctx, "GET", "fake", "hello", WithDefault()...)
	if err != nil {
		fmt.Println(err)
	}
	err = req.Send()
	if err != nil {
		fmt.Println(err)
	}
}
func WithDefault() []blitz.Options {
	return []blitz.Options{
		blitz.WithHeader("Hello", "Pouet"),
		blitz.WithCustomExecutor(blitz.ExecutorFunc(func(r *http.Request) (*http.Response, error) {
			fmt.Println("Fake exec request", r)
			return &http.Response{}, nil
		})),
		blitz.WithHeader("Yolo", "Rene"),
		blitz.WithMiddleware(FirstMid),
		blitz.WithMiddleware(SecondMid),
	}
}
func FirstMid(outer blitz.Executor) blitz.Executor {
	return blitz.ExecutorFunc(func(r *http.Request) (*http.Response, error) {
		fmt.Println("First")
		resp, err := outer.Do(r)
		return resp, err
	})
}

func SecondMid(outer blitz.Executor) blitz.Executor {
	return blitz.ExecutorFunc(func(r *http.Request) (*http.Response, error) {
		fmt.Println("Second")
		resp, err := outer.Do(r)
		return resp, err
	})
}
