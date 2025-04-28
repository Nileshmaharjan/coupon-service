# Coupon Issuance Service

A Go + ConnectRPC backend for issuing a limited number of coupons per campaign, with support for both in-memory and Postgres-backed stores.

---

## Prerequisites

- Go 1.21+ installed
- (Optional) Postgres database (e.g., Supabase) with the following schema:

```sql
CREATE TABLE campaigns (
  id         TEXT      PRIMARY KEY,
  name       TEXT      NOT NULL,
  total      INTEGER   NOT NULL,
  start_time TIMESTAMPTZ NOT NULL,
  issued     INTEGER   NOT NULL DEFAULT 0
);

CREATE TABLE coupon_codes (
  campaign_id TEXT      REFERENCES campaigns(id),
  code        TEXT      NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(campaign_id, code)
);

ALTER TABLE coupon_codes
  ADD CONSTRAINT uq_coupon UNIQUE(campaign_id, code);
```

- (Optional) `psql` or Supabase SQL editor to run migrations
- `curl` for manual testing
- (Optional) `fasthttp` for load testing:

```bash
go install github.com/valyala/fasthttp@v1.61.0
```

---

## Setup

Clone the repo:

```bash
git clone https://github.com/Nileshmaharjan/coupon-service.git
cd coupon
```

Configure `configs/config.yaml`:

```yaml
server:
  port: 8080

database:
  dsn: "postgresql://USER:PASS@HOST:5432/postgres?sslmode=require"

store:
  type: "postgres"
```

Fetch dependencies:

```bash
go mod tidy
```

---

## Running the Server

```bash
go run cmd/server/main.go --config=configs/config.yaml
```

You should see:

```text
listening :8080
```

---

## API Endpoints

All RPCs are exposed over HTTP POST at:

```
http://localhost:8080/coupon.CouponService/<Method>
```

---

### Create a Campaign

Starts immediately with total=N coupons:

```bash
curl -v http://localhost:8080/coupon.CouponService/CreateCampaign \
  -H 'Content-Type: application/json' \
  -d "{\"name\":\"SpringPromo\",\"total\":10000,\"startTime\":\"$(date -u +'%Y-%m-%dT%H:%M:%SZ')\"}"
```

**Response:**

```json
{ "campaignId": "3bece6fa" }
```

---

### Issue a Coupon

Request one coupon from an existing campaign:

```bash
curl -v http://localhost:8080/coupon.CouponService/IssueCoupon \
  -H 'Content-Type: application/json' \
  -d '{"campaignId":"b2c7f597"}' # campaign id
```

**Response (success):**

```json
{ "code": "가1나2다3", "error": "" }
```

**Response (sold out):**

```json
{ "code": "", "error": "sold out" }
```

---

### Get Campaign Details

Fetch the campaign info and all issued codes:

```bash
curl -v http://localhost:8080/coupon.CouponService/GetCampaign \
  -H 'Content-Type: application/json' \
  -d '{"campaignId":"b2c7f597"}' # campaign id
```

**Response:**

```json
{
  "campaignId": "3bece6fa",
  "name": "SpringPromo",
  "total": 5,
  "startTime": "2025-04-27T12:34:56Z",
  "issuedCodes": ["가1나2다3", "라4마5바6"]
}
```

---

## Load Testing

Place this benchmark in `scripts/load_test.go`, then run:

```bash
go get github.com/valyala/fasthttp@v1.61.0
go test -bench=BenchmarkConcurrentIssue scripts/load_test.go -benchtime=2s
```

It will:

- Create a 1-coupon campaign in memory
- Fire 10,000 parallel `IssueCoupon` calls
- Assert you never issue more than 1 coupon

You should see output like:

```text
--- PASS: BenchmarkConcurrentIssue
    load_test.go:…
        Issued: 1 codes, SoldOut: 9999, Errors: 0
PASS
```

---

## Results

### Campaigns

![Campaigns](results/campaigns.png)

### Coupon Codes

![Coupon Codes](results/couponCodes.png)

### Create Campaign

![Create Campaign](results/createCampaign.png)

### Issue Coupon

![Issue Coupon](results/issueCoupon.png)

### Get Campaigns

![Get Campaigns](results/getCampaign.png)

### Load Testing

![Load Testing](results/loadTesting.png)

---
# coupon-service
