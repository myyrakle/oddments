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
	Priority   int    // Priority: higher number = higher priority
	TTLSeconds int    // Time-To-Live: duration in seconds for the lock
}

type TryLockResult struct {
	ExpiresAt     time.Time // Expiration time of the lock
	Acquired      bool      // Whether the lock was successfully acquired
	QueuePosition int       // Position in the queue if not acquired (0 if acquired)
}

// TryLock attempts to acquire a priority-based distributed lock.
// If lock is held, the request is queued by priority.
func TryLock(ctx context.Context, db *sql.DB, params TryLockParams) (TryLockResult, error) {
	// Try to acquire the lock directly if it's free or expired
	const tryAcquireQ = `
		INSERT INTO priority_locks (name, owner, priority, expires_at)
		VALUES ($1, $2, $3, NOW() + make_interval(secs => $4))
		ON CONFLICT (name) DO UPDATE
		SET owner = EXCLUDED.owner, 
		    priority = EXCLUDED.priority,
		    expires_at = EXCLUDED.expires_at
		WHERE priority_locks.expires_at <= NOW()
		RETURNING expires_at;
	`
	var expires time.Time
	err := db.QueryRowContext(ctx, tryAcquireQ, params.Name, params.Owner, params.Priority, params.TTLSeconds).Scan(&expires)
	if err == nil {
		// Successfully acquired the lock
		return TryLockResult{ExpiresAt: expires, Acquired: true, QueuePosition: 0}, nil
	}
	if err != sql.ErrNoRows {
		return TryLockResult{}, err
	}

	// Lock is held by someone else, add to queue
	const enqueueQ = `
		INSERT INTO lock_queue (lock_name, owner, priority, requested_at, heartbeat_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (lock_name, owner) DO UPDATE
		SET priority = EXCLUDED.priority, 
		    requested_at = EXCLUDED.requested_at,
		    heartbeat_at = NOW();
	`
	_, err = db.ExecContext(ctx, enqueueQ, params.Name, params.Owner, params.Priority)
	if err != nil {
		return TryLockResult{}, err
	}

	// Get position in queue (excluding stale entries)
	const positionQ = `
		SELECT COUNT(*) + 1
		FROM lock_queue
		WHERE lock_name = $1 
		  AND heartbeat_at >= NOW() - make_interval(secs => 10)
		  AND (priority > $2 OR (priority = $2 AND requested_at < (
		      SELECT requested_at FROM lock_queue WHERE lock_name = $1 AND owner = $3
		  )));
	`
	var position int
	err = db.QueryRowContext(ctx, positionQ, params.Name, params.Priority, params.Owner).Scan(&position)
	if err != nil {
		return TryLockResult{}, err
	}

	return TryLockResult{Acquired: false, QueuePosition: position}, nil
}

type LockParams struct {
	Name             string        // Lock Name: unique identifier for the lock
	Owner            string        // Lock Owner: identifier for the entity requesting the lock
	Priority         int           // Priority: higher number = higher priority
	TTLSeconds       int           // Time-To-Live: duration in seconds for the lock
	IntervalDuration time.Duration // Retry interval duration
}

type LockResult struct {
	ExpiresAt time.Time // Expiration time of the lock
}

