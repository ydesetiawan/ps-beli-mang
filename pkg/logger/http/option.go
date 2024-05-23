package http

type Option func(c *Client)

func WithMaskingFunc(fn MaskingFunc) Option {
	return func(c *Client) {
		c.maskingFunc = fn
	}
}

func WithSkipFunc(fn SkipFunc) Option {
	return func(c *Client) {
		c.skipFunc = fn
	}
}
