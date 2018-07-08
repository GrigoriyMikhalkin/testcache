package cache

import (
  "fmt"
  "sync"
  "time"
  "errors"
  "reflect"
)

type value struct {
  val     interface{}
  created time.Time
}

type Cache struct {
  sync.RWMutex
  cache map[string]*value
  ttl   time.Duration
}

func InitCache(ttl time.Duration) *Cache {
  /*
   * Instantiates Cache
   *
   * :param: ttl: time.Duration: Time to live for cache keys in seconds
   */

  return &Cache{
    cache: make(map[string]*value),
    ttl:   ttl,
  }
}

func (c *Cache) expired(val *value) bool {
  return time.Now().Add(-c.ttl).After(val.created)
}

func (c *Cache) removeExpiredKey(key string) {
  c.Lock()
  defer c.Unlock()

  if res, ok := c.cache[key]; ok && c.expired(res) {
    delete(c.cache, key)
  }
}

func (c *Cache) Create(key string, val interface{}) (interface{}, error) {
  /*
   * Add stuff to cache
   */
  c.Lock()
  defer c.Unlock()

  if res, ok := c.cache[key]; ok && !c.expired(res) {
    return res.val, errors.New(
      fmt.Sprintf("Key %s already used", key),
    )
  }

  c.cache[key] = &value{
    val: val,
    created: time.Now(),
  }
  return c.cache[key].val, nil
}

func (c *Cache) Get(key string) (interface{}, error) {
  /*
   * Get value stored in cache
   */
   c.RLock()
   defer c.RUnlock()

   res, ok := c.cache[key]
   if !ok || c.expired(res) {
     if ok {
       go c.removeExpiredKey(key)
     }
     return nil, errors.New(
       fmt.Sprintf("Key %s not found", key),
     )
   }

   return res.val, nil
}

func (c *Cache) Update(key string, val interface{}) (interface{}, error) {
  /*
   * Update key value in cache
   */
   c.Lock()
   defer c.Unlock()

   if _, ok := c.cache[key]; !ok {
     return nil, errors.New(
       fmt.Sprintf("Key %s not stored in cache", key),
     )
   }

   c.cache[key] = &value{
     val: val,
     created: time.Now(),
   }
   return c.cache[key].val, nil
}

func (c *Cache) Remove(key string) error {
  /*
   * Update key value in cache
   */
   c.Lock()
   defer c.Unlock()

   if res, ok := c.cache[key]; !ok || c.expired(res) {
     delete(c.cache, key)
     return errors.New(
       fmt.Sprintf("Key %s not stored in cache", key),
     )
   }

   delete(c.cache, key)
   return nil
}

func (c *Cache) ListOfKeys() []string {
  c.RLock()
  defer c.RUnlock()

  var key string
  var strkeys []string
  keys := reflect.ValueOf(c.cache).MapKeys()
  for i := 0; i < len(keys); i++ {
    key = keys[i].String()
    if !c.expired(c.cache[key]) {
      strkeys = append(strkeys, key)
    } else {
      go c.removeExpiredKey(key)
    }
  }
  return strkeys
}
