// Copyright 2024 Jonah Brüchert
// Authors: jbb@kaidan.im
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package processors

import (
	"fmt"
	"os"
	"github.com/public-transport/gtfsparser"
	gtfs "github.com/public-transport/gtfsparser/gtfs"
	"slices"
)

// AgencyDuplicateRemover merges semantically equivalent routes
type AgencyFilter struct {
	NamesToRemove []string
	NamesToKeep []string
}

// Run this AgencyDuplicateRemover on some feed
func (f AgencyFilter) Run(feed *gtfsparser.Feed) {
	if len(f.NamesToKeep) > 0 {
		runKeepOnly(feed, f.NamesToKeep)
	} else if len(f.NamesToRemove) > 0 {
		runRemove(feed, f.NamesToRemove)
	}
}



func runKeepOnly(feed *gtfsparser.Feed, namesToKeep []string) {
	fmt.Fprintf(os.Stdout, "Keeping only selected agencies... ")

	keepMap := make(map[string]*gtfs.Agency)
	for id, agency := range feed.Agencies {
		if slices.Contains(namesToKeep, agency.Name) {
			keepMap[id] = agency
		}
	}
	for id := range feed.Agencies {
		if _, keep := keepMap[id]; !keep {
			delete(feed.Agencies, id)
		}
	}

	keepRoutes := make(map[string]*gtfs.Route)
	for id, route := range feed.Routes {
		for _, agency := range keepMap {
			if route.Agency == agency {
				keepRoutes[id] = route
			}
		}
	}
	for id := range feed.Routes {
		if _, keep := keepRoutes[id]; !keep {
			delete(feed.Routes, id)
		}
	}

	for tripID, trip := range feed.Trips {
		if _, keep := keepRoutes[trip.Route.Id]; !keep {
			feed.DeleteTrip(tripID)
		}
	}

	for id, attr := range feed.FareAttributes {
		keep := false
		for _, agency := range keepMap {
			if attr.Agency == agency {
				keep = true
				break
			}
		}
		if !keep {
			delete(feed.FareAttributes, id)
		}
	}

	fmt.Fprintf(os.Stdout, "done. (%d agencies kept)\n", len(keepMap))
}

func runRemove(feed *gtfsparser.Feed, namesToRemove []string) {
	fmt.Fprintf(os.Stdout, "Removing selected agencies... ")
	removedAgencies := []*gtfs.Agency{}

	for agencyId := range feed.Agencies {
		agency := feed.Agencies[agencyId]
		if slices.Contains(namesToRemove, agency.Name) {
			removedAgencies = append(removedAgencies, agency)
			delete(feed.Agencies, agencyId)
		}
	}

	deletedRoutes := []*gtfs.Route{}
	for routeId := range feed.Routes {
		route := feed.Routes[routeId]
		if slices.Contains(removedAgencies, route.Agency) {
			delete(feed.Routes, routeId)
			deletedRoutes = append(deletedRoutes, route)
		}
	}

	for trip := range feed.Trips {
		if slices.Contains(deletedRoutes, feed.Trips[trip].Route) {
			feed.DeleteTrip(trip)
		}
	}

	for fareAttr := range feed.FareAttributes {
		attr := feed.FareAttributes[fareAttr]
		if slices.Contains(removedAgencies, attr.Agency) {
			delete(feed.FareAttributes, attr.Id)
		}
	}

	fmt.Fprintf(os.Stdout, "done. (-%d agencies removed)\n", len(removedAgencies))
}


