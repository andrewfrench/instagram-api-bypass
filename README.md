# instagram-api-bypass

This repository parses Instagram's public response to GET requests, providing an API-free way to obtain basic account and media information.

## Installation

```
./build.sh
``` 

The `iab` binary will now be in your $GOBIN (`~/go/bin` by default).  If you haven't already, add this to your $PATH to be able to use `iab` seamlessly.

## Use
### Get user account information

```
iab account my_instagram_username
```

```
import "github.com/andrewfrench/instagram-api-bypass/pkg/account"

func main() {
    acc, err := account.Get("my_instagram_username")
}
```

### Get media information

```
iab media media_shortcode
```

```
import "github.com/andrewfrench/instagram-api-bypass/pkg/media"

func main() {
    med, err := media.Get("media_shortcode")
}
```
