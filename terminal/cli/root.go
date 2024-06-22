package cli

import (
	"fmt"
	"time"

	"github.com/FoxHoundTechnology/remote-control-miners/terminal/dashboard"
	"github.com/FoxHoundTechnology/remote-control-miners/terminal/query"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "termui",
		Short: "terminal ui miner dashboard",
	}
	rootCmd.AddCommand(
		StartClientCommand(),
	)
	return rootCmd
}

func StartClientCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start termui client",
		Long:  "start termui client with specified address:port",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryClient, err := query.NewQueryClient(args[0])
			if err != nil {
				return fmt.Errorf("invalid url format: %v", err)
			}

			// NOTE: implement .env to detect dev environment
			//queryClient := query.MockQueryClient{}
			dashboardModel, err := dashboard.InitModel(queryClient, time.Second)
			if err != nil {
				return err
			}

			dashboard := tea.NewProgram(dashboardModel)
			if _, err := dashboard.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
