# httpx
Lightweight and pluggable HTTP client for Go.

## Example

```go
client := httpx.NewClient("https://httpbin.org", httpx.WithLogger(printLogger{}))
var result map[string]any
err := client.Get("/get", map[string]string{"q": "hello"}, &result); 
if err != nil {
   panic(err)
}
```
