package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const defaultDSN = "postgres://postgres@localhost:5432/postgres?sslmode=disable"

type TryLockParams struct {
	Name       string // Lock Name: unique identifier for the lock
	Owner      string // Lock Owner: identifier for the entity requesting the lock
	TTLSeconds int    // Time-To-Live: duration in seconds for the lock
}

type TryLockResult struct {
	ExpiresAt time.Time // Expiration time of the lock
	Acquired  bool      // Whether the lock was successfully acquired
}

// TryLock attempts to acquire a distributed lock.
// Returns the expiration time, whether the lock was acquired, and any error.
func TryLock(ctx context.Context, db *sql.DB, params TryLockParams) (TryLockResult, error) {
	// Atomic operation: insert if not exists, or update if expired
	const q = `
		INSERT INTO locks (name, owner, expires_at)
		VALUES ($1, $2, NOW() + make_interval(secs => $3))
		ON CONFLICT (name) DO UPDATE
		SET owner = EXCLUDED.owner, expires_at = EXCLUDED.expires_at
		WHERE locks.expires_at <= NOW()
		RETURNING expires_at;
	`
	var expires time.Time
	err := db.QueryRowContext(ctx, q, params.Name, params.Owner, params.TTLSeconds).Scan(&expires)
	if err == sql.ErrNoRows {
		// Lock is held by another owner and not expired
		return TryLockResult{}, nil
	}
	if err != nil {
		return TryLockResult{}, err
	}
	return TryLockResult{ExpiresAt: expires, Acquired: true}, nil
}

type LockParams struct {
	Name             string        // Lock Name: unique identifier for the lock
	Owner            string        // Lock Owner: identifier for the entity requesting the lock
	TTLSeconds       int           // Time-To-Live: duration in seconds for the lock
	IntervalDuration time.Duration // Retry interval duration
}

type LockResult struct {
	ExpiresAt time.Time // Expiration time of the lock
}

// Lock continuously attempts to acquire a distributed lock until successful.
func Lock(ctx context.Context, db *sql.DB, params LockParams) (LockResult, error) {
	for {
		result, err := TryLock(ctx, db, TryLockParams{
			Name:       params.Name,
			Owner:      params.Owner,
			TTLSeconds: params.TTLSeconds,
		})
		if err != nil {
			return LockResult{}, err
		}
		if result.Acquired {
			return LockResult{ExpiresAt: result.ExpiresAt}, nil
		}
		// Wait before retrying
		time.Sleep(params.IntervalDuration)
	}
}

type RefreshLockParams struct {
	Name       string // Lock Name: unique identifier for the lock
	Owner      string // Lock Owner: identifier for the entity refreshing the lock
	TTLSeconds int    // Time-To-Live: new duration in seconds for the lock
}

type RefreshLockResult struct {
	ExpiresAt time.Time // New expiration time of the lock
	Refreshed bool      // Whether the refresh succeeded
}

// RefreshLock extends the lock's TTL if we still own it.
// Returns the new expiration time, whether the refresh succeeded, and any error.
func RefreshLock(ctx context.Context, db *sql.DB, params RefreshLockParams) (RefreshLockResult, error) {
	const q = `
		UPDATE locks
		SET expires_at = NOW() + make_interval(secs => $2)
		WHERE name = $1 AND owner = $3
		RETURNING expires_at;
	`
	var expires time.Time
	err := db.QueryRowContext(ctx, q, params.Name, params.TTLSeconds, params.Owner).Scan(&expires)
	if err == sql.ErrNoRows {
		// Lock was lost or expired
		return RefreshLockResult{}, nil
	}
	if err != nil {
		return RefreshLockResult{}, err
	}
	return RefreshLockResult{ExpiresAt: expires, Refreshed: true}, nil
}

type UnlockParams struct {
	Name  string // Lock Name: unique identifier for the lock
	Owner string // Lock Owner: identifier for the entity releasing the lock
}

type UnlockResult struct {
	Released bool // Whether the lock was released
}

// Unlock releases the lock if we still own it.
// Returns whether the lock was released and any error.
func Unlock(ctx context.Context, db *sql.DB, params UnlockParams) (UnlockResult, error) {
	const q = `DELETE FROM locks WHERE name = $1 AND owner = $2;`
	res, err := db.ExecContext(ctx, q, params.Name, params.Owner)
	if err != nil {
		return UnlockResult{}, err
	}
	rowsAffected, _ := res.RowsAffected()
	return UnlockResult{Released: rowsAffected > 0}, nil
}

type CheckLockStatusParams struct {
	Name  string // Lock Name: unique identifier for the lock
	Owner string // Lock Owner: identifier for the entity checking the lock
}

type CheckLockStatusResult struct {
	Valid     bool      // Whether the lock is still valid
	ExpiresAt time.Time // Expiration time of the lock
}

