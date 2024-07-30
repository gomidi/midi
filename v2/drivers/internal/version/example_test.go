package version

import (
	"fmt"
)

func Example() {

	var v1 = "v1.1.0"
	var v2 = "1.0.4"
	var invalid = ".1.0"

	vers1, err := Parse(v1)

	if err != nil {
		fmt.Printf("could not parse version %q", v1)
	}

	vers2, err := Parse(v2)

	if err != nil {
		fmt.Printf("could not parse version %q", v1)
	}

	if vers2.Less(*vers1) {
		fmt.Printf("version %s < version %s\n", vers2, vers1)
	}

	vers2.Minor = 1

	if vers1.Less(*vers2) {
		fmt.Printf("version %s < version %s\n", vers1, vers2)
	}

	last := Versions{vers1, vers2}.Sort().Last()

	fmt.Printf("last version is %s\n", last)

	_, err = Parse(invalid)

	if err != nil {
		fmt.Printf("could not parse version %q", invalid)
	}

	// Output:
	// version 1.0.4 < version 1.1.0
	// version 1.1.0 < version 1.1.4
	// last version is 1.1.4
	// could not parse version ".1.0"
}
