package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"

    "gopkg.in/yaml.v2"
    couponconnect "github.com/Nileshmaharjan/coupon-service/gen/coupon/couponconnect"
    "github.com/Nileshmaharjan/coupon-service/internal/coupon"
    "github.com/Nileshmaharjan/coupon-service/internal/store"
)

type Config struct {
    Server struct {
        Port int `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        DSN string `yaml:"dsn"`
    } `yaml:"database"`
    Store struct {
        Type string `yaml:"type"`
    } `yaml:"store"`
}

func loadConfig(path string) *Config {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatalf("read config: %v", err)
    }
    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        log.Fatalf("unmarshal config: %v", err)
    }
    return &cfg
}

func main() {
    cfgPath := flag.String("config", "configs/config.yaml", "config path")
    flag.Parse()
    cfg := loadConfig(*cfgPath)

    var repo coupon.Store
    switch cfg.Store.Type {
    case "postgres":
        pg, err := store.NewPostgresStore(cfg.Database.DSN)
        if err != nil {
            log.Fatalf("postgres connect: %v", err)
        }
        repo = pg
    default:
        repo = store.NewMemoryStore()
    }

    svc := coupon.NewService(repo)
    mux := http.NewServeMux()
    path, handler := couponconnect.NewCouponServiceHandler(svc)
    mux.Handle(path, handler)

    addr := fmt.Sprintf(":%d", cfg.Server.Port)
    log.Printf("listening %s", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}
