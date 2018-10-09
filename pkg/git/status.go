package git

import (
	"fmt"
	"strings"

	go_cmd "github.com/go-cmd/cmd"
)

func gitCmd2(command string, dir string) go_cmd.Status {
	cmdOptions := go_cmd.Options{
		Buffered: true,
	}
	execCmd := go_cmd.NewCmdOptions(cmdOptions, "sh", "-c", "cd '"+dir+"'; git "+command)
	status := <-execCmd.Start()
	return status
}

type StatusType string

const (
	StatusTypeUnchanged = "Unchanged"
	StatusTypeModified  = "Modified"
	StatusTypeAdded     = "Added"
	StatusTypeDeleted   = "Deleted"
	StatusTypeRenamed   = "Renamed"
	StatusTypeCopied    = "Copied"
	StatusTypeUntracked = "Untracked"
	StatusTypeIgnored   = "Ignored"
)

type Status struct {
	changes    bool
	files      []FileStatus
	mergeState bool
}

type FileStatus struct {
	FilePath         string
	OriginalFilePath string
	WorkTreeStatus   StatusType
	IndexStatus      StatusType
}

func getStatus(dir string) Status {
	// 1 M. N... 100755 100755 100755 df43166016aa8f73ff348175fb11c8f061ebd871 79e6ffe9583d538fdf957b5e7595ba9aa9ed9929 scripts/env.sh
	cmdStatus := gitCmd2("status --porcelain=v2 --branch", dir)

	status := Status{}

	for _, line := range cmdStatus.Stdout {
		var x, y, sub, path, origPath string // mH, mI, mW, hH, hI, score
		var merge bool = false
		switch line[0] {
		case '#':
			// header
			continue
		case '1':
			// Ordinary changed entries
			// 1 <XY> <sub> <mH> <mI> <mW> <hH> <hI> <path>
			parts := strings.SplitN(line, " ", 9)
			x = string(parts[1][0])
			y = string(parts[1][1])
			sub = parts[2]
			path = parts[8]
		case '2':
			// Renamed or copied
			// 2 <XY> <sub> <mH> <mI> <mW> <hH> <hI> <X><score> <path><sep><origPath>
			parts := strings.SplitN(line, " ", 9)
			x = string(parts[1][0])
			y = string(parts[1][1])
			sub = parts[2]
			pathParts := strings.SplitN(parts[8], "\t", 2)
			path = pathParts[0]
			origPath = pathParts[1]
		case '!':
			// ignored file
			// ! <path>
			x = "."
			y = "!"
			parts := strings.SplitN(line, " ", 2)
			path = parts[1]
		case '?':
			// untracked file
			// ? <path>
			x = "."
			y = "?"
			parts := strings.SplitN(line, " ", 2)
			path = parts[1]
		case 'u':
			// Unmerged entries
			// u <xy> <sub> <m1> <m2> <m3> <mW> <h1> <h2> <h3> <path>
			merge = true
			parts := strings.SplitN(line, " ", 11)
			x = string(parts[1][0])
			y = string(parts[1][1])
			sub = parts[2]
			path = parts[11]
		default:
			panic(fmt.Sprintf("Unknown git porcelain type '%s' in '%s'", string(line[0]), line))
		}

		if merge {
			status.mergeState = true
			continue // Skip merge changes
		}

		if sub != "N..." {
			continue // Skip submodule changes
		}

		fileStatus := FileStatus{
			FilePath:         path,
			OriginalFilePath: origPath,
			IndexStatus:      statusTypeMapping[x],
			WorkTreeStatus:   statusTypeMapping[y],
		}
		status.files = append(status.files, fileStatus)

		if isChange(fileStatus.IndexStatus) || isChange(fileStatus.WorkTreeStatus) {
			status.changes = true
		}
	}

	return status
}

var statusTypeMapping = map[string]StatusType{
	".": StatusTypeUnchanged,
	"!": StatusTypeIgnored,
	"?": StatusTypeUntracked,
	"M": StatusTypeModified,
	"A": StatusTypeAdded,
	"R": StatusTypeRenamed,
	"C": StatusTypeCopied,
	"D": StatusTypeDeleted,
}

func isChange(statusType StatusType) bool {
	switch statusType {
	case StatusTypeUnchanged:
		return false
	case StatusTypeIgnored:
		return false
	case StatusTypeUntracked:
		return true
	case StatusTypeModified:
		return true
	case StatusTypeAdded:
		return true
	case StatusTypeRenamed:
		return true
	case StatusTypeCopied:
		return true
	case StatusTypeDeleted:
		return true
	default:
		panic(fmt.Sprintf("Unhandled status type '%s'", statusType))
	}
}

// STATUS CODES:
// X          Y     Meaning
// -------------------------------------------------
//          [AMD]   not updated
// M        [ MD]   updated in index
// A        [ MD]   added to index
// D                deleted from index
// R        [ MD]   renamed in index
// C        [ MD]   copied in index
// [MARC]           index and work tree matches
// [ MARC]     M    work tree changed since index
// [ MARC]     D    deleted in work tree
// [ D]        R    renamed in work tree
// [ D]        C    copied in work tree
// -------------------------------------------------
// D           D    unmerged, both deleted
// A           U    unmerged, added by us
// U           D    unmerged, deleted by them
// U           A    unmerged, added by them
// D           U    unmerged, deleted by us
// A           A    unmerged, both added
// U           U    unmerged, both modified
// -------------------------------------------------
// ?           ?    untracked
// !           !    ignored
// -------------------------------------------------

// HEADERS:
// Line                                     Notes
// ------------------------------------------------------------
// # branch.oid <commit> | (initial)        Current commit.
// # branch.head <branch> | (detached)      Current branch.
// # branch.upstream <upstream_branch>      If upstream is set.
// # branch.ab +<ahead> -<behind>           If upstream is set and
// 	     the commit is present.
// ------------------------------------------------------------

// FIELD EXPLANATION
// Field       Meaning
// --------------------------------------------------------
// <XY>        A 2 character field containing the staged and
// unstaged XY values described in the short format,
// with unchanged indicated by a "." rather than
// a space.
// <sub>       A 4 character field describing the submodule state.
// "N..." when the entry is not a submodule.
// "S<c><m><u>" when the entry is a submodule.
// <c> is "C" if the commit changed; otherwise ".".
// <m> is "M" if it has tracked changes; otherwise ".".
// <u> is "U" if there are untracked changes; otherwise ".".
// <mH>        The octal file mode in HEAD.
// <mI>        The octal file mode in the index.
// <mW>        The octal file mode in the worktree.
// <hH>        The object name in HEAD.
// <hI>        The object name in the index.
// <X><score>  The rename or copy score (denoting the percentage
// of similarity between the source and target of the
// move or copy). For example "R100" or "C75".
// <path>      The pathname.  In a renamed/copied entry, this
// is the target path.
// <sep>       When the `-z` option is used, the 2 pathnames are separated
// with a NUL (ASCII 0x00) byte; otherwise, a tab (ASCII 0x09)
// byte separates them.
// <origPath>  The pathname in the commit at HEAD or in the index.
// This is only present in a renamed/copied entry, and
// tells where the renamed/copied contents came from.
// --------------------------------------------------------
