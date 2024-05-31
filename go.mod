module github.com/public-transport/gtfsclean

go 1.20

require (
	github.com/paulmach/go.geojson v1.5.0
	github.com/public-transport/gtfsparser v0.0.0-20240508211334-a2b010870a5e
	github.com/public-transport/gtfswriter v0.0.0-20240530234004-bf8f5e60799e
	github.com/spf13/pflag v1.0.5

	// Remove this once our minimum Go version is forced to be 1.21
	// This is only used for slices which are added in 1.21 stdlib
	golang.org/x/exp v0.0.0-20240530194437-404ba88c7ed0
)

require (
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/twotwotwo/sorts v0.0.0-20160814051341-bf5c1f2b8553 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
)
