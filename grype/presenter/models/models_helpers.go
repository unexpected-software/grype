package models

import (
	"testing"

	"github.com/anchore/grype/grype"

	"github.com/anchore/syft/syft/distro"

	"github.com/anchore/stereoscope/pkg/image"

	syftPkg "github.com/anchore/syft/syft/pkg"
	syftSource "github.com/anchore/syft/syft/source"

	"github.com/anchore/grype/grype/vulnerability"

	"github.com/anchore/grype/grype/match"
	"github.com/anchore/grype/grype/pkg"
)

func GenerateAnalysis(t *testing.T) grype.Analysis {
	t.Helper()

	packages := generatePackages(t)
	matches := generateMatches(t, packages[0])
	context := generateContext(t)

	return grype.Analysis{
		Matches:          matches,
		Packages:         packages,
		Context:          context,
		MetadataProvider: NewMetadataMock(),
	}
}

func generateMatches(t *testing.T, p pkg.Package) match.Matches {
	t.Helper()

	matches := []match.Match{
		{
			Type: match.ExactDirectMatch,
			Vulnerability: vulnerability.Vulnerability{
				ID:             "CVE-1999-0001",
				RecordSource:   "source-1",
				FixedInVersion: "the-next-version",
			},
			Package: p,
			Matcher: match.DpkgMatcher,
			SearchKey: map[string]interface{}{
				"distro": map[string]string{
					"type":    "ubuntu",
					"version": "20.04",
				},
			},
			SearchMatches: map[string]interface{}{
				"constraint": ">= 20",
			},
		},
		{
			Type: match.ExactIndirectMatch,
			Vulnerability: vulnerability.Vulnerability{
				ID:           "CVE-1999-0002",
				RecordSource: "source-2",
			},
			Package: p,
			Matcher: match.DpkgMatcher,
			SearchKey: map[string]interface{}{
				"cpe": "somecpe",
			},
			SearchMatches: map[string]interface{}{
				"constraint": "somecpe",
			},
		},
	}

	collection := match.NewMatches()
	collection.Add(p, matches...)

	return collection
}

func generatePackages(t *testing.T) []pkg.Package {
	t.Helper()

	return []pkg.Package{
		{
			Name:    "package-1",
			Version: "1.1.1",
			Type:    syftPkg.DebPkg,
			Locations: []syftSource.Location{
				{RealPath: "/somefile-1.txt"},
			},
			CPEs: []syftPkg.CPE{
				{
					Part:     "a",
					Vendor:   "anchore",
					Product:  "engine",
					Version:  "0.9.2",
					Language: "python",
				},
			},
		},
		{
			Name:    "package-2",
			Version: "2.2.2",
			Type:    syftPkg.DebPkg,
			Locations: []syftSource.Location{
				{RealPath: "/somefile-2.txt"},
			},
		},
	}
}

func generateContext(t *testing.T) pkg.Context {
	img := image.Image{}

	src, err := syftSource.NewFromImage(&img, "user-input")
	if err != nil {
		t.Fatalf("failed to create scope: %+v", err)
	}

	d, err := distro.NewDistro(distro.CentOS, "8.0", "rhel")
	if err != nil {
		t.Fatalf("could not make distro: %+v", err)
	}

	return pkg.Context{
		Source: &src.Metadata,
		Distro: &d,
	}
}
