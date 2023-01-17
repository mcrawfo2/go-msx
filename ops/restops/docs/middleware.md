# Middleware

You can use a Mediator/Middleware component to augment the functionality of 
go-msx components such as endpoints:

> Middleware is software that's assembled into an app pipeline to handle requests and responses.
> Each component:
> - chooses whether to pass execution to the next component in the pipeline.
> - can perform work before and after the next component in the pipeline.
>
> -- [ASP.NET Core Middleware](https://learn.microsoft.com/en-us/aspnet/core/fundamentals/middleware/?tabs=aspnetcore2x&view=aspnetcore-7.0), Microsoft

Go HTTP Middleware factories implement a de facto function signature:

```go
type Middleware func(next http.Handler) http.Handler

func myMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Middleware BEFORE logic goes here...
        next.ServeHTTP(w, r)
        // Middleware AFTER logic goes here...
    })
}
```

The factory accepts the subsequent Handler in the middleware chain, and returns a new
Handler which wraps it with the desired added functionality.  You can find more details and
examples on [this blog post](https://www.alexedwards.net/blog/making-and-using-middleware).

## Available HTTP Middleware

go-msx does not currently define any HTTP Middleware, however there are many libraries available,
including:

- [github.com/gorilla/handlers](https://github.com/gorilla/handlers)
- [github.com/throttled/throttled/v2](https://github.com/throttled/throttled)
