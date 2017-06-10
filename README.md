# instagram-api-bypass

This repository parses Instagram's public response to GET requests, providing an API-free way to obtain basic account and media information.

## Installation

```
go get github.com/andrewfrench/instagram-api-bypass
```

## Use
### Get user account information

```
import "github.com/andrewfrench/instagram-api-bypass/account"

func main() {
    acc, err := account.Get("my_instagram_username")
}
```

### Get media information

```
import "github.com/andrewfrench/instagram-api-bypass/media"

func main() {
    med, err := media.Get("media_shortcode")
}
```