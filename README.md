# flickr.go

A go library to easily consume Flickr API.
The project is currently under heavy development, so API coverage should get better on a daily basis.

[![GoDoc](https://godoc.org/github.com/masci/flickr.go?status.svg)](https://godoc.org/github.com/masci/flickr.go)
[![Build Status](https://travis-ci.org/masci/flickr.go.svg)](https://travis-ci.org/masci/flickr.go)
[![Coverage Status](https://coveralls.io/repos/masci/flickr.go/badge.svg)](https://coveralls.io/r/masci/flickr.go)


## Extra-API Methods
 * Get OAuth request token
 * Get OAuth authorize URL
 * Get OAuth access token
 * Upload photo

## API Methods

### auth.oauth
 * flickr.auth.oauth.checkToken

### photos (in progress, 1/26)
 * ~~flickr.photos.addTags~~
 * flickr.photos.delete
 * ~~flickr.photos.getAllContexts~~
 * ~~flickr.photos.getContactsPhotos~~
 * ~~flickr.photos.getContactsPublicPhotos~~
 * ~~flickr.photos.getContext~~
 * ~~flickr.photos.getCounts~~
 * ~~flickr.photos.getExif~~
 * ~~flickr.photos.getFavorites~~
 * ~~flickr.photos.getInfo~~
 * ~~flickr.photos.getNotInSet~~
 * ~~flickr.photos.getPerms~~
 * ~~flickr.photos.getRecent~~
 * ~~flickr.photos.getSizes~~
 * ~~flickr.photos.getUntagged~~
 * ~~flickr.photos.getWithGeoData~~
 * ~~flickr.photos.getWithoutGeoData~~
 * ~~flickr.photos.recentlyUpdated~~
 * ~~flickr.photos.removeTag~~
 * ~~flickr.photos.search~~
 * ~~flickr.photos.setContentType~~
 * ~~flickr.photos.setDates~~
 * ~~flickr.photos.setMeta~~
 * ~~flickr.photos.setPerms~~
 * ~~flickr.photos.setSafetyLevel~~
 * ~~flickr.photos.setTags~~

### photosets (in progress, 4/14)
 * flickr.photosets.addPhoto
 * flickr.photosets.create
 * flickr.photosets.delete
 * ~~flickr.photosets.editMeta~~
 * ~~flickr.photosets.editPhotos~~
 * ~~flickr.photosets.getContext~~
 * ~~flickr.photosets.getInfo~~
 * flickr.photosets.getList
 * ~~flickr.photosets.getPhotos~~
 * ~~flickr.photosets.orderSets~~
 * ~~flickr.photosets.removePhoto~~
 * ~~flickr.photosets.removePhotos~~
 * ~~flickr.photosets.reorderPhotos~~
 * ~~flickr.photosets.setPrimaryPhoto~~

### test
 * flickr.test.echo
 * flickr.test.login
 * flickr.test.null

