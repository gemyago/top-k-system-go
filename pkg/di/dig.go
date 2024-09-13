package di

import (
	"fmt"

	"go.uber.org/dig"
)

// Dig is used as DI toolkit https://pkg.go.dev/go.uber.org/dig
// we are not creating any abstraction over it, but we do have a set of tools to make it easier to use

type ConstructorWithOpts struct {
	Constructor interface{}
	Options     []dig.ProvideOption
}

func ProvideAll(container *dig.Container, providers ...interface{}) error {
	for i, provider := range providers {
		switch p := provider.(type) {
		case ConstructorWithOpts:
			if err := container.Provide(p.Constructor, p.Options...); err != nil {
				return fmt.Errorf("failed to provide %d-th dependency: %w", i, err)
			}
		default:
			if err := container.Provide(provider); err != nil {
				return fmt.Errorf("failed to provide %d-th dependency: %w", i, err)
			}
		}
	}
	return nil
}

// ProvideValue will create a constructor (e.g func) from a given value.
func ProvideValue[T any](val T, opts ...dig.ProvideOption) ConstructorWithOpts {
	return ConstructorWithOpts{
		Constructor: func() T { return val },
		Options:     opts,
	}
}

// ProvideWithArg will create a constructor with a first arg explicitly provided
// supposed return no error.
func ProvideWithArg[
	TArg any,
	TConstructorArg any,
	TRes any,
](
	arg TArg,
	constructor func(arg TArg, cArg TConstructorArg) TRes,
) func(TConstructorArg) TRes {
	return func(cArg TConstructorArg) TRes {
		return constructor(arg, cArg)
	}
}

// ProvideWithArgErr will create a constructor with a first arg explicitly provided
// supposed return an error.
func ProvideWithArgErr[
	TArg any,
	TConstructorArg any,
	TDep any,
](
	arg TArg,
	constructor func(arg TArg, cArg TConstructorArg) (TDep, error),
) func(TConstructorArg) (TDep, error) {
	return func(cArg TConstructorArg) (TDep, error) {
		return constructor(arg, cArg)
	}
}
