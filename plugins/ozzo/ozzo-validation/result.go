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

// Monad for validation results, inspired by github.com/samber/mo.Result
type Result[T any] struct {
	isErr bool
	value T
	err   error
}

func (r Result[T]) Unwrap() T {
	return r.value
}

func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

func (r Result[T]) IsOk() bool {
	return !r.isErr
}

func (r Result[T]) IsError() bool {
	return r.isErr
}

func (r Result[T]) Error() error {
	return r.err
}

func (r Result[T]) Get() (T, error) {
	if r.isErr {
		var t T
		return t, r.err
	}

	return r.value, nil
}

func Ok[T any](value T) Result[T] {
	return Result[T]{
		isErr: false,
		value: value,
		err:   nil,
	}
}

func Err[T any](err error) Result[T] {
	return Result[T]{
		isErr: true,
		err:   err,
	}
}
