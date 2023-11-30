package main

import (
	"github.com/spf13/cobra"
)

func main() {
	v := &cobra.Command{}

	_ = v.Execute()
}
