package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/narqo/go-badge"
	"github.com/xanzy/go-gitlab"
)

type GitLabMergeRequests struct{}

func (b *GitLabMergeRequests) Render(column string, project Project) string {
	if !isGitLab(project) {
		return ""
	}
	gl := gitlab.NewClient(nil, GITLAB_ACCESS_TOKEN)

	state := "opened"
	options := &gitlab.ListProjectMergeRequestsOptions{
		State: &state,
	}

	id := project.Namespace + "/" + project.Name
	id = strings.Trim(id, "/")
	_, response, err := gl.MergeRequests.ListProjectMergeRequests(id, options)
	if err != nil {
		return svgBadge("merge requests", err.Error(), badge.ColorLightgrey, project.URL)
	}

	color := badge.ColorBrightgreen
	if response.TotalItems > 0 {
		color = badge.ColorYellow
	}
	return svgBadge("merge requests", fmt.Sprintf("%d open", response.TotalItems), color, project.URL)
}

type GitLabBranches struct{}

func (b *GitLabBranches) Render(column string, project Project) string {
	if !isGitLab(project) {
		return ""
	}
	gl := gitlab.NewClient(nil, GITLAB_ACCESS_TOKEN)

	options := &gitlab.ListBranchesOptions{}

	id := project.Namespace + "/" + project.Name
	id = strings.Trim(id, "/")
	_, response, err := gl.Branches.ListBranches(id, options)
	if err != nil {
		return svgBadge("branches", err.Error(), badge.ColorLightgrey, project.URL)
	}

	color := badge.ColorBrightgreen
	if response.TotalItems > 1 {
		color = badge.ColorGreen
	}
	if response.TotalItems > 2 {
		color = badge.ColorYellow
	}
	return svgBadge("branches", fmt.Sprintf("%d open", response.TotalItems), color, project.URL)
}

type GitLabTag struct{}

func (b *GitLabTag) Render(column string, project Project) string {
	if !isGitLab(project) {
		return ""
	}
	gl := gitlab.NewClient(nil, GITLAB_ACCESS_TOKEN)

	options := &gitlab.ListTagsOptions{}

	id := project.Namespace + "/" + project.Name
	id = strings.Trim(id, "/")
	tags, _, err := gl.Tags.ListTags(id, options)
	if err != nil {
		return svgBadge("tag", err.Error(), badge.ColorLightgrey, project.URL)
	}
	tag := tags[0]
	return svgBadge("tag", tag.Name, badge.ColorBlue, project.URL)
}

type GitLabProject struct {
	field string
}

func (b *GitLabProject) Render(column string, project Project) string {
	if !isGitLab(project) {
		return ""
	}
	gl := gitlab.NewClient(nil, GITLAB_ACCESS_TOKEN)

	t := true
	options := &gitlab.GetProjectOptions{
		Statistics: &t,
	}

	id := project.Namespace + "/" + project.Name
	id = strings.Trim(id, "/")

	gitlabProject, _, err := gl.Projects.GetProject(id, options)
	if err != nil {
		return svgBadge("gitlab", err.Error(), badge.ColorLightgrey, project.URL)
	}

	switch b.field {
	case "issues":
		color := badge.ColorBrightgreen
		if gitlabProject.OpenIssuesCount > 0 {
			color = badge.ColorYellow
		}
		return svgBadge("issues", fmt.Sprintf("%d open", gitlabProject.OpenIssuesCount), color, project.URL)
	case "lastcommit":
		return svgBadge("last commit", humanize.Time(*gitlabProject.LastActivityAt), badge.ColorLightgrey, project.URL)
	case "stars":
		return svgBadge("stars", fmt.Sprint(gitlabProject.StarCount), badge.ColorBlue, project.URL)
	case "forks":
		return svgBadge("Fork", fmt.Sprint(gitlabProject.ForksCount), badge.ColorBlue, project.URL)
	case "size":
		return svgBadge("repo size", humanize.Bytes(uint64(gitlabProject.Statistics.RepositorySize)), badge.ColorBlue, project.URL)
	}

	return svgBadge("error", "unknown field", badge.ColorLightgrey, project.URL)
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
