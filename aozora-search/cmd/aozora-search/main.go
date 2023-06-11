package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aozora-search",
	Short: "Aozora Search is a CLI for searching Aozora Bunko",
	Long:  `Aozora Search is a CLI for searching Aozora Bunko. It provides several sub-commands for different types of searches.`,
}

func main() {
	var dsn string
	rootCmd.PersistentFlags().StringVarP(&dsn, "database", "d", "database.sqlite", "database")
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authorsCmd := &cobra.Command{
		Use:   "authors",
		Short: "Show authors",
		Run: func(cmd *cobra.Command, args []string) {
			err := showAuthors(db)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	titlesCmd := &cobra.Command{
		Use:   "titles [AuthorID]",
		Short: "Show titles",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := showTitles(db, args[0])
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	contentCmd := &cobra.Command{
		Use:   "content [AuthorID] [TitleID]",
		Short: "Show content",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			err := showContent(db, args[0], args[1])
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	queryCmd := &cobra.Command{
		Use:   "query [Query]",
		Short: "Show contents hit with query",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := queryContent(db, args[0])
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(authorsCmd, titlesCmd, contentCmd, queryCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
