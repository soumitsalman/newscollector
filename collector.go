package main

import (
	"encoding/csv"
	"log"
	"os"

	ds "github.com/soumitsalman/beansack/sdk"
	datautils "github.com/soumitsalman/data-utils"
	"github.com/soumitsalman/document-loader/document"
	"github.com/soumitsalman/document-loader/loaders"
)

const _SITEMAPS_CSV = "./sitemaps.csv"

type NewsSiteCollector struct {
	site_loaders []*loaders.WebLoader
	store_func   func([]ds.Bean)
}

func NewCollector(store_func func([]ds.Bean)) NewsSiteCollector {
	return NewsSiteCollector{
		site_loaders: createSiteLoaders(),
		store_func:   store_func,
	}
}

func (collector NewsSiteCollector) Collect() {
	for _, loader := range collector.site_loaders {
		docs := loader.LoadSite()
		log.Println(len(docs), "new beans found from", loader.Config.Sitemap)
		// storeNewBeans(docs)
		collector.store_func(toBeans(docs))
	}
}

func readSitemapsCSV() [][]string {
	file, _ := os.Open(_SITEMAPS_CSV)
	defer file.Close()
	sitemaps, _ := csv.NewReader(file).ReadAll()
	// ignore the header
	return sitemaps[1:]
}

func createSiteLoaders() []*loaders.WebLoader {
	site_loaders := datautils.Transform(readSitemapsCSV(), func(item *[]string) *loaders.WebLoader {
		return loaders.NewDefaultNewsSitemapLoader(2, (*item)[0])
	})
	return append(site_loaders,
		// this is a specialied loader
		loaders.NewYCHackerNewsSiteLoader(),
	)
}

func toBeans(docs []*document.Document) []ds.Bean {
	beans := make([]ds.Bean, len(docs))
	for i, doc := range docs {
		beans[i].Url = doc.URL
		beans[i].Source = doc.Source
		beans[i].Title = doc.Title
		beans[i].Kind = doc.Kind
		beans[i].Text = doc.Text
		beans[i].Author = doc.Author
		beans[i].Created = doc.PublishDate
		beans[i].Keywords = doc.Keywords
		if doc.Comments > 0 || doc.Likes > 0 {
			beans[i].MediaNoise = &ds.MediaNoise{
				BeanUrl:       doc.URL,
				Source:        doc.Source,
				Comments:      doc.Comments,
				ThumbsupCount: doc.Likes,
			}
		}
	}
	return beans
}
