package repositories

import (
	"codelit/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetAllMembers(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewDBRepository(db)

	columns := []string{"id", "name", "type", "role", "duration", "tags"}
	rows := sqlmock.NewRows(columns).
		AddRow(1, "John Doe", "employee", "Software Engineer", 5, pq.Array([]string{"tag1", "tag2"})).
		AddRow(2, "Jane Smith", "employee", "Project Manager", 7, pq.Array([]string{"tag3", "tag4"}))

	// Act
	mock.ExpectQuery("SELECT id, name, type, role, duration, tags FROM members").WillReturnRows(rows)

	members, err := repo.GetAllMembers()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, members, 2)
	assert.Equal(t, 1, members[0].ID)
	assert.Equal(t, "John Doe", members[0].Name)
	assert.Equal(t, []string{"tag1", "tag2"}, members[0].Tags)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMemberByID(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewDBRepository(db)

	columns := []string{"id", "name", "type", "role", "duration", "tags"}
	row := sqlmock.NewRows(columns).
		AddRow(1, "John Doe", "employee", "Software Engineer", 5, pq.Array([]string{"tag1", "tag2"}))

	mock.ExpectQuery("SELECT \\* FROM members WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(row)

	// Act
	member, err := repo.GetMemberByID(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, 1, member.ID)
	assert.Equal(t, "John Doe", member.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateMember(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewDBRepository(db)

	query := "INSERT INTO members \\(name, type, role, duration, tags\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\) RETURNING id"
	mock.ExpectQuery(query).
		WithArgs("John Doe", "employee", "Software Engineer", 5, pq.Array([]string{"tag1", "tag2"})).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	member := &models.Member{
		Name:     "John Doe",
		Type:     "employee",
		Role:     "Software Engineer",
		Duration: 5,
		Tags:     []string{"tag1", "tag2"},
	}

	// Act
	err := repo.CreateMember(member)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, member.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateMember(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewDBRepository(db)

	query := "UPDATE members SET name = \\$1, type = \\$2, role = \\$3, duration = \\$4, tags = \\$5 WHERE id = \\$6"
	mock.ExpectExec(query).
		WithArgs("John Doe", "employee", "Software Engineer", 5, pq.Array([]string{"tag1", "tag2"}), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	member := &models.Member{
		ID:       1,
		Name:     "John Doe",
		Type:     "employee",
		Role:     "Software Engineer",
		Duration: 5,
		Tags:     []string{"tag1", "tag2"},
	}

	// Act
	err := repo.UpdateMember(member)

	// Assert the results
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteMember(t *testing.T) {
	// Arrange
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewDBRepository(db)

	query := "DELETE FROM members WHERE id = \\$1"
	mock.ExpectExec(query).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Act
	err := repo.DeleteMember(1)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
