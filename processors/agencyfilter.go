// Copyright 2024 Jonah Brüchert
// Authors: jbb@kaidan.im
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package processors

import (
	"fmt"
	"os"
	"slices"

	"github.com/public-transport/gtfsparser"
)

// AgencyFilter removes agencies specified by command line arguments
type AgencyFilter struct {
	NamesToRemove []string
	NamesToKeep   []string
}

// Run this AgencyFilter on some feed
func (f AgencyFilter) Run(feed *gtfsparser.Feed) {
	tripsBefore := len(feed.Trips)
	agenciesBefore := len(feed.Agencies)

	removedAgencies, removedTrips := 0, 0
	if len(f.NamesToKeep) > 0 {
		removedAgencies, removedTrips = runKeepOnly(feed, f.NamesToKeep)
	} else if len(f.NamesToRemove) > 0 {
		removedAgencies, removedTrips = runRemove(feed, f.NamesToRemove)
	}

	fmt.Fprintf(os.Stdout, "done. (-%d agencies [-%.2f%%], %d trips [-%.2f%%])\n", removedAgencies, 100.0*float64(removedAgencies)/float64(agenciesBefore),
		removedTrips, 100.0*float64(removedTrips)/float64(tripsBefore))
}

func runKeepOnly(feed *gtfsparser.Feed, namesToKeep []string) (int, int) {
	fmt.Fprintf(os.Stdout, "Keeping only selected agencies... ")
	removedAgencies := []string{}

	for agencyId, agency := range feed.Agencies {
		if !slices.Contains(namesToKeep, agency.Name) {
			delete(feed.Agencies, agencyId)
			removedAgencies = append(removedAgencies, agency.Id)
		}
	}

	removedTrips := removeAgencyTrips(feed, removedAgencies)
	return len(removedAgencies), removedTrips
}

func runRemove(feed *gtfsparser.Feed, namesToRemove []string) (int, int) {
	fmt.Fprintf(os.Stdout, "Removing selected agencies... ")
	removedAgencies := []string{}

	for agencyId, agency := range feed.Agencies {
		if slices.Contains(namesToRemove, agency.Name) {
			delete(feed.Agencies, agencyId)
			removedAgencies = append(removedAgencies, agency.Id)
		}
	}

	removedTrips := removeAgencyTrips(feed, removedAgencies)
	return len(removedAgencies), removedTrips
}

func removeAgencyTrips(feed *gtfsparser.Feed, removedAgencies []string) int {
	removedTrips := 0
	for tripId := range feed.Trips {
		if slices.Contains(removedAgencies, feed.Trips[tripId].Route.Agency.Id) {
			feed.DeleteTrip(tripId)
			removedTrips += 1
		}
	}
	return removedTrips
}
