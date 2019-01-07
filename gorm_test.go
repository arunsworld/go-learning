package learning

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	sqlite3 "github.com/mattn/go-sqlite3"
)

func TestGormQuickStart(t *testing.T) {
	type User struct {
		ID        uint   `gorm:"primary_key"`
		Email     string `gorm:"unique;not null"`
		Password  string
		FirstName string
		LastName  string
		IsActive  bool
	}

	db, err := gorm.Open("sqlite3", "/tmp/test.db")
	if err != nil {
		t.Error("Could not open DB:", err)
	}
	defer db.Close()
	db.LogMode(false)

	defer func() {
		err = os.Remove("/tmp/test.db")
		if err != nil {
			t.Error("Could not delete the database: ", err)
			return
		}
	}()

	db.AutoMigrate(&User{})

	// Test User Creation
	u := User{Email: "arunsworld@gmail.com", FirstName: "Arun", IsActive: true}
	errors := db.Create(&u).GetErrors()
	if len(errors) > 0 {
		t.Error("Encountered errors while creating user.", errors)
		return
	}
	if u.ID != 1 {
		t.Error("New user created. Expected ID=1. Got:", u.ID)
		return
	}

	// Test User Search - one user
	var user User
	db.First(&user, "email = ?", "arunsworld@gmail.com")
	if user.ID != 1 {
		t.Error("New user created. Expected ID=1. Got:", user.ID)
		return
	}

	// Update user to deactivate
	db.Model(&user).Update("IsActive", false)
	if user.IsActive {
		t.Error("Deactivated the user but it's still active...")
		return
	}

	// Test User Creation - Second Time
	u2 := User{Email: "arunsworld@gmail.com", FirstName: "Arun", IsActive: true}
	errors = db.Create(&u2).GetErrors()
	if len(errors) != 1 {
		t.Error("Expected 1 unique constraint error. Got: ", errors)
		return
	}
	e, ok := errors[0].(sqlite3.Error)
	if !ok {
		t.Error("Error expected to be sqlite3.Error but didn't get that.")
		return
	}
	if e.Code.Error() != "constraint failed" {
		t.Error("Expected constraint failed but got:", e.Code)
		return
	}

	// Test Successful User Creation - Second Time
	err = db.Create(&User{Email: "arun@e2open.com", FirstName: "Arun", IsActive: true}).Error
	if err != nil {
		t.Error("Encountered errors while creating second user.", err)
		return
	}

	// Test user search multiple users
	users := []User{}
	db.Find(&users)
	if len(users) != 2 {
		t.Error("Expected 2 user records but did not get that. Got:", len(users))
		return
	}

	// Basic count directly on table
	var count int
	db.Table("users").Count(&count)
	if count != 2 {
		t.Error("Expected 2 records during basic count but did not get that. Got:", count)
		return
	}

	// Record not found
	uu := User{}
	notFound := db.First(&uu, "20").RecordNotFound()
	if !notFound {
		t.Error("Did not expect to find record with index 20.")
		return
	}

	// Select using IN
	moreUsers := []User{}
	db.Find(&moreUsers, "id IN (?)", []int{1, 2})
	if count != 2 {
		t.Error("Expected 2 records during select using IN. Got:", count)
		return
	}

}
