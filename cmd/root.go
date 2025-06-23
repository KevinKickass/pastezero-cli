package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var VersionInfo = "PasteZero CLI"

var rootCmd = &cobra.Command{
	Use:   "pastezero",
	Short: "PasteZero CLI Tool",
	Long:  "Verschlüsselter Dateitransfer mit pastezero.de (Zero Trust, E2EE).",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version anzeigen",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VersionInfo)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringP("api", "a", "https://api.pastezero.de", "PasteZero API URL")
}
