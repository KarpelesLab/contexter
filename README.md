# Contexter

A dirty piece of code for fetchint context.Context from the stack when it
cannot be passed normally.

This can be useful for example when called in a MarshalJSON() method and having
context information is needed to marshal the object right, but the json package
provides no way to pass a context.
