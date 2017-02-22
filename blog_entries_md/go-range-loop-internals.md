Go Range Loop Internals
February 22, 2017

While they are very convenient, I always found Go's range loops a bit mystifying. I'm not alone in this:
<blockquote class="twitter-tweet" data-lang="en"><p lang="en" dir="ltr"><a href="https://twitter.com/hashtag/golang?src=hash">#golang</a> pop quiz: does this program terminate?<br><br>func main() {<br>   v := []int{1, 2, 3}<br>   for i := range v {<br>       v = append(v, i)<br>  }<br>}</p>&mdash; Dαve Cheney (@davecheney) <a href="https://twitter.com/davecheney/status/819759166617108481">January 13, 2017</a></blockquote>
<blockquote class="twitter-tweet" data-lang="en"><p lang="en" dir="ltr">Today&#39;s <a href="https://twitter.com/hashtag/golang?src=hash">#golang</a> gotcha: the two-value range over an array does a copy.  Avoid by ranging over the pointer instead.<a href="https://t.co/SbK667osvA">https://t.co/SbK667osvA</a></p>&mdash; Damian Gryski (@dgryski) <a href="https://twitter.com/dgryski/status/816226596835225600">January 3, 2017</a></blockquote>
<script async src="//platform.twitter.com/widgets.js" charset="utf-8"></script>

Now I could take these facts and try to remember them but it is likely I will forget. To have a better chance of remembering this I'll need to find out _why_ the range loop behaves the way it does. So here goes.

## Step 1: RTFM

