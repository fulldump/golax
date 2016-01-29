
<!-- MarkdownTOC autolink=true bracket=round depth=4 -->

- [Implementation decisions](#implementation-decisions)
    - [Decision #1](#decision-1)

<!-- /MarkdownTOC -->

# Implementation decisions

This part cover some of the implementation decisions taken along the development process.

## Decision #1

Handler functions has 3 parameters:

```
func (w http.ResponseWriter, r *http.Request, c *lax.Context) {
    
}
```

* w - Response writer (from standard http library)
* r - Request  (from standard http library)
* c - The context object, allow pass data between middlewares and methods

Why not embed `w` and `r` into `c` and have a simpler function signature?

```
func (c *lax.Context) {
    
}
```

The main reason is to maintain familiarity between the old handlers and Lax. You could copy & paste the code 

The secondary reason is code readability. Since `w` and `r` are used in all (or almost all) cases, `w.Write(...)` is more readable than `c.w.Write(...)`.

