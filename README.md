# flickr

A go library to easily consume Flickr API.
The project is currently under heavy development, so it hasn't a version number yet.

[![GoDoc](https://godoc.org/gopkg.in/masci/flickr.v2?status.svg)](https://godoc.org/gopkg.in/masci/flickr.v2)
[![Build Status](https://travis-ci.org/masci/flickr.svg)](https://travis-ci.org/masci/flickr)
[![Coverage Status](https://coveralls.io/repos/masci/flickr/badge.svg)](https://coveralls.io/r/masci/flickr)

## Usage

`flickr` aims to expose a Go Api matching Flickr REST Api, so that you don't need
to build HTTP requests and parse HTTP response manually. For example, the Flickr
method `flickr.photosets.create` is implemented with the `Create` function in the `flickr/photosets`
package:

```go
import "fmt"
import "gopkg.in/masci/flickr.v2"
import "gopkg.in/masci/flickr.v2/photosets"

// create an API client with credentials
client := flickr.NewFlickrClient("your_apikey", "your_apisecret")
client.OAuthToken = "your_token"
client.OAuthTokenSecret = "your_tokenSecret"

response, _ := photosets.Create(client, "My Set", "Description", "primary_photo_id")
fmt.Println("New photoset created:", response.Photoset.Id)
```

`flickr` responses implement `flickr.FlickrResponse` interface. A response contains error codes
and error messages (if any) produced by Flickr or the specific data returned by the api call.
Different methods may return different kind of responses.

### Upload a photo

There are a number of functions that don't map any actual Flickr Api method
(see below for the detailed list). For example, to upload a photo, you call the
`UploadFile` or `UploadReader` functions in the `flickr` package:

```go
import "gopkg.in/masci/flickr.v2"


// upload the image file with default (nil) options
resp, err := flickr.UploadFile(client, "/path/to/image", nil)
```
Files are uploaded through an io.Pipe fueled in a separate goroutine, so the process is pretty efficient.

### Authentication (or how to retrieve OAuth credentials)

Several api calls must be authenticated and authorized: `flickr` only supports
OAuth since the original token-based method has been deprecated by Flickr. This
is an example describing the OAuth worflow from a command line application:

```go
import "gopkg.in/masci/flickr.v2"

client := flickr.NewFlickrClient("your_apikey", "your_apisecret")

// first, get a request token
requestTok, _ := flickr.GetRequestToken(client)

// build the authorizatin URL
url, _ := flickr.GetAuthorizeUrl(client, requestTok)

// ask user to hit the authorization url with
// their browser, authorize this application and coming
// back with the confirmation token

// finally, get the access token, setup the client and start making requests
accessTok, err := flickr.GetAccessToken(client, requestTok, "oauth_confirmation_code")
client.OAuthToken = accessTok.OAuthToken
client.OAuthTokenSecret = accessTok.OAuthTokenSecret
```

### Api coverage

Only a small part of the Flickr Api is implemented as Go functions: even if it's quite
simple to write the code for the mapping, I only did it for methods I actually need in my projects
(contributions well accepted). Anyway, if you need to call a Flickr Api method that wasn't
already mapped, you can do it manually:

```go
import "fmt"
import "gopkg.in/masci/flickr.v2"

client := flickr.NewFlickrClient("your_apikey", "your_apisecret")
client.Init()
client.Args.Set("method", "flickr.cameras.getBrandModels")
client.Args.Set("brand", "nikon")

client.OAuthSign()
response := &flickr.BasicResponse{}
err := flickr.DoGet(client, response)

if err != nil {
    fmt.Printf("Error: %s", err)
} else {
    fmt.Println("Api response:", response.Extra)
}
```

Checkout the `example` folder and the docs pages for more details.

## Note on Go versions

The latest version `v2` only supports go `1.6` and above, for Go `< 1.6` use the `v1` package:
```
go get gopkg.in/masci/flickr.v1
```

## API Methods

### Extra-API Methods
These are methods that are not actually part of the Flickr API

 * Get OAuth request token
 * Get OAuth authorize URL
 * Get OAuth access token
 * Upload photo

### auth.oauth
 * flickr.auth.oauth.checkToken

### photos
 * flickr.photos.delete
 * flickr.photos.getInfo
 * flickr.photos.setDates

### photosets
 * flickr.photosets.addPhoto
 * flickr.photosets.create
 * flickr.photosets.delete
 * flickr.photosets.editMeta
 * flickr.photosets.editPhotos
 * flickr.photosets.getInfo
 * flickr.photosets.getList
 * flickr.photosets.getPhotos
 * flickr.photosets.orderSets
 * flickr.photosets.removePhoto
 * flickr.photosets.removePhotos
 * flickr.photosets.reorderPhotos
 * flickr.photosets.setPrimaryPhoto

### people
 * flickr.people.getPhotos

### test
 * flickr.test.echo
 * flickr.test.login
 * flickr.test.null