The first order of business is to go read the range loop documentation. The Go language spec documents range loops in [the for statement section](https://golang.org/ref/spec#For_statements) under "For statements with `range` clause". I will not copy the entire spec here, but let me summarise some of the interesting bits.

First, let's remind ourselves what we're looking at here:
```
for i := range a {
    fmt.Println(i)
}
```

### The range variables

Most of you will know that to the left of the range clause (`i` in the example above) you can assign the loop variables using:

* assigment (`=`) 
* short variable declaration (`:=`)

You can also choose to put nothing to ignore the loop variables altogether.

If you use the short variable declaration style assignment (`:=`), Go will reuse the variables for each iteration of the loop (only in scope inside the loop).

### The range expression

On the right side of the range clause (`a` in the example above) you find what they call the _range expression_. It can hold any expression that evaluates to one of the following:

 * array
 * pointer to an array
 * slice
 * string
 * map
 * channels that allow receiving, e.g. `chan int` or `chan<- int`

**The range expression is evaluated once before beginning the loop**. Note that there is an exception to this rule: if you range over an array (or pointer to) and you only assign the index: then only `len(a)` is evaluated. Evaluating just `len(a)` means that the expression `a` may be evaluated at compile time and replaced with a constant by the compiler. [The spec for the `len` function](https://golang.org/ref/spec#Length_and_capacity) explains:

>The expressions len(s) and cap(s) are constants if the type of s is an array or pointer to an array and the expression s does not contain channel receives or (non-constant) function calls; in this case s is not evaluated. Otherwise, invocations of len and cap are not constant and s is evaluated.

Now **what  exactly does "evaluated" mean?** Unfortunatly I can't find this information in the spec. Of course I can guess that it means to execute the expression completely until it can not be reduced further. In any case, the high order bit here is that the range expression is evaluated _once_ before the beginning of the loop. **How would you evaluate an expression just once? By assigning it to a variable!** Could that be what is happening here?

Interestingly the spec mentions something specific about adding and removing to/from maps (no mention of slices):

>If map entries that have not yet been reached are removed during iteration, the corresponding iteration values will not be produced. If map entries are created during iteration, that entry may be produced during the iteration or may be skipped. 

I will come back to maps later.

## Step 2: Data types supported by range

If we assume for a minute that the _range expression_ gets assigned to a variable once before the start of the loop, what does that mean? The answer is that it depends on the data type, so let's have a closer look at the data types supported by `range`.

Before we do that, do remember this: **in Go, everything you assign, you copy**. If you assign a pointer, you copy the pointer. If you assign a struct, you copy the struct. The same is true when passing arguments to a function. Anyway, here goes:

| type    | syntactic sugar for                                     |
|---------|---------------------------------------------------------|
| array   |the array                                                |
| string  |struct holding len + a pointer to the backing array      |
| slice   |struct holding len, cap + a pointer to the backing array |
| map     |pointer to a struct                                      |
| channel |pointer to a struct                                      |

See the references at the bottom of this post to learn more about the internal structure of these data types.

So what does this mean? These examples highlight some of the differences:

```
// copies the entire array
var a [10]int
acopy := a 

// copies the slice header struct only, NOT the backing array
s := make([]int, 10)
scopy := s

// copies the map pointer only
m := make(map[string]int)
mcopy := m
```

So if at the beginning of a `range` loop you would assign an array expression to a variable (to ensure it evaluates only once), you would be copying the entire array. We might be onto something here.

## Step 3: Go compiler source code

Lazy as I am I simply Googled for the Go compiler source. The first thing I found was the GCC version of the compiler. The interesting bits, as far as the `range` clause is concerned, are in `statements.cc`, [as in this comment](https://github.com/golang/gofrontend/blob/e387439bfd24d5e142874b8e68e7039f74c744d7/go/statements.cc#L5384):

```
// Arrange to do a loop appropriate for the type.  We will produce
//   for INIT ; COND ; POST {
//           ITER_INIT
//           INDEX = INDEX_TEMP
//           VALUE = VALUE_TEMP // If there is a value
//           original statements
//   }
```

Now we're getting somewhere. Not entirely unsurprising, range loops are just syntactic sugar for C-style loops internally. The code has specific "desugarings" for earch type supported by range. For example, [arrays](https://github.com/golang/gofrontend/blob/e387439bfd24d5e142874b8e68e7039f74c744d7/go/statements.cc#L5501):

```
// The loop we generate:
//   len_temp := len(range)
//   range_temp := range
//   for index_temp = 0; index_temp < len_temp; index_temp++ {
//           value_temp = range_temp[index_temp]
//           index = index_temp
//           value = value_temp
//           original body
//   }
```

[Slices](https://github.com/golang/gofrontend/blob/e387439bfd24d5e142874b8e68e7039f74c744d7/go/statements.cc#L5593):

```
//   for_temp := range
//   len_temp := len(for_temp)
//   for index_temp = 0; index_temp < len_temp; index_temp++ {
//           value_temp = for_temp[index_temp]
//           index = index_temp
//           value = value_temp
//           original body
//   }
```

The common theme here is that

* **Everything is just a C-style for loop**
* **The thing you iterate over is assigned to a temporary variable**

That said, this was in the GCC frontend. Most people I know use the gc compiler that comes with the Go distribution. It looks like that compiler does [pretty much the same thing](https://github.com/golang/go/blob/ea020ff3de9482726ce7019ac43c1d301ce5e3de/src/cmd/compile/internal/gc/range.go#L169).

## What we know

1. Loop variables are reused and assigned to at each iteration.
2. The range expression gets evaluated once before the loop starts by assigning to a variable.
3. You can delete or add values to a map while iterating. Adds may or may not be visible in the loop.

With this in hand, let's go back to the examples listed at the beginning of this post

## Dave's tweet

<blockquote class="twitter-tweet" data-lang="en"><p lang="en" dir="ltr"><a href="https://twitter.com/hashtag/golang?src=hash">#golang</a> pop quiz: does this program terminate?<br><br>func main() {<br>   v := []int{1, 2, 3}<br>   for i := range v {<br>       v = append(v, i)<br>  }<br>}</p>&mdash; Dαve Cheney (@davecheney) <a href="https://twitter.com/davecheney/status/819759166617108481">January 13, 2017</a></blockquote>

The reason this terminates is because it roughly translates to something like this:

```
for_temp := v
len_temp := len(for_temp)
for index_temp = 0; index_temp < len_temp; index_temp++ {
        value_temp = for_temp[index_temp]
        index = index_temp
        value = value_temp
        v = append(v, index)
}
```

We know slices are syntactic sugar for a struct holding a pointer to the backing array. The loop is iterating over `for_temp`, which is a copy of that struct taken before the start of the loop. Any changes the variable `v` itself are thus not relevant as it is another copy of the struct. The backing array is still shared as it is just a pointer in that struct, so things like `v[i] = 1` would work.

## Damian's tweet
<blockquote class="twitter-tweet" data-lang="en"><p lang="en" dir="ltr">Today&#39;s <a href="https://twitter.com/hashtag/golang?src=hash">#golang</a> gotcha: the two-value range over an array does a copy.  Avoid by ranging over the pointer instead.<a href="https://t.co/SbK667osvA">https://t.co/SbK667osvA</a></p>&mdash; Damian Gryski (@dgryski) <a href="https://twitter.com/dgryski/status/816226596835225600">January 3, 2017</a></blockquote>

Again, similar to the case above, the array gets assigned to a temporary variable before the loop starts, which in the case of an array means taking a copy of the entire array. The reason it works with a pointer is that in that case the pointer will be copied, not the array.

## Extra: maps

In the spec we read that

* it is safe to add to and remove from maps in a range loop
* if you add an element it may or may not see it in an upcoming iteration

Why does it work like that? For one, we know that maps are pointers to a struct. Before the loop starts, the pointer will be copied and not the internal data structure, hence why it is possible to add or remove keys inside the loop. This makes sense!

So why might you not see the element you added in an upcoming iteration? Well if you know about how hash tables work, which is what a map really is, then you'll know that inside the backing array for a hash table the items are in no particular order. The item you add last might hash to index zero in the backing array. So if you assume that Go reserves the right to iterate over this array in any order, it is indeed impossible to predict whether or not you will see the item you added inside the loop. After all, you might already be past index zero in the backing array. This might not be exactly what happens in the case of the Go map, but it makes sense to leave the decision to the compiler writer for this reason.

## References

1. [The Go Programming Language Specification](https://golang.org/ref/spec)
2. [Go slices: usage and internals](https://blog.golang.org/go-slices-usage-and-internals)
3. [Go Data Structures](https://research.swtch.com/godata)
4. Inside the map implementation: [slides](https://docs.google.com/presentation/d/1CxamWsvHReswNZc7N2HMV7WPFqS8pvlPVZcDegdC_T4/edit#slide=id.g153a5e64a5_1_0) | [video](https://www.youtube.com/watch?v=Tl7mi9QmLns)
5. Understanding nil: [slides](https://speakerdeck.com/campoy/understanding-nil) | [video](https://www.youtube.com/watch?v=ynoY2xz-F8s) 
6. [string source code](https://golang.org/src/runtime/string.go)
7. [slice source code](https://golang.org/src/runtime/slice.go)
8. [map source code](https://golang.org/src/runtime/hashmap.go)
9. [channel source code](https://golang.org/src/runtime/chan.go)
