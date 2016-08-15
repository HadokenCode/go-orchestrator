package types

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/colorstring"

	ermcmd "go-orchestrator/cmd"
	"strings"
)

// Catalog describes the data for the project to release.
type Catalog []Project

// Project to release
type Project struct {
	Name            string        `json:"name"`
	GitOrganization string        `json:"git_organization"`
	Labels          string        `json:"labels"`
	MavenProperty   string        `json:"maven_property_version"`
	ReleaseParams   ReleaseParams `json:"release"`
	Container       Container     `json:"container"`
}

// Container to use to release the project
type Container struct {
	Image string `json:"image"`
	Step  int    `json:"step"`
}

// ReleaseParams specific to the project for the release
type ReleaseParams struct {
	Branch                 string `json:"branch"`
	Version                string `json:"version"`
	CurrentSnapshotVersion string `json:"current_snapshot_version"`
	NextSnapshotVersion    string `json:"next_snapshot_version"`
	Patches                string `json:"patches"`
	PatchesAfterRelease    string `json:"patches_after_release"`
	NexusHost              string `json:"nexus_host"`
	NexusStagingProfile    string `json:"nexus_staging_profile"`
}

// DisplayCatalog show all prokjects available into the catalog
// filtered by label
func DisplayCatalog(catalog Catalog, label string) {
	var projects = Catalog{}

	if label == ermcmd.LabelAll {
		projects = catalog
	} else {
		projects = catalog.FilterByLabel(label)
	}
    fmt.Printf(colorstring.Color("[blue] \nProjects lits: \n"))
	for _, p := range projects {
		fmt.Printf(colorstring.Color("[blue] * %s - %s - %s \n"), p.Name, p.ReleaseParams.Version, p.ReleaseParams.Branch)
		fmt.Printf("JSON %+v", p)
		fmt.Println("\n\n-----")
	}

}

// GetCatalog Get performs a request to get the catalog data for a release.
func GetCatalog(uri string) (catalog Catalog, err error) {

	// send the request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return catalog, nil
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return catalog, fmt.Errorf("Http request to %s failed: %s", req.URL, err.Error())
	}
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&catalog); err != nil {
		return catalog, fmt.Errorf("Decoding catalog response failed: %v", err)
	}

	return catalog, nil
}

func (slice Catalog) Len() int {
	return len(slice)
}

func (slice Catalog) Less(i, j int) bool {
	return slice[i].Container.Step < slice[j].Container.Step
}

func (slice Catalog) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// FilterByLabel returns a list of projects that contains the given label. An
// empty name returns all projects.
func (projects Catalog) FilterByLabel(label string) Catalog {
	filtered := make(Catalog, 0)

	for _, p := range projects {
		if strings.Contains(p.Labels, label) {
			filtered = append(filtered, p)
		}

	}

	return filtered
}
