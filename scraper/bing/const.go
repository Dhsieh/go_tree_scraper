package bingscraper

// This is where all the const values in BingScraper are stored at
const (
	// Template for bing's API to grab images.
	bingPhotoUrl              string = "https://www.bing.com/images/async?q=%s+tree&first=1&mmasync=1"
	bingKeywordPhotoUrl       string = "https://www.bing.com/images/async?q=%s&first=1&mmasync=1&count=100"
	bingKeywordPhotoNumberUrl string = "https://www.bing.com/images/async?q=%s&first=%d&mmasync=1"
)
