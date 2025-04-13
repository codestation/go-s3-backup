/*
Copyright 2025 codestation

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

package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCensor(t *testing.T) {
	r := require.New(t)

	long := []string{"--foo", "--bar", "1234", "--baz"}
	short := []string{"--foo", "-b=1234", "--baz"}

	res := censorArg(long, "--bar")
	r.NotContains(res, "1234")

	res = censorArg(long, "--none")
	r.Equal(res, long)

	res = censorArg(short, "-b")
	r.NotContains(res, "-b=1234")

	res = censorArg(short, "-c")
	r.Contains(res, "-b=1234")

	res = censorArg(long, "")
	r.Equal(res, long)
}
