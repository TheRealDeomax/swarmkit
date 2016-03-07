package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/swarm-v2/agent"
	"github.com/docker/swarm-v2/identity"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	agentCmd = &cobra.Command{
		Use:   "agent",
		Short: "Run the swarm agent",
		Long: `Start a swarm agent with the provided path. If starting from an
empty path, the agent will allocate an identity and startup. If data is
already present, the agent will recover and startup.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			if name == "" {
				hn, err := os.Hostname()
				if err != nil {
					return err
				}

				log.Debugf("agent: fallback to hostname as name")
				name = hn
			}

			id, err := cmd.Flags().GetString("id")
			if err != nil {
				return err
			}

			if id == "" {
				log.Debugf("agent: generated random identifier")
				id = identity.NewID()
			}

			managerAddrs, err := cmd.Flags().GetStringSlice("manager")
			if err != nil {
				return err
			}

			log.Debugf("managers:", managerAddrs)
			managers := agent.NewManagers(managerAddrs...)

			ag, err := agent.New(&agent.Config{
				ID:       id,
				Name:     name,
				Managers: managers,
			})
			if err != nil {
				log.Fatalln(err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := ag.Run(ctx); err != nil {
				log.Fatalln(err)
			}

			return nil
		},
	}
)

func init() {
	agentCmd.Flags().String("id", "", "Specifies the identity of the node")
	agentCmd.Flags().String("name", "", "Agent name. Defaults to hostname")
	agentCmd.Flags().StringSliceP("manager", "m", []string{"localhost:4242"}, "Specify one or more manager addresses")
}
