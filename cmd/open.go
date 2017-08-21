package cmd

import (
	"errors"
	"fmt"

	"github.com/exercism/cli/browser"
	"github.com/exercism/cli/comms"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
)

// openCmd opens the designated exercise in the browser.
var openCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"o"},
	Short:   "Open an exercise on the website.",
	Long: `Open the specified exercise to the solution page on the Exercism website.

Pass either the name of an exercise, or the path to the directory that contains
the solution you want to see on the website.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewUserConfig()
		BailOnError(err)
		ws := workspace.New(cfg.Workspace)

		if len(args) != 1 {
			// TODO: usage
			return
		}

		paths, err := ws.Locate(args[0])
		BailOnError(err)

		solutions, err := workspace.NewSolutions(paths)
		BailOnError(err)

		if len(solutions) == 0 {
			return
		}

		if len(solutions) > 1 {
			var mine []*workspace.Solution
			for _, s := range solutions {
				if s.IsRequester {
					mine = append(mine, s)
				}
			}
			solutions = mine
		}

		selection := comms.NewSelection()
		for _, solution := range solutions {
			selection.Items = append(selection.Items, solution)
		}
		for {
			prompt := `
We found more than one. Which one did you mean?
Type the number of the one you want to select.

%s
> `
			option, err := selection.Pick(prompt)
			if err != nil {
				fmt.Println(err)
				continue
			}
			solution, ok := option.(*workspace.Solution)
			if ok {
				browser.Open(solution.URL)
				break
			}
			BailOnError(errors.New("should never happen"))
		}

	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}
