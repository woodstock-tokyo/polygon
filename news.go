package polygon

import (
	"context"
)

// News models a news item either for the market or for an individual stock.
type News struct {
	Results []struct {
		ID        string `json:"id"`
		Publisher struct {
			Name     string `json:"name"`
			Homepage string `json:"homepage_url"`
			Logo     string `json:"logo_url"`
			Favicon  string `json:"favicon_url"`
		} `json:"publisher"`
		Title        string   `json:"title"`
		Author       string   `json:"author"`
		PublishedUTC string   `json:"published_utc"`
		ArticleURL   string   `json:"article_url"`
		Tickers      []string `json:"tickers"`
		ImageURL     string   `json:"image_url"`
		Description  string   `json:"description"`
		Keywords     []string `json:"keywords"`
	} `json:"results"`
	Status  string `json:"status"`
	Count   int    `json:"count"`
	NextURL string `json:"next_url"`
}

// NewsOption option for fetching news
type NewsOption struct {
	Ticker                     string `url:"ticker"`
	Published_LessThan         string `url:"published_utc.lt,omitempty"`
	Published_LessThanEqual    string `url:"published_utc.lte,omitempty"`
	Published_GreaterThan      string `url:"published_utc.gt,omitempty"`
	Published_GreaterThanEqual string `url:"published_utc.gte,omitempty"`
	Order                      Order  `url:"order,omitempty"`
	Sort                       string `url:"sort,omitempty"`
}

// News retrieves the given number of news articles for the given stock symbol.
func (c Client) News(ctx context.Context, ticker string, opt *NewsOption) (News, error) {
	if opt == nil {
		opt = &NewsOption{}
	}
	opt.Ticker = ticker

	n := News{}
	endpoint, err := c.endpointWithOpts("/reference/news", opt)
	if err != nil {
		return n, err
	}

	err = c.GetJSON(ctx, endpoint, &n)
	return n, err
}
