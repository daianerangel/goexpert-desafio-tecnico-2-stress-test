/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goexpert-desafio-tecnico-2-stress-test",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("url", "", "Url to be tested")
	rootCmd.PersistentFlags().Int("requests", 0, "Number of requests")
	rootCmd.PersistentFlags().Int("concurrency", 0, "Number of concurrent requests")
	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.MarkPersistentFlagRequired("requests")
	rootCmd.MarkPersistentFlagRequired("concurrency")
	loaderTest := newLoaderTest()
	rootCmd.AddCommand(loaderTest)
}

func newLoaderTest() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Run a new loader test",
		Long:  `Run a new loader test`,
		RunE:  runLoaderTest(),
	}
}
func runLoaderTest() RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		url, _ := cmd.Flags().GetString("url")
		fmt.Println("Url:", url)
		totalReqs, _ := cmd.Flags().GetInt("requests")
		fmt.Println("Requests:", totalReqs)
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		fmt.Println("Concurrency:", concurrency)

		if totalReqs <= 0 || concurrency <= 0 {
			var err = errors.New("parameters --requests and --concurrency must be provided and greater than 0")
			return err
		}

		start := time.Now()
		results := make(chan int, totalReqs)
		var wg sync.WaitGroup

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < totalReqs/concurrency; j++ {
					resp, err := http.Get(url)
					if err != nil {
						results <- 0
						continue
					}
					results <- resp.StatusCode
					resp.Body.Close()
				}
			}()
		}

		wg.Wait()
		close(results)

		totalTime := time.Since(start)
		statusCount := make(map[int]int)
		for status := range results {
			statusCount[status]++
		}

		fmt.Printf("Total time: %v\n", totalTime)
		fmt.Printf("Total requests: %d\n", totalReqs)
		fmt.Printf("HTTP 200 responses: %d\n", statusCount[200])
		for status, count := range statusCount {
			if status != 200 {
				fmt.Printf("HTTP %d responses: %d\n", status, count)
			}
		}

		return nil
	}

}
