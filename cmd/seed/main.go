package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/dasler-fw/bookcrossing/internal/models"
)

const (
	batchSize     = 1000
	softDeletePct = 5 // percent
)

// Defaults for seeding; can be overridden via env vars SEED_*
var (
	genresTotal    = getEnvInt("SEED_GENRES", 0)
	usersTotal     = getEnvInt("SEED_USERS", 1000000)
	booksTotal     = getEnvInt("SEED_BOOKS", 0)
	reviewsTotal   = getEnvInt("SEED_REVIEWS", 0)
	exchangesTotal = getEnvInt("SEED_EXCHANGES", 0)
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db := mustConnect()

	// Silence logger for bulk ops
	db = db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})

	fmt.Println("Clearing existing data...")
	truncateAll(db)
	fmt.Println("Data cleared ✓")

	gofakeit.Seed(0)

	genreIDs := seedGenres(db)
	userIDs := seedUsers(db)
	bookIDs, bookOwners := seedBooks(db, userIDs)
	seedBookGenres(db, bookIDs, genreIDs)
	seedReviews(db, userIDs, bookIDs)
	seedExchanges(db, bookIDs, bookOwners)

	fmt.Println("\n=== Seeding completed ===")
	fmt.Printf("Genres:    %d\n", len(genreIDs))
	fmt.Printf("Users:     %d\n", len(userIDs))
	fmt.Printf("Books:     %d\n", len(bookIDs))
	fmt.Printf("Reviews:   %d\n", reviewsTotal)
	fmt.Printf("Exchanges: %d\n", exchangesTotal)
}

func mustConnect() *gorm.DB {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
		if err != nil {
			log.Fatalf("failed to connect using DATABASE_URL: %v", err)
		}
		return db
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	if pass == "" {
		pass = os.Getenv("DB_PASS")
	}
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, dbname, sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return db
}

func truncateAll(db *gorm.DB) {
	// Most dependent -> least dependent
	// exchanges, reviews, book_genres, books, genres, users
	// Use CASCADE to handle FKs and restart identities
	stmt := "TRUNCATE TABLE exchanges, reviews, book_genres, books, genres, users RESTART IDENTITY CASCADE"
	db.Exec(stmt)
}

func seedGenres(db *gorm.DB) []uint {
	total := genresTotal
	ids := make([]uint, 0, total)
	buf := make([]models.Genre, 0, batchSize)

	fmt.Printf("Seeding genres... 0/%d", total)
	for i := 0; i < total; i++ {
		g := models.Genre{
			Name: fmt.Sprintf("Genre %03d %s", i+1, gofakeit.BookGenre()),
		}
		// ~5% soft delete
		if gofakeit.Number(1, 100) <= softDeletePct {
			deletedAt := gofakeit.DateRange(time.Now().AddDate(-2, 0, 0), time.Now().AddDate(-1, 0, 0))
			g.DeletedAt = gorm.DeletedAt{Time: deletedAt, Valid: true}
		}
		buf = append(buf, g)
		if len(buf) >= batchSize {
			insertAndCollect(db, &buf, &ids)
			fmt.Printf("\rSeeding genres... %d/%d", i+1, total)
		}
	}
	if len(buf) > 0 {
		insertAndCollect(db, &buf, &ids)
	}
	fmt.Println(" ✓")
	return ids
}

func seedUsers(db *gorm.DB) []uint {
	total := usersTotal
	ids := make([]uint, 0, total)
	buf := make([]models.User, 0, batchSize)

	fmt.Printf("Seeding users... 0/%d", total)
	for i := 0; i < total; i++ {
		u := models.User{
			Name:         gofakeit.Name(),
			Email:        fmt.Sprintf("user_%06d@test.com", i),
			PasswordHash: gofakeit.UUID(),
			City:         gofakeit.City(),
			Address:      gofakeit.Address().Address,
		}
		if gofakeit.Number(1, 100) <= softDeletePct {
			deletedAt := gofakeit.DateRange(time.Now().AddDate(-2, 0, 0), time.Now().AddDate(-1, 0, 0))
			u.DeletedAt = gorm.DeletedAt{Time: deletedAt, Valid: true}
		}
		buf = append(buf, u)
		if len(buf) >= batchSize {
			insertAndCollect(db, &buf, &ids)
			fmt.Printf("\rSeeding users... %d/%d", i+1, total)
		}
	}
	if len(buf) > 0 {
		insertAndCollect(db, &buf, &ids)
	}
	fmt.Println(" ✓")
	return ids
}

