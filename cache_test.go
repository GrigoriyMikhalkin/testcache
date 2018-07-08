package cache

import (
  "time"
  "sync"
  "testing"
  "strconv"
)

func TestCreate(t *testing.T) {
  duration := time.Duration(1) * time.Second
  cache := InitCache(duration)
  testVal := "test"

  // Add entry successfully
  res, err := cache.Create(testVal, testVal)
  if (err != nil) || (res != testVal) {
    t.Error("Expected 'test', got ", res)
  }

  // Key already used and didn't expired
  res, err = cache.Create(testVal, testVal)
  if (err == nil) || (res != testVal) {
    t.Error("Expected 'test', got ", res)
  }

  // Key used, but expired
  time.Sleep(duration)
  testVal2 := "test2"
  res, err = cache.Create(testVal, testVal2)
  if (err != nil) || (res != testVal2) {
    t.Error("Expected 'test2', got ", res)
  }
}

func TestGetConcurrent(t *testing.T) {
  var wg sync.WaitGroup

  duration := time.Duration(1) * time.Second
  cache := InitCache(duration)
  testVal := "test"

  for i := 1; i <= 10; i++ {
    wg.Add(1)
    go func (i int) {
      val, _ := cache.Create(testVal, testVal + strconv.Itoa(i))
      val2, err := cache.Get(testVal)
      if (err != nil) || !(val == val2) {
        t.Error(val, " not equal to ", val2)
      }

      wg.Done()
    }(i)
  }

  wg.Wait()
}

func TestUpdateConcurrent(t *testing.T) {
  var wg sync.WaitGroup

  duration := time.Duration(1) * time.Second
  cache := InitCache(duration)
  testVal := "test"
  cache.Create(testVal, testVal)

  for i := 1; i <= 10; i++ {
    wg.Add(1)
    go func (i int) {
      val, err := cache.Update(testVal, testVal + strconv.Itoa(i))
      if (err != nil) || !(val == testVal + strconv.Itoa(i)) {
        t.Error(val, " not equal to ", testVal + strconv.Itoa(i))
      }

      wg.Done()
    }(i)
  }

  wg.Wait()
}

func TestRemoveConcurrent(t *testing.T) {
  var wg sync.WaitGroup

  duration := time.Duration(1) * time.Second
  cache := InitCache(duration)
  testVal := "test"

  for i := 1; i <= 10; i++ {
    wg.Add(1)
    go func (i int) {
      cache.Create(testVal + strconv.Itoa(i), testVal)
      err := cache.Remove(testVal + strconv.Itoa(i))
      if (err != nil) {
        t.Error("Couldn't remove ", testVal)
      }

      wg.Done()
    }(i)
  }

  wg.Wait()
}

func TestListOfKeys(t *testing.T) {
  duration := time.Duration(1) * time.Second
  cache := InitCache(duration)
  testVal := "test"
  testVal2 := "test2"
  cache.Create(testVal, testVal)
  time.Sleep(duration)
  cache.Create(testVal2, testVal2)

  keys := cache.ListOfKeys()
  if !(len(keys) == 1) {
    t.Error("ListOfKeys returned wrong number of elements")
  }
}
