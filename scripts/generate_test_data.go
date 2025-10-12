package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/kiin21/go-rest/internal/config"
	"github.com/kiin21/go-rest/internal/initialize/db"
	"gorm.io/gorm"
)

type Starter struct {
	Domain        string `gorm:"column:domain"`
	Name          string `gorm:"column:name"`
	Email         string `gorm:"column:email"`
	Mobile        string `gorm:"column:mobile"`
	WorkPhone     string `gorm:"column:work_phone"`
	JobTitle      string `gorm:"column:job_title"`
	DepartmentID  *int64 `gorm:"column:department_id"`
	LineManagerID *int64 `gorm:"column:line_manager_id"`
}

func (Starter) TableName() string {
	return "starters"
}

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Init DB
	dbConn, err := db.InitDB(&cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Generate and insert 100k records
	totalRecords := 100000
	batchSize := 1000

	jobTitles := []string{
		"Software Engineer", "Senior Software Engineer", "Lead Software Engineer",
		"Product Manager", "Senior Product Manager", "Data Engineer", "Senior Data Engineer",
		"DevOps Engineer", "QA Engineer", "UI/UX Designer", "Technical Manager",
		"Business Analyst", "Project Manager", "Scrum Master", "System Administrator",
		"Frontend Developer", "Backend Developer", "Full Stack Developer",
		"Mobile Developer", "Game Developer", "AI Engineer", "Security Engineer",
	}

	firstNames := []string{"Nguyen", "Tran", "Le", "Pham", "Hoang", "Huynh", "Phan", "Vu", "Vo", "Dang", "Bui", "Do"}
	lastNames := []string{"Anh", "Minh", "Duc", "Tuan", "Hieu", "Long", "Nam", "Khoa", "Thang", "Hai", "Tien", "Duy"}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < totalRecords; i += batchSize {
		batch := make([]Starter, 0, batchSize)
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		for j := i; j < end; j++ {
			domain := fmt.Sprintf("user%d", j+13) // Start from 13 since 1-12 exist
			firstName := firstNames[rand.Intn(len(firstNames))]
			lastName := lastNames[rand.Intn(len(lastNames))]
			name := fmt.Sprintf("%s %s", firstName, lastName)
			email := fmt.Sprintf("%s@vng.com.vn", domain)
			mobile := generateVietnamesePhone()
			jobTitle := jobTitles[rand.Intn(len(jobTitles))]
			deptID := int64(rand.Intn(12) + 1) // 1-12
			var lineMgrID *int64
			if rand.Float32() < 0.3 { // 30% have line manager
				mgr := rand.Int63n(12) + 1
				lineMgrID = &mgr
			}

			starter := Starter{
				Domain:        domain,
				Name:          name,
				Email:         email,
				Mobile:        mobile,
				JobTitle:      jobTitle,
				DepartmentID:  &deptID,
				LineManagerID: lineMgrID,
			}
			batch = append(batch, starter)
		}

		// Insert batch
		if err := insertBatch(dbConn, batch); err != nil {
			log.Fatalf("Failed to insert batch starting at %d: %v", i, err)
		}

		fmt.Printf("Inserted %d-%d records\n", i+1, end)
	}

	fmt.Println("Successfully inserted 100,000 test records!")
}

func insertBatch(db *gorm.DB, batch []Starter) error {
	return db.CreateInBatches(batch, len(batch)).Error
}

func generateVietnamesePhone() string {
	// Generate Vietnamese mobile number: (+84) 09xxxxxxxx or similar
	prefixes := []string{"09", "08", "07", "05", "03"}
	prefix := prefixes[rand.Intn(len(prefixes))]
	number := rand.Intn(900000000) + 100000000 // 8 digits
	return fmt.Sprintf("(+84) %s%d", prefix, number)
}
