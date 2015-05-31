I/O With Go: io.Pipe()
May 30, 2015

I find that Go's [I/O framework](https://golang.org/pkg/io/) is one of its major strengths:
 * The [`io.Reader`](https://golang.org/pkg/io/#Reader) and [`io.Writer`](https://golang.org/pkg/io/#Writer) abstractions make it easy to create composable programs
 * It's a great example of how to use interfaces in your own programs

One of my recent discoveries is [`io.Pipe()`](https://golang.org/pkg/io/#Pipe).

Let's for example encode some JSON and send it as an HTTP POST body. You could use a [`bytes.Buffer`](https://golang.org/pkg/bytes/#Buffer) to store the result of the encoding and then pass it as the HTTP POST body:

**BEFORE**

```go
package main

import (
  "bytes"
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"
)

type msg struct {
  Text string
}

func handleErr(err error) {
  if err != nil {
    log.Fatalf("%s\n", err)
  }
}

func main() {
  m := msg{Text: "brought to you by bytes.Buffer"}
  var buf bytes.Buffer
  err := json.NewEncoder(&buf).Encode(&m)
  handleErr(err)

  resp, err := http.Post("https://httpbin.org/post", "application/json", &buf)
  handleErr(err)
  defer resp.Body.Close()

  b, err := ioutil.ReadAll(resp.Body)
  handleErr(err)

  log.Printf("%s\n", b)
}
```

This is easy to understand but we are **allocating an unnecessary buffer** which is the kind of pattern that can become a problem at scale. `io.Pipe` allows you to eliminate the temporary buffer and connect the JSON encoder directly to the HTTP POST:

**AFTER**

```go
package main

import (
  "encoding/json"
  "io"
  "io/ioutil"
  "log"
  "net/http"
)

type msg struct {
  Text string
}

func handleErr(err error) {
  if err != nil {
    log.Fatalf("%s\n", err)
  }
}

// use a io.Pipe to connect a JSON encoder to an HTTP POST: this way you do
// not need a temporary buffer to store the JSON bytes
func main() {
  r, w := io.Pipe()

  // writing without a reader will deadlock so write in a goroutine
  go func() {
    // it is important to close the writer or reading from the other end of the
    // pipe will never finish
    defer w.Close()

    m := msg{Text: "brought to you by io.Pipe()"}
    err := json.NewEncoder(w).Encode(&m)
    handleErr(err)
  }()

  resp, err := http.Post("https://httpbin.org/post", "application/json", r)
  handleErr(err)
  defer resp.Body.Close()

  b, err := ioutil.ReadAll(resp.Body)
  handleErr(err)

  log.Printf("%s\n", b)
}
```

Of course in this trivial example it is overkill to use `io.Pipe`. But when the buffers are getting larger and you have lots of goroutines doing this kind of stuff concurrently **`io.Pipe` can help you reduce memory usage**!
