package store

import (
    "context"
    "database/sql"
    "errors"
    "time"

    "github.com/Nileshmaharjan/coupon-service/internal/coupon"
    _ "github.com/lib/pq"
    "github.com/lib/pq"
)

const (
    maxTxnRetries = 3  
)

type PostgresStore struct {
    db *sql.DB
}


func NewPostgresStore(dsn string) (*PostgresStore, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    if err := db.PingContext(context.Background()); err != nil {
        return nil, err
    }
    return &PostgresStore{db: db}, nil
}

func (p *PostgresStore) Create(c *coupon.Campaign) error {
    _, err := p.db.ExecContext(
        context.Background(),
        `INSERT INTO campaigns(id, name, total, start_time) VALUES($1,$2,$3,$4)`,
        c.ID, c.Name, c.Total, c.StartTime,
    )
    return err
}

func (p *PostgresStore) Get(id string) (*coupon.Campaign, error) {
    var (
        name      string
        total     int32
        startTime time.Time
        issued    int32
    )
    err := p.db.QueryRowContext(
        context.Background(),
        `SELECT name, total, start_time, issued
           FROM campaigns
          WHERE id = $1`,
        id,
    ).Scan(&name, &total, &startTime, &issued)
    if err == sql.ErrNoRows {
        return nil, errors.New("not found")
    } else if err != nil {
        return nil, err
    }

    c := &coupon.Campaign{
        ID:        id,
        Name:      name,
        Total:     total,
        StartTime: startTime,
    }
    c.SetIssued(int(issued))

    rows, err := p.db.QueryContext(
        context.Background(),
        `SELECT code FROM coupon_codes WHERE campaign_id = $1`,
        id,
    )
    if err == nil {
        defer rows.Close()
        for rows.Next() {
            var code string
            if rows.Scan(&code) == nil {
                c.AppendCode(code)
            }
        }
    }
    return c, nil
}

func (p *PostgresStore) Issue(id string, now time.Time) (string, error) {
    return p.issueWithKey(id, now, "")
}

func (p *PostgresStore) issueWithKey(
    campaignID string,
    now time.Time,
    idempotencyKey string,
) (string, error) {
    if idempotencyKey != "" {
        var existing string
        err := p.db.QueryRowContext(
            context.Background(),
            `SELECT code
               FROM idempotency_keys
              WHERE campaign_id=$1 AND key=$2`,
            campaignID, idempotencyKey,
        ).Scan(&existing)
        if err == nil {
            return existing, nil
        }
        if err != sql.ErrNoRows {
            return "", err
        }
    }

    var lastErr error
    for attempt := 1; attempt <= maxTxnRetries; attempt++ {
        code, err := p.issueTx(campaignID, now, idempotencyKey)
        if err == nil {
            return code, nil
        }
        lastErr = err
        var pgErr *pq.Error
        if errors.As(err, &pgErr) && (pgErr.Code == "40001" || pgErr.Code == "40P01") {
            time.Sleep(50 * time.Millisecond)
            continue
        }
        break
    }
    return "", lastErr
}

func (p *PostgresStore) issueTx(
    campaignID string,
    now time.Time,
    idempotencyKey string,
) (string, error) {
    ctx := context.Background()
    tx, err := p.db.BeginTx(ctx, nil)
    if err != nil {
        return "", err
    }
    defer tx.Rollback()

    // Lock campaign row
    var total, issued int32
    var startTime time.Time
    if err := tx.QueryRowContext(
        ctx,
        `SELECT total, issued, start_time
           FROM campaigns
          WHERE id=$1 FOR UPDATE`,
        campaignID,
    ).Scan(&total, &issued, &startTime); err != nil {
        if err == sql.ErrNoRows {
            return "", errors.New("not found")
        }
        return "", err
    }

    // enforce start time & stock
    if now.Before(startTime) {
        return "", errors.New("not started")
    }
    if issued >= total {
        return "", errors.New("sold out")
    }

    // increment issued counter
    if _, err := tx.ExecContext(
        ctx,
        `UPDATE campaigns SET issued = issued + 1 WHERE id = $1`,
        campaignID,
    ); err != nil {
        return "", err
    }

    // generate & insert coupon code, retry collisions
    var code string
    for i := 0; i < 5; i++ {
        code = coupon.GenCode(10)
        _, err := tx.ExecContext(
            ctx,
            `INSERT INTO coupon_codes(campaign_id, code, created_at)
                    VALUES($1,$2,$3)`,
            campaignID, code, now,
        )
        if err == nil {
            break
        }
        if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
            tx.ExecContext(ctx,
                `UPDATE campaigns SET issued = issued - 1 WHERE id = $1`,
                campaignID,
            )
            continue
        }
        return "", err
    }

    // record idempotency key if provided
    if idempotencyKey != "" {
        if _, err := tx.ExecContext(
            ctx,
            `INSERT INTO idempotency_keys(campaign_id, key, code, created_at)
                    VALUES($1, $2, $3, $4)`,
            campaignID, idempotencyKey, code, now,
        ); err != nil {
            return "", err
        }
    }

    if err := tx.Commit(); err != nil {
        return "", err
    }
    return code, nil
}
