package main

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

func error01() error {
	return errors.Wrap(sql.ErrNoRows, "error01 failed")
}

func error02() error {
	return errors.WithMessage(error01(), "error02 failed")
}

func main() {
	err := error02()
	if errors.Cause(err) == sql.ErrNoRows {
		fmt.Printf("data not found, %v\n", err)
		fmt.Printf("%+v\n", err)
		return
	}
	if err != nil {
		// unknown error
	}
}