func seedBooks(db *gorm.DB, userIDs []uint) ([]uint, []uint) {
	total := booksTotal
	bookIDs := make([]uint, 0, total)
	bookOwners := make([]uint, 0, total)
	buf := make([]models.Book, 0, batchSize)

	statuses := []string{"available", "reserved"}

	fmt.Printf("Seeding books... 0/%d", total)
	for i := 0; i < total; i++ {
		owner := userIDs[gofakeit.Number(0, len(userIDs)-1)]
		b := models.Book{
			Title:       gofakeit.BookTitle(),
			Author:      gofakeit.Name(),
			Description: gofakeit.Paragraph(1, 3, 12, " "),
			AISummary:   gofakeit.Sentence(12),
			Status:      statuses[gofakeit.Number(0, len(statuses)-1)],
			UserID:      owner,
		}
		if gofakeit.Number(1, 100) <= softDeletePct {
			deletedAt := gofakeit.DateRange(time.Now().AddDate(-2, 0, 0), time.Now().AddDate(-1, 0, 0))
			b.DeletedAt = gorm.DeletedAt{Time: deletedAt, Valid: true}
		}
		buf = append(buf, b)
		if len(buf) >= batchSize {
			// insert batch
			sess := db.Session(&gorm.Session{SkipDefaultTransaction: true, SkipHooks: true})
			if err := sess.CreateInBatches(&buf, batchSize).Error; err != nil {
				log.Fatalf("failed to insert books batch: %v", err)
			}
			for _, rec := range buf {
				bookIDs = append(bookIDs, rec.ID)
				bookOwners = append(bookOwners, rec.UserID)
			}
			buf = buf[:0]
			fmt.Printf("\rSeeding books... %d/%d", i+1, total)
		}
	}
	if len(buf) > 0 {
		sess := db.Session(&gorm.Session{SkipDefaultTransaction: true, SkipHooks: true})
		if err := sess.CreateInBatches(&buf, batchSize).Error; err != nil {
			log.Fatalf("failed to insert books batch: %v", err)
		}
		for _, rec := range buf {
			bookIDs = append(bookIDs, rec.ID)
			bookOwners = append(bookOwners, rec.UserID)
		}
	}
	fmt.Println(" ✓")
	return bookIDs, bookOwners
}

func seedBookGenres(db *gorm.DB, bookIDs, genreIDs []uint) {
	fmt.Printf("Seeding book genres (many-to-many)... 0/%d", len(bookIDs))

	// We'll insert into join table directly for speed
	// Build batched VALUES list
	insertPrefix := "INSERT INTO book_genres (book_id, genre_id) VALUES "
	pairs := make([]string, 0, batchSize*3) // avg 3 genres per book
	countBooks := 0
	for i, bookID := range bookIDs {
		// choose 1-5 unique genres for this book
		n := gofakeit.Number(1, 5)
		if n > len(genreIDs) {
			n = len(genreIDs)
		}
		picked := make(map[uint]struct{}, n)
		for len(picked) < n {
			gid := genreIDs[gofakeit.Number(0, len(genreIDs)-1)]
			if _, ok := picked[gid]; ok {
				continue
			}
			picked[gid] = struct{}{}
			pairs = append(pairs, fmt.Sprintf("(%d,%d)", bookID, gid))
		}

		countBooks++
		// Flush periodically by books or when pairs large
		if countBooks%batchSize == 0 || len(pairs) >= batchSize*5 {
			flushJoinPairs(db, insertPrefix, &pairs)
			fmt.Printf("\rSeeding book genres (many-to-many)... %d/%d", i+1, len(bookIDs))
		}
	}
	if len(pairs) > 0 {
		flushJoinPairs(db, insertPrefix, &pairs)
	}
	fmt.Println(" ✓")
}

func flushJoinPairs(db *gorm.DB, insertPrefix string, pairs *[]string) {
	if len(*pairs) == 0 {
		return
	}
	// Postgres limit for single statement is large, but keep manageable
	query := insertPrefix + strings.Join(*pairs, ",")
	if err := db.Exec(query).Error; err != nil {
		log.Fatalf("failed to insert book_genres batch: %v", err)
	}
	*pairs = (*pairs)[:0]
}

