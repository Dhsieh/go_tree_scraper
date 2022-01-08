package cloud_functions

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	treeConfig "github.com/Dhsieh/tree_scraper/config"
	"github.com/Dhsieh/tree_scraper/data"
	"github.com/spf13/viper"
)

func Test(ctx context.Context, event data.GCSEvent) error {
	viper.SetConfigType("yaml")

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Error creating storage client %\nv", err)
	}

	bucketName := event.Bucket
	bucket := client.Bucket(bucketName)

	config := bucket.Object(event.Name)
	reader, err := config.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Error reading %v file %v\n", event.Name, err)
	}

	defer reader.Close()
	readFile, err := ioutil.ReadAll(reader)

	if err != nil {
		return fmt.Errorf("Error reading file %v %v\n", event.Name, err)
	}

	viper.ReadConfig(bytes.NewBuffer(readFile))

	var configuration treeConfig.Configuration
	viper.Unmarshal(&configuration)

	bingScraper := configuration.CreateBingScraper()
	fmt.Println(bingScraper)
	bingScraper.ScrapeImagesToGCS(ctx, configuration.Keyword)

	return nil
}
