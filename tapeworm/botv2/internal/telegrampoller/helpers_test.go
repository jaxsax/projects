package telegrampoller_test

import (
	"testing"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/telegrampoller"
	"github.com/stretchr/testify/require"
)

func TestURLExtractor(t *testing.T) {
	t.Run("emojis in text", func(t *testing.T) {
		text := "🆙 https://github.com/upscayl/upscayl"
		offset := 3
		length := 34

		link := telegrampoller.ExtractURL(text, offset, length)
		require.Equal(t, "https://github.com/upscayl/upscayl", link)
	})
}
