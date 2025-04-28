package main

import (
    "encoding/json"
    "sync"
    "testing"
    "time"

    "github.com/valyala/fasthttp"
)

const (
    baseURL       = "http://localhost:8080"
    createPath    = "/coupon.CouponService/CreateCampaign"
    issuePath     = "/coupon.CouponService/IssueCoupon"
    totalCoupons  = 50
    concurrency   = 10000
)

func BenchmarkConcurrentIssue(b *testing.B) {
    client := &fasthttp.Client{}

    // create a fresh campaign
    now := time.Now().UTC().Format(time.RFC3339)
    createPayload, _ := json.Marshal(map[string]interface{}{
        "name":      "ColdPlayConcertSeoul",
        "total":     totalCoupons,
        "startTime": now,
    })

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    req.SetRequestURI(baseURL + createPath)
    req.Header.SetMethod("POST")
    req.Header.Set("Content-Type", "application/json")
    req.SetBody(createPayload)

    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)
    if err := client.Do(req, resp); err != nil || resp.StatusCode() != fasthttp.StatusOK {
        b.Fatalf("CreateCampaign failed: status=%d err=%v body=%s", resp.StatusCode(), err, resp.Body())
    }

    var createResp struct{ CampaignId string `json:"campaignId"` }
    if err := json.Unmarshal(resp.Body(), &createResp); err != nil {
        b.Fatalf("failed to parse CreateCampaign response: %v", err)
    }
    campaignID := createResp.CampaignId
    b.Logf("Created campaign %s (total=%d)", campaignID, totalCoupons)

    // prepare issue coupon payload
    issuePayload, _ := json.Marshal(map[string]string{
        "campaignId": campaignID,
    })

    // run the benchmark iterations
    for i := 0; i < b.N; i++ {
        var wg sync.WaitGroup
        results := make(chan []byte, concurrency)
        wg.Add(concurrency)

        // launch 1,000 concurrent issue coupon calls
        for j := 0; j < concurrency; j++ {
            go func() {
                defer wg.Done()
                r := fasthttp.AcquireRequest()
                defer fasthttp.ReleaseRequest(r)
                r.SetRequestURI(baseURL + issuePath)
                r.Header.SetMethod("POST")
                r.Header.Set("Content-Type", "application/json")
                r.SetBody(issuePayload)

                rs := fasthttp.AcquireResponse()
                defer fasthttp.ReleaseResponse(rs)
                if err := client.Do(r, rs); err != nil {
                    results <- []byte(`{"error":"network"}`)
                } else {
                    results <- append([]byte(nil), rs.Body()...)
                }
            }()
        }

        wg.Wait()
        close(results)

        // tally results
        codes := make(map[string]struct{})
        soldOut, errs := 0, 0
        for raw := range results {
            var res struct {
                Code  string `json:"code"`
                Error string `json:"error"`
            }
            if err := json.Unmarshal(raw, &res); err != nil {
                errs++
            } else if res.Error != "" {
                soldOut++
            } else {
                codes[res.Code] = struct{}{}
            }
        }

        // assert no over-issuance
        if len(codes) > totalCoupons {
            b.Fatalf("over-issuance: got %d unique codes, limit %d", len(codes), totalCoupons)
        }
        b.Logf("Issued: %d codes, SoldOut: %d, Errors: %d", len(codes), soldOut, errs)

        // brief pause
        time.Sleep(100 * time.Millisecond)
    }
}
