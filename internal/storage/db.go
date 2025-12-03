package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	storage := &PostgresStorage{db: db}

	// Auto-create tables if they don't exist
	if err := storage.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return storage, nil
}

func (s *PostgresStorage) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		registration_no INTEGER UNIQUE NOT NULL,
		phone_number BIGINT NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_students_email ON students(email);
	CREATE INDEX IF NOT EXISTS idx_students_registration_no ON students(registration_no);
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateStudent(ctx context.Context, student *Student) error {
	query := `
		INSERT INTO students (first_name, last_name, registration_no, phone_number, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	now := time.Now()
	err := s.db.QueryRowContext(
		ctx,
		query,
		student.FirstName,
		student.LastName,
		student.RegistrationNo,
		student.PhoneNumber,
		student.Email,
		student.Password,
		now,
		now,
	).Scan(&student.ID)

	if err != nil {
		return fmt.Errorf("failed to create student: %w", err)
	}

	student.CreatedAt = now
	student.UpdatedAt = now
	return nil
}

func (s *PostgresStorage) GetStudentByID(ctx context.Context, id int64) (*Student, error) {
	query := `
		SELECT id, first_name, last_name, registration_no, phone_number, email, password, created_at, updated_at
		FROM students
		WHERE id = $1
	`
	var student Student
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.RegistrationNo,
		&student.PhoneNumber,
		&student.Email,
		&student.Password,
		&student.CreatedAt,
		&student.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("student not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	return &student, nil
}

func (s *PostgresStorage) GetStudentByEmail(ctx context.Context, email string) (*Student, error) {
	query := `
		SELECT id, first_name, last_name, registration_no, phone_number, email, password, created_at, updated_at
		FROM students
		WHERE email = $1
	`
	var student Student
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.RegistrationNo,
		&student.PhoneNumber,
		&student.Email,
		&student.Password,
		&student.CreatedAt,
		&student.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("student not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	return &student, nil
}

func (s *PostgresStorage) UpdateStudent(ctx context.Context, student *Student) error {
	query := `
		UPDATE students
		SET first_name = $1, last_name = $2, registration_no = $3, phone_number = $4, email = $5, updated_at = $6
		WHERE id = $7
	`
	result, err := s.db.ExecContext(
		ctx,
		query,
		student.FirstName,
		student.LastName,
		student.RegistrationNo,
		student.PhoneNumber,
		student.Email,
		time.Now(),
		student.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update student: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}

func (s *PostgresStorage) DeleteStudent(ctx context.Context, id int64) error {
	query := `DELETE FROM students WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}

func (s *PostgresStorage) ListStudents(ctx context.Context, limit, offset int) ([]*Student, error) {
	query := `
		SELECT id, first_name, last_name, registration_no, phone_number, email, password, created_at, updated_at
		FROM students
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}
	defer rows.Close()

	var students []*Student
	for rows.Next() {
		var student Student
		err := rows.Scan(
			&student.ID,
			&student.FirstName,
			&student.LastName,
			&student.RegistrationNo,
			&student.PhoneNumber,
			&student.Email,
			&student.Password,
			&student.CreatedAt,
			&student.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student: %w", err)
		}
		students = append(students, &student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return students, nil
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