// CheckLockStatus checks if a lock is still valid for the given owner.
// Returns whether the lock is still valid and any error.
func CheckLockStatus(ctx context.Context, db *sql.DB, params CheckLockStatusParams) (CheckLockStatusResult, error) {
	const q = `
SELECT expires_at 
FROM locks 
WHERE name = $1 AND owner = $2 AND expires_at > NOW();
`
	var expires time.Time
	err := db.QueryRowContext(ctx, q, params.Name, params.Owner).Scan(&expires)
	if err == sql.ErrNoRows {
		// Lock expired or lost
		return CheckLockStatusResult{Valid: false, ExpiresAt: time.Time{}}, nil
	}
	if err != nil {
		return CheckLockStatusResult{Valid: false, ExpiresAt: time.Time{}}, err
	}
	return CheckLockStatusResult{Valid: true, ExpiresAt: expires}, nil
}

// initDB creates the locks table if it doesn't exist
func initDB(ctx context.Context, db *sql.DB) error {
	const createTableSQL = `
		CREATE TABLE IF NOT EXISTS locks (
		name TEXT PRIMARY KEY,
		owner TEXT NOT NULL,
		expires_at TIMESTAMPTZ NOT NULL
		);
	`
	_, err := db.ExecContext(ctx, createTableSQL)
	return err
}

func worker(ctx context.Context, db *sql.DB, workerID int, lockName string) {
	owner := fmt.Sprintf("worker-%d", workerID)

	log.Printf("[Worker %d] Starting, attempting to acquire lock '%s'", workerID, lockName)

	// Try to acquire lock with blocking retry
	result, err := Lock(ctx, db, LockParams{
		Name:             lockName,
		Owner:            owner,
		TTLSeconds:       5,
		IntervalDuration: 500 * time.Millisecond,
	})
	if err != nil {
		log.Printf("[Worker %d] ❌ Failed to acquire lock: %v", workerID, err)
		return
	}

	log.Printf("[Worker %d] ✅ Acquired lock! Expires at %s", workerID, result.ExpiresAt.Format("15:04:05"))

	// Simulate work while holding the lock
	workDuration := time.Duration(2+workerID%3) * time.Second
	log.Printf("[Worker %d] 🔨 Working for %v...", workerID, workDuration)

	// Simulate some work with periodic status checks
	workStart := time.Now()
	for time.Since(workStart) < workDuration {
		time.Sleep(1 * time.Second)

		// Check if we still own the lock
		status, err := CheckLockStatus(ctx, db, CheckLockStatusParams{
			Name:  lockName,
			Owner: owner,
		})
		if err != nil {
			log.Printf("[Worker %d] ⚠️ Error checking lock status: %v", workerID, err)
			return
		}
		if !status.Valid {
			log.Printf("[Worker %d] ❌ Lock expired! Lost ownership", workerID)
			return
		}
		log.Printf("[Worker %d] ⏳ Still working... lock valid until %s", workerID, status.ExpiresAt.Format("15:04:05"))
	}

	log.Printf("[Worker %d] ✅ Work completed!", workerID)

	// Release the lock
	releaseResult, err := Unlock(ctx, db, UnlockParams{
		Name:  lockName,
		Owner: owner,
	})
	if err != nil {
		log.Printf("[Worker %d] ⚠️ Error releasing lock: %v", workerID, err)
		return
	}
	if releaseResult.Released {
		log.Printf("[Worker %d] 🔓 Lock released successfully", workerID)
	} else {
		log.Printf("[Worker %d] ⚠️ Failed to release lock (already released or expired)", workerID)
	}
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = defaultDSN
		log.Printf("Using default DSN (set DATABASE_URL to override): %s", dsn)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Test connection
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL")

	// Initialize database
	if err := initDB(ctx, db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("✅ Database initialized")

	// Clean up any existing locks
	_, _ = db.ExecContext(ctx, "DELETE FROM locks")
	log.Println("🧹 Cleaned up existing locks")
	log.Println("")

	// Demo: Multiple goroutines competing for the same lock
	lockName := "shared-resource"
	numWorkers := 5

	log.Printf("🚀 Starting %d workers competing for lock '%s'\n", numWorkers, lockName)
	log.Println("═══════════════════════════════════════════════════════════")

	// Launch workers concurrently
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, db, id, lockName)
		}(i)
		// Stagger the start slightly to make logs clearer
		time.Sleep(100 * time.Millisecond)
	}

	// Wait for all workers to complete
	wg.Wait()

	log.Println("\n═══════════════════════════════════════════════════════════")
	log.Println("✅ All workers completed!")

	// Clean up locks before next demo
	_, _ = db.ExecContext(ctx, "DELETE FROM locks")
	time.Sleep(1 * time.Second)
}
