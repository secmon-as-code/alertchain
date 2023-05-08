# HTTP

## `http.fetch`

This action fetches data from the specified URL using the specified HTTP method and returns the result.

### Arguments

Example policy:

```rego
run[res] {
  res := {
    id: "your-action",
    uses: "http.fetch",
    args: {
      "method": "GET",
      "url": "https://api.example.com/data",
      "header": {
        "Authorization": "Bearer " + input.env.API_TOKEN,
      },
    },
  },
}
```

- `method` (string, required): Specifies the HTTP method to use for the request (e.g., "GET", "POST", "PUT", "DELETE").
- `url` (string, required): Specifies the URL to fetch the data from.
- `header` (map of strings, optional): Specifies the HTTP headers to include in the request.
- `data` (string, optional): Specifies the request body to include in the request (for methods like "POST" and "PUT").

### Response

The response depends on the `Content-Type` of the HTTP response. The following content types are supported:

- `application/json`: The response body will be parsed as JSON and returned as an object.
- `application/octet-stream`: The response body will be returned as a byte array.
- Other content types: The response body will be returned as a string.

In case of an error while making the HTTP request or processing the response, an error will be returned.