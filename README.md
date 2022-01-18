[![GoDoc](https://godoc.org/github.com/KarpelesLab/contexter?status.svg)](https://godoc.org/github.com/KarpelesLab/contexter)

# Contexter

A dirty piece of code for fetching `context.Context` from the stack when it
cannot be passed normally.

## But why?

I turned out to need access to context.Context from a `MarshalJSON()` method
and found out there is no way to pass it easily. This is one of my attempts
at passing ctx across the stack.

So if your method that accepts a `context.Context` calls `json.Marshal()`, it
is likely the `MarshalJSON()` methods will be able to fetch the context using
`contexter.Context()`.

### My soul doesn't hurt enough

More stack based golang dark magic can be found online:

* https://github.com/jtolio/gls

## It doesn't work

There can be various reasons why this doesn't work, or stopped working.

* The method receiving a `context.Context` has been inlined and its parameters aren't available
* The variable containing the `context.Context` isn't used and was overwritten
* You're into a different goroutine (which means a new stack)
* Something changed in Go's runtime, stack display format, etc
* The architecture you're using does something that's not accounted for in this module

This is not an exact science, and using this package may result in the end of
the world, or your program crashing in ways go can't recover. You've been
warned.

# Usage

```go
ctx := contexter.Context()
```

This will automatically fetch the closest context.Context object found in the
stack that was passed as a context.Context object, or nil if none were found.

```go
var ctx context.Context
if contexter.Find(&ctx) {
	// use ctx
}
```

This alternative version can find other kind of interfaces on the stack.
