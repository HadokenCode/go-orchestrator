package types

// release.json file in each container
type Release struct {
	ID                         string        `json:"id"`
	ReleaseStep                ReleaseStep   `json:"step"`
	JiraID                     string        `json:"jira.issue_id"`
	NexusStagingRepositoryId   string        `json:"nexus.staged_repository_id"`
	NexusStagingRepositoryUrl  string        `json:"nexus.staged_repository_url"`
}
type ReleaseStep struct {
	Name          string        `json:"name"`
	Status        string        `json:"status"`
}