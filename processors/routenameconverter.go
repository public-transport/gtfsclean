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
	MoveHeadsignRegex  *regexp.Regexp
}

func (f RouteNameConverter) Run(feed *gtfsparser.Feed) {
	fmt.Fprintf(os.Stdout, "Copying selected trip names and headsigns to route names...")

	for _, trip := range feed.Trips {
		route := trip.Route
		if route == nil {
			continue
		}
		if route.Short_name != "" &&
			f.KeepRouteNameRegex != nil &&
			f.KeepRouteNameRegex.MatchString(route.Short_name) {
			continue
		} else {
			newName := ""

			if trip.Short_name != nil &&
				*trip.Short_name != "" &&
				f.CopyTripNameRegex != nil &&
				f.CopyTripNameRegex.MatchString(*trip.Short_name) {
				newName = *trip.Short_name
			}

			if trip.Headsign != nil &&
				*trip.Headsign != "" &&
				f.MoveHeadsignRegex != nil &&
				f.MoveHeadsignRegex.MatchString(*trip.Headsign) {
				newName = *trip.Headsign

				replacementHeadsign := ""
				if len(trip.StopTimes) > 0 {
					replacementHeadsign = trip.StopTimes[len(trip.StopTimes)-1].Stop().Name
				}
				trip.Headsign = &replacementHeadsign
			}

			if newName != "" {
				newRouteId := route.Id + trip.Id
				newRoute := (*route)
				newRoute.Id = newRouteId
				newRoute.Short_name = newName
				feed.Routes[newRouteId] = &newRoute
				trip.Route = &newRoute
			}
		}
	}

	fmt.Fprintf(os.Stdout, " done.\n")
}
