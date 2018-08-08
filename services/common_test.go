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
	r.Contains(res, "1234")

	res = censorArg(short, "-b")
	r.NotContains(res, "-b=1234")

	res = censorArg(short, "-c")
	r.Contains(res, "-b=1234")

	res = censorArg(long, "")
	r.Contains(res, "1234")
}
