package main

import (
	"log"
	"net/http"

	"github.com/osm/pingtraq"
	"github.com/spf13/cobra"
)

func main() {
	var port string
	var database string
	var rootCmd = &cobra.Command{
		Use: "ptd",
		Run: func(cmd *cobra.Command, args []string) {
			if err := pingtraq.Init(database); err != nil {
				log.Fatal(err)
			}

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)

				url := r.URL.Path
				if len(url) < 2 {
					return
				}

				name := url[1:]
				var id string
				if id = pingtraq.IsPing(name); id == "" {
					return
				}

				if err := pingtraq.AddPingRecord(id, r); err != nil {
					log.Printf("AddPingRecord: %v\n", err)
					return
				}

				log.Printf("Ping %s from %s\n", name, r.RemoteAddr)

				return
			})

			log.Printf("Listening on *:%s\n", port)
			log.Fatal(http.ListenAndServe(":"+port, nil))
		},
	}
	rootCmd.Flags().StringVarP(&port, "port", "p", "", "Listen port (required)")
	rootCmd.Flags().StringVarP(&database, "database", "d", "", "Database file (required)")
	rootCmd.MarkFlagRequired("port")
	rootCmd.MarkFlagRequired("database")
	rootCmd.Execute()
}
