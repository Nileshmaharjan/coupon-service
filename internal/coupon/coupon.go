
package coupon

import (
  "errors"
  "sync"
  "sync/atomic"
  "time"

  mrand "math/rand"
)

func init() {
  mrand.Seed(time.Now().UnixNano())
}

type Store interface {
  Create(*Campaign) error
  Get(string) (*Campaign, error)
  Issue(string, time.Time) (string, error)
}
    
type Campaign struct {
  ID        string
  Name      string
  Total     int32
  StartTime time.Time

  issuedCnt int32     
  codes     []string  
  mu        sync.Mutex // make sure thread safe
}


func (c *Campaign) SetIssued(n int) {
  atomic.StoreInt32(&c.issuedCnt, int32(n))
}

func (c *Campaign) Issue(now time.Time) (string, error) {
  // check if the campaign has started
  if now.Before(c.StartTime) {
    return "", errors.New("not started")
  }
  // atomically reserve one coupon
  cur := atomic.AddInt32(&c.issuedCnt, 1)
  if cur > c.Total {
    // roll back if we exceed the total
    atomic.AddInt32(&c.issuedCnt, -1)
    return "", errors.New("sold out")
  }
  // generate a unique code
  code := GenCode(10)
  c.mu.Lock()
  c.codes = append(c.codes, code)
  c.mu.Unlock()
  return code, nil
}

// Return copy of all issued coupon codes
func (c *Campaign) Codes() []string {
  c.mu.Lock()
  defer c.mu.Unlock()
  dup := make([]string, len(c.codes))
  copy(dup, c.codes)
  return dup
}


// safely append a single code to the campaign
func (c *Campaign) AppendCode(code string) {
  c.mu.Lock()
  defer c.mu.Unlock()
  c.codes = append(c.codes, code)
}


// generate a random n-character coupon code
func GenCode(n int) string {
  korean := []rune("가나다라마바사아자차카타파하")
  digits := []rune("0123456789")

  b := make([]rune, n)
  for i := 0; i < n; i++ {
    if mrand.Intn(2) == 0 {
      b[i] = korean[mrand.Intn(len(korean))]
    } else {
      b[i] = digits[mrand.Intn(len(digits))]
    }
  }
  return string(b)
}
