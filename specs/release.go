package specs

import (
	"github.com/2zqa/ssot-specs-collector/aggregator"
	"github.com/dekobon/distro-detect/linux"
)

type Release struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Codename string `json:"codename"`
}

func FetchRelease() (Release, error) {
	a := aggregator.New("release")
	var release Release
	ld := linux.DiscoverDistro()

	if name := ld.LsbRelease["DISTRIB_ID"]; a.CheckPredicateAndAdd(name != "", "could not retrieve release name") {
		release.Name = name
	}

	if version := ld.LsbRelease["DISTRIB_RELEASE"]; a.CheckPredicateAndAdd(version != "", "could not retrieve release version") {
		release.Version = version
	}

	if codename := ld.LsbRelease["DISTRIB_CODENAME"]; a.CheckPredicateAndAdd(codename != "", "could not retrieve release codename") {
		release.Codename = codename
	}

	return release, a.ErrorOrNil()
}
