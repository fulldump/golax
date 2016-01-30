
<!-- MarkdownTOC autolink=true bracket=round depth=4 -->

- [Implementation decisions](#implementation-decisions)
    - [Decision #1](#decision-1)

<!-- /MarkdownTOC -->

# Implementation decisions

This part cover some of the implementation decisions taken along the development process.

## Decision #1

Handler functions has 1 parameter:

```
func (c *lax.Context) {
    
}
```

Why not `w`, `r` and `c` and maintain developer compatibility?

We would ended up with the following signature:

```
func (w http.ResponseWriter, r *http.Request, c *lax.Context) {
    
}
```

Old code is not going to work by doing copy&paste, but you only have to replace:

* `w` by `c.Response`
* `r` by `c.Request`

Making this decision is hard but `c *lax.Context` is much easier to remember.

About code readability, `w.Write(...)` is shorter but `c.Response.Write(...)` is more semantic.

