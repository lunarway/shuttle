package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lunarway/shuttle/pkg/errors"
	"github.com/lunarway/shuttle/pkg/git"
)

// DocumentationURL returns a URL pointing to plan documentation if any is
// available. Plan reference and plan documentation field is inspected and
// parsed.
func (p *ShuttleProjectContext) DocumentationURL() (string, error) {
	var ref string
	switch {
	case p.Plan.Documentation != "":
		ref = p.Plan.Documentation
	case p.Config.Plan != "":
		ref = p.Config.Plan
	default:
		return "", errors.NewExitCode(1, "Could not find any plan documentation")
	}

	switch {
	case git.IsPlan(ref):
		return normalizeGitPlan(git.ParsePlan(ref))
	case isHTTPSPlan(ref):
		return ref, nil
	case filepath.IsAbs(ref), strings.HasPrefix(ref, "./"), strings.HasPrefix(ref, "../"):
		return "", errors.NewExitCode(2, "Local plan has no documentation")
	default:
		return "", errors.NewExitCode(1, "Could not detect protocol for plan '%s'", ref)
	}
}

func normalizeGitPlan(p git.Plan) (string, error) {
	switch p.Protocol {
	case "https":
		return fmt.Sprintf("%s://%s", p.Protocol, p.Repository), nil
	case "ssh":
		repoSlug := strings.TrimPrefix(p.Repository, fmt.Sprintf("%s:", p.Host))
		return fmt.Sprintf("https://%s/%s", p.Host, repoSlug), nil
	default:
		// this should never happen as parsed git plans always has a protocol of ssh
		// or https
		return "", errors.NewExitCode(1, "Could not parse git plan reference")
	}
}
