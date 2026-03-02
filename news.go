// Copyright (c) 2026 Woodstock K.K.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
		AmpURL       string   `json:"amp_url"`
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
	Limit                      uint   `url:"limit,omitempty"`
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
