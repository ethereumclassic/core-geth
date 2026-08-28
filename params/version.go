// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"fmt"
)

const (
	VersionMajor = 1          // Major version component of the current release
	VersionMinor = 13         // Minor version component of the current release
	VersionPatch = 0          // Patch version component of the current release
	VersionMeta  = "unstable" // Version metadata to append to the version string
	// VersionName is the devp2p FAMILY name and must not carry the version.
	// node.Config.NodeName() appends "/v<VersionWithMeta>" itself, and network
	// censuses group peers by this first segment -- etcnodes.org counts 512 of
	// 536 ETC peers under "CoreGeth". A version-bearing family name changes
	// every release, so those nodes would be counted as a different client and
	// the version would appear twice in one identity string.
	VersionName = "CoreGeth"
)

// VersionMeta carries the release stage and must never be empty. It advances
// "unstable" -> "RC1", "RC2", ... -> "stable", set at the moment each tag is
// cut. Both ArchiveVersion and VersionWithCommit below branch on
// VersionMeta != "stable" rather than on emptiness, so an empty value takes
// the not-a-release path and appends nothing after its separator, producing
// "1.13.0--<commit>". Emptying this field to suppress a suffix therefore
// yields a malformed archive name rather than a clean one.

// Version holds the textual version string.
var Version = func() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}()

// VersionWithMeta holds the textual version string including the metadata.
var VersionWithMeta = func() string {
	v := Version
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

// ArchiveVersion holds the textual version string used for Geth archives. e.g.
// "1.8.11-dea1ce05" for stable releases, or "1.8.13-unstable-21c059b6" for unstable
// releases.
func ArchiveVersion(gitCommit string) string {
	vsn := Version
	if VersionMeta != "stable" {
		vsn += "-" + VersionMeta
	}
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	return vsn
}

func VersionWithCommit(gitCommit, gitDate string) string {
	vsn := VersionWithMeta
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	if (VersionMeta != "stable") && (gitDate != "") {
		vsn += "-" + gitDate
	}
	return vsn
}
