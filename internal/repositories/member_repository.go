package repositories

import (
	grpcclient "codelit/internal/client"
	"codelit/internal/client/pb"
	"codelit/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/lib/pq"
)

type MemberRepository interface {
	GetAllMembers() ([]*models.Member, error)
	GetMemberByID(id int) (*models.Member, error)
	CreateMember(member *models.Member) error
	UpdateMember(member *models.Member) error
	DeleteMember(id int) error
}

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) *DBRepository {
	return &DBRepository{
		db: db,
	}
}

func (r *DBRepository) GetAllMembers() ([]*models.Member, error) {
	rows, err := r.db.Query("SELECT id, name, type, role, duration, tags FROM members")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []*models.Member{}
	for rows.Next() {
		member := &models.Member{}
		var tags pq.StringArray // Use pq.StringArray to store tags as an array of strings
		err := rows.Scan(&member.ID, &member.Name, &member.Type, &member.Role, &member.Duration, &tags)
		if err != nil {
			return nil, err
		}
		member.Tags = []string(tags) // Convert pq.StringArray to []string
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *DBRepository) GetMemberByID(id int) (*models.Member, error) {
	member := &models.Member{}

	var tags pq.StringArray // Use pq.StringArray to store tags as an array of strings

	row := r.db.QueryRow("SELECT * FROM members WHERE id = $1", id)
	err := row.Scan(&member.ID, &member.Name, &member.Type, &member.Role, &member.Duration, &tags)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("member not found")
		}
		return nil, err
	}

	member.Tags = []string(tags) // Convert pq.StringArray to []string

	return member, nil
}

func (r *DBRepository) CreateMember(member *models.Member) error {
	query := `INSERT INTO members (name, type, role, duration, tags)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	tagsArray := pq.Array(member.Tags) // Convert slice of strings to pq.Array
	err := r.db.QueryRow(query, member.Name, member.Type, member.Role, member.Duration, tagsArray).Scan(&member.ID)
	if err != nil {
		return err
	}
	conn := grpcclient.StartGRPC()
	client := pb.NewGreeterClient(conn)

	req := &pb.HelloRequest{
		Name: member.Name,
	}

	res, err := client.SayHello(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
	//TODO: Add validate member logic
	return nil
}

func (r *DBRepository) UpdateMember(member *models.Member) error {
	query := `UPDATE members SET name = $1, type = $2, role = $3, duration = $4, tags = $5
	WHERE id = $6`
	tagsArray := pq.Array(member.Tags) // Convert slice of strings to pq.Array
	_, err := r.db.Exec(query, member.Name, member.Type, member.Role, member.Duration, tagsArray, member.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBRepository) DeleteMember(id int) error {
	query := "DELETE FROM members WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
