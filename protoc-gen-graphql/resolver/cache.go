package resolver

// Cache is struct for checking key in stacks but never keep its value.
type Cache struct {
	c map[string]struct{}
}

func NewCache() *Cache {
	return &Cache{
		c: make(map[string]struct{}),
	}
}

// Exists returns true if key exists in stack
func (c *Cache) Exists(key string) bool {
	_, ok := c.c[key]
	return ok
}

// Add adds stack with key
func (c *Cache) Add(key string) {
	c.c[key] = struct{}{}
}
