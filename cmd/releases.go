// Copyright © 2017 suzan2go <ksuzuki180@gmail.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

type releasesClean struct {
	repository string
	dry        bool
}

// releasesCmd represents the releases command
func releasesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "releases",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("releases called")
		},
	}

	return cmd
}

func reasesCleanCmd() *cobra.Command {
	c := &releasesClean{}
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "clean reases tag",
		Long: `clean up reasess tags:
raiden clean -r "suzan2go/many_releases_tag_repo"`,
		Run: c.clean,
	}

	flags := cmd.Flags()
	flags.StringVarP(&c.repository, "repository", "r", "", "Set repository name like suzan2go/hoge")
	flags.BoolVarP(&c.dry, "dry", "d", false, "Just get reases tag and not delete")

	return cmd
}

func init() {
	rc := releasesCmd()
	RootCmd.AddCommand(rc)

	rc.AddCommand(reasesCleanCmd())
}

func (c *releasesClean) clean(cmd *cobra.Command, args []string) {
	if len(c.repository) == 0 {
		log.Fatal("repository not specified")
	}
	if c.dry {
		fmt.Println("[dry run] clean releases tags for " + c.repository)
		return
	}
	fmt.Println("clean releases tags for " + c.repository)
}
