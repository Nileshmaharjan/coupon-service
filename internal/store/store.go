package store

import (
  "time"
  "github.com/Nileshmaharjan/coupon-service/internal/coupon"
)

type Store interface {
  Create(*coupon.Campaign) error
  Get(string) (*coupon.Campaign, error)
  Issue(string, time.Time) (string, error)
}