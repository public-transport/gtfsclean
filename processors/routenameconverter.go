// Copyright 2025 Jonah Brüchert
// Authors: jbb@kaidan.im
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package processors

import (
	"fmt"
	"os"
	"regexp"

	"github.com/public-transport/gtfsparser"
)

// RouteNameConverter copies the trip_short_name field into route_short_name based on a set of rules
type RouteNameConverter struct {
	KeepRouteNameRegex *regexp.Regexp
	CopyTripNameRegex  *regexp.Regexp
}

func (f RouteNameConverter) Run(feed *gtfsparser.Feed) {
	fmt.Fprintf(os.Stdout, "Copying selected trip names to route names...")

	for _, trip := range feed.Trips {
		route := trip.Route
		if route == nil {
			continue
		}
		if route.Short_name != "" && f.KeepRouteNameRegex != nil && f.KeepRouteNameRegex.MatchString(route.Short_name) {
			continue
		} else {
			if trip.Short_name != nil && *trip.Short_name != "" && f.CopyTripNameRegex != nil && f.CopyTripNameRegex.MatchString(*trip.Short_name) {
				newRouteId := route.Id + trip.Id
				// todo copy
				newRoute := (*route)
				newRoute.Id = newRouteId
				newRoute.Short_name = *trip.Short_name
				feed.Routes[newRouteId] = &newRoute
				trip.Route = &newRoute
			}
		}
	}

	fmt.Fprintf(os.Stdout, " done.")
}
