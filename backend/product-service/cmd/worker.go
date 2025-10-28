package cmd

import (
	"fmt"
	"product-service/internal/adapter/message"

	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Menjalankan worker untuk consume RabbitMQ dan index ke Elasticsearch",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Worker untuk Elasticsearch Indexing sedang berjalan...")
		message.StartConsumer()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
