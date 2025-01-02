// Copyright 2024 Jonah Brüchert
// Authors: jbb@kaidan.im
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package processors

import (
	"github.com/public-transport/gtfsparser"
	gtfs "github.com/public-transport/gtfsparser/gtfs"
	"slices"
)

// AgencyDuplicateRemover merges semantically equivalent routes
type AgencyFilter struct {
	NamesToRemove []string
}

// Run this AgencyDuplicateRemover on some feed
func (f AgencyFilter) Run(feed *gtfsparser.Feed) {
	agencies := []*gtfs.Agency{}

	for agencyId := range feed.Agencies {
		agency := feed.Agencies[agencyId]
		if slices.Contains(f.NamesToRemove, agency.Name) {
			agencies = append(agencies, agency)

			delete(feed.Agencies, agencyId)
		}
	}

	// Drop everything referencing the agency
	for routeId := range feed.Routes {
		route := feed.Routes[routeId]
		if slices.Contains(agencies, route.Agency) {
			delete(feed.Routes, routeId)
		}

		for trip := range feed.Trips {
			if feed.Trips[trip].Route == route {
				delete(feed.Trips, trip)
			}
		}
	}

	for fareAttributes := range feed.FareAttributes {
		attribute := feed.FareAttributes[fareAttributes]
		if slices.Contains(agencies, attribute.Agency) {
			delete(feed.FareAttributes, attribute.Id)
		}
	}
}
