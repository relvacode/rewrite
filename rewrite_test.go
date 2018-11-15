package rewrite

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

type StreamTest struct {
	Replace  string
	With     string
	Using    string
	Produces string
}

func (tc StreamTest) Run(t *testing.T) {
	t.Run(fmt.Sprintf("%s=>%s", tc.Replace, tc.With), func(t *testing.T) {
		var r = strings.NewReader(tc.Using)
		var rp = New(r, []byte(tc.Replace), []byte(tc.With))
		var buf bytes.Buffer

		_, err := io.Copy(&buf, rp)
		if err != nil {
			t.Fatal(err)
		}

		var out = buf.String()
		if out != tc.Produces {
			t.Fatalf("Bad replace: expected %q but generated %q", tc.Produces, out)
		}
	})
}

func TestStream(t *testing.T) {
	tests := []StreamTest{
		{
			Replace:  "a",
			With:     "b",
			Using:    "a",
			Produces: "b",
		},
		{
			Replace:  "hello",
			With:     "hello, world!",
			Using:    "hello",
			Produces: "hello, world!",
		},
		{
			Replace:  "a",
			With:     "b",
			Using:    "aaaaa",
			Produces: "bbbbb",
		},
		{
			Replace:  "string",
			With:     "string",
			Using:    "stringstring",
			Produces: "stringstring",
		},
		{
			Replace:  "string",
			With:     "newstring",
			Using:    "strin",
			Produces: "strin",
		},
		{
			Replace:  "b",
			With:     "a",
			Using:    "abababab",
			Produces: "aaaaaaaa",
		},
		{
			Replace:  "foo",
			With:     "bar",
			Using:    "foobarfoobarfoofo",
			Produces: "barbarbarbarbarfo",
		},
		{
			Replace:  "hello world",
			With:     "world",
			Using:    "hello=hello world=hello",
			Produces: "hello=world=hello",
		},
	}

	for _, tc := range tests {
		tc.Run(t)
	}
}

var what = "a longer string to match and replace"
var with = "some other content"

func Generate() []byte {
	var b bytes.Buffer
	for i := 0; i < 100000000; i++ {
		b.WriteString("foobar")
		if i%25 == 0 {
			b.WriteString(what)
		}
	}
	return b.Bytes()
}

func BenchmarkStream_Read(b *testing.B) {
	b.Run("vanilla", func(b *testing.B) {
		var data = Generate()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			io.Copy(ioutil.Discard, bytes.NewReader(data))
		}
	})

	b.Run("rewrite", func(b *testing.B) {
		var data = Generate()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var r = New(bytes.NewReader(data), []byte(what), []byte(with))
			io.Copy(ioutil.Discard, r)
		}

	})
}
