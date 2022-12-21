# Analytics

[![Go Reference](https://pkg.go.dev/badge/github.com/facily-tech/go-core/analytics.svg)](https://pkg.go.dev/github.com/facily-tech/go-core/analytics)

The purpose of the analytics package is to create a simple interface to use mixpanel
event tracking feature.

You can learn more about mixpanel at https://wwww.mixpanel.com

## Usage

### Default analytics

This require you to have a environment variable ANALYTICS_MIXPANEL_TOKEN for 
supplying token and api url will default to api.mixpanel.com.

```go
package main

import (
	"log"

	"github.com/facily-tech/go-core/analytics"
)

func main() {
	if err := analytics.TrackSync("sample event", map[string]interface{}{"id": 123}); err != nil {
		log.Println(err)
	}
}
```

### Custom analytics

You may want to change token, logger, url etc.

```go
package main

import (
	"log"

	"github.com/facily-tech/go-core/analytics"
)

func main() {
	a := &analytics.Analytics{}
	a.WithMixpanelURL("no token", "")

	if err := a.TrackSync("sample event", map[string]interface{}{"id": 123}); err != nil {
		log.Println(err)
	}
}
```