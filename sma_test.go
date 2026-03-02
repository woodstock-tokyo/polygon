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
	"errors"
	"fmt"
	"testing"
)

func Test_SMA(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	sma, err := c.SimpleMovingAverage(ctx, "AAPL", &SMAOption{
		Timespan: Day,
		Window:   50,
		Limit:    1,
	})

	if err != nil {
		t.Fatal(fmt.Errorf("get sma: %w", err))
	}

	if len(sma.Results.Values) == 0 {
		t.Error("unexpected sma:", sma)
	}
}

func Test_SMAErrorNotFound(t *testing.T) {
	ctx := context.Background()
	c := NewClient(token)
	_, err := c.SimpleMovingAverage(ctx, "NOT_A_SYMBOL", &SMAOption{
		Timespan: Day,
		Window:   50,
		Limit:    1,
	})

	if !errors.Is(err, ErrSMANoResults) {
		t.Fatal("unexpected error:", err)
	}
}