// Lock continuously attempts to acquire a priority-based distributed lock until successful.
func Lock(ctx context.Context, db *sql.DB, params LockParams) (LockResult, error) {
	// Heartbeat timeout: if a queue entry hasn't been updated in this time, it's considered dead
	const queueHeartbeatTimeout = 10 // seconds

	// First attempt: try to acquire immediately without queuing
	result, err := TryLock(ctx, db, TryLockParams{
		Name:       params.Name,
		Owner:      params.Owner,
		Priority:   params.Priority,
		TTLSeconds: params.TTLSeconds,
	})
	if err != nil {
		return LockResult{}, err
	}
	if result.Acquired {
		// Got it on first try! No need to queue
		return LockResult{ExpiresAt: result.ExpiresAt}, nil
	}

	// We're now in the queue, keep trying
	for {
		// Clean up stale queue entries before checking
		const cleanupQ = `
			DELETE FROM lock_queue 
			WHERE lock_name = $1 
			  AND heartbeat_at < NOW() - make_interval(secs => $2);
		`
		db.ExecContext(ctx, cleanupQ, params.Name, queueHeartbeatTimeout)

		// Check if we're the highest priority waiter (excluding stale entries)
		const checkTopQ = `
			SELECT owner, priority
			FROM lock_queue
			WHERE lock_name = $1
			  AND heartbeat_at >= NOW() - make_interval(secs => $2)
			ORDER BY priority DESC, requested_at ASC
			LIMIT 1;
		`
		var topOwner string
		var topPriority int
		err := db.QueryRowContext(ctx, checkTopQ, params.Name, queueHeartbeatTimeout).Scan(&topOwner, &topPriority)

		// Only try to acquire if we're the highest priority waiter
		if err == sql.ErrNoRows || (topOwner == params.Owner && topPriority == params.Priority) {
			// Try to acquire the lock (will update heartbeat via TryLock)
			result, err := TryLock(ctx, db, TryLockParams{
				Name:       params.Name,
				Owner:      params.Owner,
				Priority:   params.Priority,
				TTLSeconds: params.TTLSeconds,
			})
			if err != nil {
				return LockResult{}, err
			}
			if result.Acquired {
				// Success! Remove from queue and return
				db.ExecContext(ctx, "DELETE FROM lock_queue WHERE lock_name = $1 AND owner = $2", params.Name, params.Owner)
				return LockResult{ExpiresAt: result.ExpiresAt}, nil
			}
		} else {
			// We're waiting: update heartbeat to show we're still alive
			const updateHeartbeatQ = `
				UPDATE lock_queue 
				SET heartbeat_at = NOW() 
				WHERE lock_name = $1 AND owner = $2;
			`
			db.ExecContext(ctx, updateHeartbeatQ, params.Name, params.Owner)
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
func RefreshLock(ctx context.Context, db *sql.DB, params RefreshLockParams) (RefreshLockResult, error) {
	const q = `
		UPDATE priority_locks
		SET expires_at = NOW() + make_interval(secs => $2)
		WHERE name = $1 AND owner = $3 AND expires_at > NOW()
		RETURNING expires_at;
	`
	var expires time.Time
	err := db.QueryRowContext(ctx, q, params.Name, params.TTLSeconds, params.Owner).Scan(&expires)
	if err == sql.ErrNoRows {
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
func Unlock(ctx context.Context, db *sql.DB, params UnlockParams) (UnlockResult, error) {
	const q = `DELETE FROM priority_locks WHERE name = $1 AND owner = $2;`
	res, err := db.ExecContext(ctx, q, params.Name, params.Owner)
	if err != nil {
		return UnlockResult{}, err
	}
	rowsAffected, _ := res.RowsAffected()

	// Also remove from queue if present
	db.ExecContext(ctx, "DELETE FROM lock_queue WHERE lock_name = $1 AND owner = $2", params.Name, params.Owner)

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
func CheckLockStatus(ctx context.Context, db *sql.DB, params CheckLockStatusParams) (CheckLockStatusResult, error) {
	const q = `
		SELECT expires_at 
		FROM priority_locks 
		WHERE name = $1 AND owner = $2 AND expires_at > NOW();
	`
	var expires time.Time
	err := db.QueryRowContext(ctx, q, params.Name, params.Owner).Scan(&expires)
	if err == sql.ErrNoRows {
		return CheckLockStatusResult{}, nil
	}
	if err != nil {
		return CheckLockStatusResult{}, err
	}
	return CheckLockStatusResult{Valid: true, ExpiresAt: expires}, nil
}

// initDB creates the priority locks and queue tables if they don't exist
func initDB(ctx context.Context, db *sql.DB) error {
	const createTablesSQL = `
		CREATE TABLE IF NOT EXISTS priority_locks (
			name TEXT PRIMARY KEY,
			owner TEXT NOT NULL,
			priority INTEGER NOT NULL,
			expires_at TIMESTAMPTZ NOT NULL
		);

		CREATE TABLE IF NOT EXISTS lock_queue (
			lock_name TEXT NOT NULL,
			owner TEXT NOT NULL,
			priority INTEGER NOT NULL,
			requested_at TIMESTAMPTZ NOT NULL,
			heartbeat_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (lock_name, owner)
		);

		CREATE INDEX IF NOT EXISTS idx_queue_priority 
		ON lock_queue(lock_name, priority DESC, requested_at ASC);
	`
	_, err := db.ExecContext(ctx, createTablesSQL)
	return err
}

func worker(ctx context.Context, db *sql.DB, workerID int, priority int, lockName string) {
	owner := fmt.Sprintf("worker-%d-p%d", workerID, priority)

	log.Printf("[Worker %d | Priority %d] 🚀 Starting, attempting to acquire lock '%s'", workerID, priority, lockName)

	// Try to acquire lock with blocking retry
	startTime := time.Now()
	result, err := Lock(ctx, db, LockParams{
		Name:             lockName,
		Owner:            owner,
		Priority:         priority,
		TTLSeconds:       4,
		IntervalDuration: 300 * time.Millisecond,
	})
	if err != nil {
		log.Printf("[Worker %d | Priority %d] ❌ Failed to acquire lock: %v", workerID, priority, err)
		return
	}
	waitTime := time.Since(startTime)

	log.Printf("[Worker %d | Priority %d] ✅ Acquired lock after waiting %v! Expires at %s",
		workerID, priority, waitTime.Round(100*time.Millisecond), result.ExpiresAt.Format("15:04:05"))

	// Simulate work while holding the lock
	workDuration := 2 * time.Second
	log.Printf("[Worker %d | Priority %d] 🔨 Working for %v...", workerID, priority, workDuration)

	time.Sleep(workDuration)

	log.Printf("[Worker %d | Priority %d] ✅ Work completed!", workerID, priority)

	// Release the lock
	unlockResult, err := Unlock(ctx, db, UnlockParams{
		Name:  lockName,
		Owner: owner,
	})
	if err != nil {
		log.Printf("[Worker %d | Priority %d] ⚠️  Error releasing lock: %v", workerID, priority, err)
		return
	}
	if unlockResult.Released {
		log.Printf("[Worker %d | Priority %d] 🔓 Lock released successfully", workerID, priority)
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
	db.SetMaxOpenConns(15)
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
	_, _ = db.ExecContext(ctx, "DELETE FROM priority_locks")
	_, _ = db.ExecContext(ctx, "DELETE FROM lock_queue")
	log.Println("🧹 Cleaned up existing locks")
	log.Println("")

	// Demo: Workers with different priorities competing for the same lock
	lockName := "priority-resource"

	log.Println("╔═══════════════════════════════════════════════════════════╗")
	log.Println("║  Priority-Based Distributed Lock Demo                    ║")
	log.Println("║  Higher priority workers should acquire lock first       ║")
	log.Println("╚═══════════════════════════════════════════════════════════╝")
	log.Println("")

	// Launch workers with different priorities
	// Worker priorities: [1=low, 5=medium, 10=high]
	workers := []struct {
		id       int
		priority int
		delay    time.Duration // Stagger start time
	}{
		{1, 1, 0 * time.Millisecond},    // Low priority, starts first
		{2, 5, 100 * time.Millisecond},  // Medium priority
		{3, 10, 200 * time.Millisecond}, // High priority
		{4, 1, 300 * time.Millisecond},  // Low priority
		{5, 10, 400 * time.Millisecond}, // High priority
		{6, 5, 500 * time.Millisecond},  // Medium priority
	}

	var wg sync.WaitGroup
	for _, w := range workers {
		wg.Add(1)
		go func(id, priority int, delay time.Duration) {
			defer wg.Done()
			time.Sleep(delay)
			worker(ctx, db, id, priority, lockName)
		}(w.id, w.priority, w.delay)
	}

	// Wait for all workers to complete
	wg.Wait()

	log.Println("\n═══════════════════════════════════════════════════════════")
	log.Println("✅ All workers completed!")
	log.Println("")
	log.Println("📌 Expected order: Worker 1 (first), then high priority")
	log.Println("   workers (3, 5) should go before low priority (2, 4, 6)")
}
