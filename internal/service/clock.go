package service

type Clock interface{ Now() string }

type FixedClock struct{ Value string }

func (c FixedClock) Now() string {
	if c.Value == "" {
		return "2026-01-01T00:00:00Z"
	}
	return c.Value
}

type SequenceClock struct {
	Values []string
	index  int
}

func (c *SequenceClock) Now() string {
	if len(c.Values) == 0 {
		return "2026-01-01T00:00:00Z"
	}
	if c.index >= len(c.Values) {
		return c.Values[len(c.Values)-1]
	}
	value := c.Values[c.index]
	c.index++
	return value
}
