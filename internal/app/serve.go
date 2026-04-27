package app

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func newServeCmd() *cobra.Command {
	var dir string
	var port int
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve generated docs locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := fmt.Sprintf(":%d", port)
			fmt.Printf("Serving %s on http://localhost%s\n", dir, addr)
			return http.ListenAndServe(addr, http.FileServer(http.Dir(dir)))
		},
	}
	cmd.Flags().StringVar(&dir, "dir", "./docs", "Directory to serve")
	cmd.Flags().IntVar(&port, "port", 8899, "Port")
	return cmd
}
