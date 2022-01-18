[![GoDoc](https://godoc.org/github.com/KarpelesLab/contexter?status.svg)](https://godoc.org/github.com/KarpelesLab/contexter)

# Contexter

A dirty piece of code for fetching context.Context from the stack when it
cannot be passed normally.

This can be useful for example when called in a MarshalJSON() method and having
context information is needed to marshal the object right, but the json package
provides no way to pass a context.

## It doesn't work

There can be various reasons why this doesn't work. Either golang inlined the
function call (in which case parameters are not on the stack), or the variable
isn't used and go used the memory for something else. This is not an exact
science, and using this package may result in the end of the world. You've been
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
