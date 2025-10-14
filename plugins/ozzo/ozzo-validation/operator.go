// Copyright 2025 samber.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://github.com/samber/ro/blob/main/licenses/LICENSE.apache.md
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


package roozzovalidation

import (
	"context"
	"errors"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/samber/ro"
)

var (
	ErrValidatable            = errors.New("value does not implement ozzo.Validatable")
	ErrValidatableWithContext = errors.New("value does not implement ozzo.ValidatableWithContext")
)

func Validate[T any](rules ...ozzo.Rule) func(ro.Observable[T]) ro.Observable[Result[T]] {
	return ro.Map(func(v T) Result[T] {
		err := ozzo.Validate(v, rules...)
		if err != nil {
			return Err[T](err)
		}
		return Ok(v)
	})
}

func ValidateStruct[T any]() func(ro.Observable[T]) ro.Observable[Result[T]] {
	var t T
	if _, ok := any(t).(ozzo.Validatable); !ok {
		panic(ErrValidatable)
	}

	return ro.Map(func(v T) Result[T] {
		err := any(v).(ozzo.Validatable).Validate()
		if err != nil {
			return Err[T](err)
		}
		return Ok(v)
	})
}

func ValidateWithContext[T any](rules ...ozzo.Rule) func(ro.Observable[T]) ro.Observable[Result[T]] {
	return ro.MapWithContext(func(ctx context.Context, v T) (context.Context, Result[T]) {
		err := ozzo.ValidateWithContext(ctx, v, rules...)
		if err != nil {
			return ctx, Err[T](err)
		}
		return ctx, Ok(v)
	})
}

func ValidateStructWithContext[T any]() func(ro.Observable[T]) ro.Observable[Result[T]] {
	var t T
	if _, ok := any(t).(ozzo.ValidatableWithContext); !ok {
		panic(ErrValidatableWithContext)
	}

	return ro.MapWithContext(func(ctx context.Context, v T) (context.Context, Result[T]) {
		err := any(v).(ozzo.ValidatableWithContext).ValidateWithContext(ctx)
		if err != nil {
			return ctx, Err[T](err)
		}
		return ctx, Ok(v)
	})
}

func ValidateOrError[T any](rules ...ozzo.Rule) func(ro.Observable[T]) ro.Observable[T] {
	return ro.MapErr(func(v T) (T, error) {
		err := ozzo.Validate(v, rules...)
		return v, err
	})
}

func ValidateStructOrError[T any]() func(ro.Observable[T]) ro.Observable[T] {
	var t T
	if _, ok := any(t).(ozzo.Validatable); !ok {
		panic(ErrValidatable)
	}

	return ro.MapErr(func(v T) (T, error) {
		err := any(v).(ozzo.Validatable).Validate()
		return v, err
	})
}

func ValidateOrErrorWithContext[T any](rules ...ozzo.Rule) func(ro.Observable[T]) ro.Observable[T] {
	return ro.MapErrWithContext(func(ctx context.Context, v T) (T, context.Context, error) {
		err := ozzo.ValidateWithContext(ctx, v, rules...)
		return v, ctx, err
	})
}

func ValidateStructOrErrorWithContext[T any]() func(ro.Observable[T]) ro.Observable[T] {
	var t T
	if _, ok := any(t).(ozzo.ValidatableWithContext); !ok {
		panic(ErrValidatableWithContext)
	}

	return ro.MapErrWithContext(func(ctx context.Context, v T) (T, context.Context, error) {
		err := any(v).(ozzo.ValidatableWithContext).ValidateWithContext(ctx)
		return v, ctx, err
	})
}

func ValidateOrSkip[T any](rules ...ozzo.Rule) func(ro.Observable[T]) ro.Observable[T] {
	return ro.Filter(func(v T) bool {
		err := ozzo.Validate(v, rules...)
		return err == nil
	})
}

func ValidateStructOrSkip[T any]() func(ro.Observable[T]) ro.Observable[T] {
	var t T
	if _, ok := any(t).(ozzo.Validatable); !ok {
		panic(ErrValidatable)
	}

	return ro.Filter(func(v T) bool {
		err := any(v).(ozzo.Validatable).Validate()
		return err == nil
	})
}

func ValidateOrSkipWithContext[T any](rules ...ozzo.Rule) func(ro.Observable[T]) ro.Observable[T] {
	return ro.FilterWithContext(func(ctx context.Context, v T) (context.Context, bool) {
		err := ozzo.ValidateWithContext(ctx, v, rules...)
		return ctx, err == nil
	})
}

func ValidateStructOrSkipWithContext[T any]() func(ro.Observable[T]) ro.Observable[T] {
	var t T
	if _, ok := any(t).(ozzo.ValidatableWithContext); !ok {
		panic(ErrValidatableWithContext)
	}

	return ro.FilterWithContext(func(ctx context.Context, v T) (context.Context, bool) {
		err := any(v).(ozzo.ValidatableWithContext).ValidateWithContext(ctx)
		return ctx, err == nil
	})
}