func seedReviews(db *gorm.DB, userIDs, bookIDs []uint) {
	total := reviewsTotal
	buf := make([]models.Review, 0, batchSize)

	fmt.Printf("Seeding reviews... 0/%d", total)
	for i := 0; i < total; i++ {
		author := userIDs[gofakeit.Number(0, len(userIDs)-1)]
		targetUser := userIDs[gofakeit.Number(0, len(userIDs)-1)]
		// Avoid self-target occasionally
		if targetUser == author && len(userIDs) > 1 {
			if idx := gofakeit.Number(0, len(userIDs)-2); uint(idx) < uint(len(userIDs)-1) {
				targetUser = userIDs[idx]
			} else {
				targetUser = userIDs[len(userIDs)-1]
			}
		}
		targetBook := bookIDs[gofakeit.Number(0, len(bookIDs)-1)]
		r := models.Review{
			AuthorID:     author,
			TargetUserID: targetUser,
			TargetBookID: targetBook,
			Text:         gofakeit.Paragraph(1, 2, 10, " "),
			Rating:       gofakeit.Number(1, 5),
		}
		if gofakeit.Number(1, 100) <= softDeletePct {
			deletedAt := gofakeit.DateRange(time.Now().AddDate(-2, 0, 0), time.Now().AddDate(-1, 0, 0))
			r.DeletedAt = gorm.DeletedAt{Time: deletedAt, Valid: true}
		}
		buf = append(buf, r)
		if len(buf) >= batchSize {
			sess := db.Session(&gorm.Session{SkipDefaultTransaction: true, SkipHooks: true})
			if err := sess.CreateInBatches(&buf, batchSize).Error; err != nil {
				log.Fatalf("failed to insert reviews batch: %v", err)
			}
			buf = buf[:0]
			fmt.Printf("\rSeeding reviews... %d/%d", i+1, total)
		}
	}
	if len(buf) > 0 {
		sess := db.Session(&gorm.Session{SkipDefaultTransaction: true, SkipHooks: true})
		if err := sess.CreateInBatches(&buf, batchSize).Error; err != nil {
			log.Fatalf("failed to insert reviews batch: %v", err)
		}
	}
	fmt.Println(" ✓")
}

func seedExchanges(db *gorm.DB, bookIDs, bookOwners []uint) {
	total := exchangesTotal
	buf := make([]models.Exchange, 0, batchSize)

	statuses := []string{"pending", "accepted", "completed", "cancelled"}

	fmt.Printf("Seeding exchanges... 0/%d", total)
	for i := 0; i < total; i++ {
		// pick two distinct books
		idx1 := gofakeit.Number(0, len(bookIDs)-1)
		idx2 := gofakeit.Number(0, len(bookIDs)-1)
		for idx2 == idx1 || bookOwners[idx2] == bookOwners[idx1] {
			idx2 = gofakeit.Number(0, len(bookIDs)-1)
		}
		status := statuses[gofakeit.Number(0, len(statuses)-1)]
		var completedAt *time.Time
		if status == "completed" {
			ct := gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now())
			completedAt = &ct
		}
		e := models.Exchange{
			InitiatorID:     bookOwners[idx1],
			RecipientID:     bookOwners[idx2],
			InitiatorBookID: bookIDs[idx1],
			RecipientBookID: bookIDs[idx2],
			Status:          status,
			CompletedAt:     completedAt,
		}
		if gofakeit.Number(1, 100) <= softDeletePct {
			deletedAt := gofakeit.DateRange(time.Now().AddDate(-2, 0, 0), time.Now().AddDate(-1, 0, 0))
			e.DeletedAt = gorm.DeletedAt{Time: deletedAt, Valid: true}
		}
		buf = append(buf, e)
		if len(buf) >= batchSize {
			sess := db.Session(&gorm.Session{SkipDefaultTransaction: true, SkipHooks: true})
			if err := sess.CreateInBatches(&buf, batchSize).Error; err != nil {
				log.Fatalf("failed to insert exchanges batch: %v", err)
			}
			buf = buf[:0]
			fmt.Printf("\rSeeding exchanges... %d/%d", i+1, total)
		}
	}
	if len(buf) > 0 {
		sess := db.Session(&gorm.Session{SkipDefaultTransaction: true, SkipHooks: true})
		if err := sess.CreateInBatches(&buf, batchSize).Error; err != nil {
			log.Fatalf("failed to insert exchanges batch: %v", err)
		}
	}
	fmt.Println(" ✓")
}

// insertAndCollect inserts a batch of records and collects their IDs into ids slice.
func insertAndCollect[T any](db *gorm.DB, batch *[]T, ids *[]uint) {
	sess := db.Session(&gorm.Session{SkipDefaultTransaction: true, SkipHooks: true})
	if err := sess.CreateInBatches(batch, batchSize).Error; err != nil {
		log.Fatalf("failed to insert batch: %v", err)
	}
	// Reflect over batch to read ID field using GORM's Model interface
	// Since we know our structs embed gorm.Model, we can query back IDs using rows if reflection is complex.
	// But GORM populates the IDs in the slice elements after Create.
	switch v := any(*batch).(type) {
	case []models.User:
		for _, r := range v {
			*ids = append(*ids, r.ID)
		}
	case []models.Genre:
		for _, r := range v {
			*ids = append(*ids, r.ID)
		}
	case []models.Book:
		for _, r := range v {
			*ids = append(*ids, r.ID)
		}
	default:
		// Fallback: try to scan IDs with RETURNING via LASTVAL is unsafe; do nothing
	}
	*batch = (*batch)[:0]
}

func getEnvInt(name string, def int) int {
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		var n int
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil && n >= 0 {
			return n
		}
	}
	return def
}
