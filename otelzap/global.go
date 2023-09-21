/*
Copyright Agoda Services Co.,Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package otelzap

import (
	"context"
	"go.uber.org/zap"
)

// L returns the global Logger
func L() *Logger {
	return &Logger{
		zap.L(),
	}
}

func S() *SugaredLogger {
	return L().Sugar()
}

// Ctx is a shortcut for L().Ctx(ctx).
func Ctx(ctx context.Context) *Logger {
	return L().Ctx(ctx)
}
