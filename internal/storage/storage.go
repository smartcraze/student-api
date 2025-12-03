package storage

import (
	"context"
	"time"
)

type Student struct {
	ID             int64     `db:"id"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	RegistrationNo int       `db:"registration_no"`
	PhoneNumber    int64     `db:"phone_number"`
	Email          string    `db:"email"`
	Password       string    `db:"password"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type Storage interface {
	CreateStudent(ctx context.Context, student *Student) error
	GetStudentByID(ctx context.Context, id int64) (*Student, error)
	GetStudentByEmail(ctx context.Context, email string) (*Student, error)
	UpdateStudent(ctx context.Context, student *Student) error
	DeleteStudent(ctx context.Context, id int64) error
	ListStudents(ctx context.Context, limit, offset int) ([]*Student, error)
}
