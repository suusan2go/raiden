// Copyright © 2017 suusan2go <ksuzuki180@gmail.com>
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
	"github.com/suusan2go/raiden/github"
)

type releasesClean struct {
	repository string
	owner      string
	prefix     string
	dry        bool
	year       int
	months     int
	days       int
}

// releasesCmd represents the releases command
func releasesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "releases",
		Short: "raiden releases command",
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
raiden releases clean -r many_releases_tag_repo -o user_or_org_name --months 1`,
		Run: c.clean,
	}

	flags := cmd.Flags()
	flags.StringVarP(&c.prefix, "prefix", "p", "", "set Prefix of release tag name")
	flags.StringVarP(&c.repository, "repository", "r", "", "Set repository name like hoge")
	flags.StringVarP(&c.owner, "owner", "o", "", "Set owner name of repository like suusan2go")
	flags.BoolVarP(&c.dry, "dry", "d", false, "Just get reases tag and not delete")
	flags.IntVar(&c.year, "year", 0, "clean releases year before")
	flags.IntVar(&c.months, "months", 0, "clean releases Month before")
	flags.IntVar(&c.days, "days", 0, "clean releases year before")

	return cmd
}

func init() {
	rc := releasesCmd()
	RootCmd.AddCommand(rc)

	rc.AddCommand(reasesCleanCmd())
}

func (c *releasesClean) clean(cmd *cobra.Command, args []string) {
	log.Printf("start clean releases tags for %s/%s", c.owner, c.repository)
	g := github.Initialize(c.owner, c.repository)
	g.DeleteReleases(c.dry, c.year, c.months, c.days, c.prefix)
	log.Println("clean releases tags for " + c.owner + "/" + c.repository)
}
