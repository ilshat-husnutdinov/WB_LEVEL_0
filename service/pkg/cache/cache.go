package cache

import (
	"errors"
	"fmt"
	"service/service/pkg/database"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
}

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

// New инициирует кэш
func New(defaultExpiration, cleanupInterval time.Duration, db *sqlx.DB) (*Cache, error) {
	const op = "cache.cache.New"

	items := make(map[string]Item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	orders, err := database.GetAllOrders(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for order_id, order := range orders {
		cache.Set(order_id, order, 0)
	}
	return &cache, nil
}

// Set помещает значение в кэш по переданному ключу
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {

	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

// Get получает значение из кэша по переданному ключу
func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return nil, false
	}

	if item.Expiration > 0 {

		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}

	}

	return item.Value, true
}

// Delete удаляет значение из кэша по переданному ключу
func (c *Cache) Delete(key string) error {

	c.Lock()

	defer c.RUnlock()

	if _, found := c.items[key]; !found {
		return errors.New("key not found")
	}

	delete(c.items, key)

	return nil
}

// StartGC() запускает GC()
func (c *Cache) StartGC() {
	go c.GC()
}

// GC() реализует сборку мусора в кэше
func (c *Cache) GC() {

	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)

		}

	}

}

// expiredKeys возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() (keys []string) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}

// ClearCache очищает кэш от данных (используется перед остановкой/падением сервиса)
func (c *Cache) ClearCache() {
	clear(c.items)
}

// GetItems возвращает все значения из кэша
func (c *Cache) GetItems() map[string]Item {
	return c.items
}
