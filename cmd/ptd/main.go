package main

import (
	"log"
	"net/http"
	"os/exec"

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
				var afterHook *string
				if id, afterHook = pingtraq.IsPing(name); id == "" {
					return
				}

				if err := pingtraq.AddPingRecord(id, r); err != nil {
					log.Printf("AddPingRecord: %v\n", err)
					return
				}

				if afterHook != nil {
					cmd := exec.Command(*afterHook)
					if out, err := cmd.CombinedOutput(); err != nil {
						log.Printf("After hook %s error: %v\n", *afterHook, err)
					} else {
						o := string(out)
						if o[len(o)-1] == '\n' {
							o = o[:len(o)-1]
						}
						log.Printf("After hook %s: %s\n", *afterHook, o)
					}
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
