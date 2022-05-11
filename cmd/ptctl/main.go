package main

import (
	"fmt"
	"os"

	"github.com/osm/pingtraq"
	"github.com/spf13/cobra"
)

func main() {
	var database string

	var addCmd = &cobra.Command{
		Use:   "add <name>",
		Short: "Add new ping endpoint",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := pingtraq.Init(database); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return
			}

			name := args[0]
			if id, _ := pingtraq.IsPing(name); id != "" {
				fmt.Fprintf(os.Stderr, "%v does already exist\n", name)
				return
			}

			pingtraq.AddPing(name)
		},
	}
	addCmd.Flags().StringVarP(&database, "database", "d", "", "Database file (required)")
	addCmd.MarkFlagRequired("database")

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all ping endpoints",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := pingtraq.Init(database); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return
			}

			names, err := pingtraq.ListPing()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return
			}

			for _, n := range names {
				fmt.Printf("%s\n", n)
			}
		},
	}
	listCmd.Flags().StringVarP(&database, "database", "d", "", "Database file (required)")
	listCmd.MarkFlagRequired("database")

	var getCmd = &cobra.Command{
		Use:   "get <name>",
		Short: "Get all records for a ping endpoint",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := pingtraq.Init(database); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return
			}

			name := args[0]
			if id, _ := pingtraq.IsPing(name); id == "" {
				fmt.Fprintf(os.Stderr, "%s does not exist\n", name)
				return
			}

			records, err := pingtraq.ListPingRecords(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return
			}

			for _, r := range records {
				fmt.Printf("%s,%s,%s,%s\n", r.Client, r.BatteryLevel, r.Address, r.CreatedAt)
			}
		},
	}
	getCmd.Flags().StringVarP(&database, "database", "d", "", "Database file (required)")
	getCmd.MarkFlagRequired("database")

	var rootCmd = &cobra.Command{Use: "ptctl"}
	rootCmd.AddCommand(addCmd, listCmd, getCmd)
	rootCmd.Execute()
}
