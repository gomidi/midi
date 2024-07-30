package version

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Version represents a version (e.g. a release version)
type Version struct {
	Major uint16
	Minor uint16
	Patch uint16
}

// String returns the string representation of the version (without a preceeding v
func (v Version) String() string {
	return fmt.Sprintf("%v.%v.%v", v.Major, v.Minor, v.Patch)
}

// Equals is true when the version exactly equals the given version
func (v Version) Equals(o Version) bool {
	return v.Major == o.Major && v.Minor == o.Minor && v.Patch == o.Patch
}

// EqualsMajor is true when the version equals the given version on the Major
func (v Version) EqualsMajor(o Version) bool {
	return v.Major == o.Major
}

// EqualsMinor is true when the version equals the given version on the Major and Minor
func (v Version) EqualsMinor(o Version) bool {
	return v.Major == o.Major && v.Minor == o.Minor
}

// Less returns true, if the Version is smaller than the given one.
func (v Version) Less(o Version) bool {
	if v.Major != o.Major {
		return v.Major < o.Major
	}

	if v.Minor != o.Minor {
		return v.Minor < o.Minor
	}

	return v.Patch < o.Patch
}

// Parse parses the version out of the given string. Valid strings are "v0.0.1" or "1.0" or "12" etc.
func Parse(v string) (*Version, error) {
	v = strings.TrimLeft(v, "v")
	nums := strings.Split(v, ".")
	var ns = make([]uint16, len(nums))

	var ver Version

	for i, n := range nums {
		nn, err := strconv.Atoi(n)
		if err != nil {
			return nil, err
		}
		if nn < 0 {
			return nil, fmt.Errorf("number must not be < 0")
		}
		ns[i] = uint16(nn)
	}

	switch len(ns) {
	case 1:
		if ns[0] == 0 {
			return nil, fmt.Errorf("invalid version %q", v)
		}
		ver.Major = ns[0]
	case 2:
		if (ns[0] + ns[1]) == 0 {
			return nil, fmt.Errorf("invalid version %q", v)
		}
		ver.Major = ns[0]
		ver.Minor = ns[1]
	case 3:
		if (ns[0] + ns[1] + ns[2]) == 0 {
			return nil, fmt.Errorf("invalid version %q", v)
		}
		ver.Major = ns[0]
		ver.Minor = ns[1]
		ver.Patch = ns[2]
	default:
		return nil, fmt.Errorf("invalid version string %q", v)
	}

	return &ver, nil

}

// Versions is a sortable slice of *Version
type Versions []*Version

// Less returns true, if the version of index a is less than the version of index b
func (v Versions) Less(a, b int) bool {
	return v[a].Less(*v[b])
}

// Len returns the number of *Version inside the slice
func (v Versions) Len() int {
	return len(v)
}

// Swap swaps the *Version of the index a with that of the index b
func (v Versions) Swap(a, b int) {
	v[a], v[b] = v[b], v[a]
}

// Sort sorts the slice and returns it
func (v Versions) Sort() Versions {
	sort.Sort(v)
	return v
}

// Last returns the last *Version of the slice.
func (v Versions) Last() *Version {
	if len(v) == 0 {
		return nil
	}

	return v[len(v)-1]
}

// First returns the first *Version of the slice.
func (v Versions) First() *Version {
	if len(v) == 0 {
		return nil
	}

	return v[0]
}
