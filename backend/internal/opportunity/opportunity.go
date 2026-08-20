package opportunity

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Type string
type ApprovalStatus string

const (
	Pending  ApprovalStatus = "pending"
	Approved ApprovalStatus = "approved"
)

const (
	Scholarship Type = "scholarship"
	Hackathon   Type = "hackathon"
	Internship  Type = "internship"
	Research    Type = "research"
	Extras      Type = "extras"
)

func (t Type) Valid() bool {
	return t == Scholarship || t == Hackathon || t == Internship || t == Research || t == Extras
}

type Opportunity struct {
	ID             int64 `json:"id"`
	recordID       uuid.UUID
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Types          []Type         `json:"types"`
	Eligibility    string         `json:"eligibility"`
	Steps          string         `json:"steps"`
	Benefits       string         `json:"benefits"`
	Link           string         `json:"link"`
	Referral       string         `json:"referral"`
	CreatedBy      uuid.UUID      `json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	ApprovalStatus ApprovalStatus `json:"approval_status"`
}
type Input struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Types       []Type `json:"types"`
	Eligibility string `json:"eligibility"`
	Steps       string `json:"steps"`
	Benefits    string `json:"benefits"`
	Link        string `json:"link"`
	Referral    string `json:"referral"`
}

func (i Input) Validate() error {
	i.Title = strings.TrimSpace(i.Title)
	i.Description = strings.TrimSpace(i.Description)
	if i.Title == "" || len(i.Title) > 250 {
		return errors.New("title must be between 1 and 250 characters")
	}
	if i.Description == "" || len(i.Description) > 10000 {
		return errors.New("description must be between 1 and 10000 characters")
	}
	if len(i.Types) == 0 {
		return errors.New("types must contain at least one opportunity type")
	}
	if len(i.Types) > 5 {
		return errors.New("types cannot contain more than 5 opportunity types")
	}
	seen := make(map[Type]struct{}, len(i.Types))
	for _, opportunityType := range i.Types {
		if !opportunityType.Valid() {
			return errors.New("types must contain scholarship, hackathon, internship, research, or extras")
		}
		if _, exists := seen[opportunityType]; exists {
			return errors.New("types cannot contain duplicate opportunity types")
		}
		seen[opportunityType] = struct{}{}
	}
	if i.Link != "" {
		u, err := url.ParseRequestURI(i.Link)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return errors.New("link must be an absolute URL")
		}
	}
	return nil
}

type Page struct {
	Items   []Opportunity `json:"items"`
	Page    int           `json:"page"`
	PerPage int           `json:"per_page"`
	Total   int           `json:"total"`
}
type Store struct{ db *pgxpool.Pool }

func NewStore(db *pgxpool.Pool) *Store { return &Store{db: db} }

const columns = "display_id,id,title,description,types,eligibility,steps,benefits,link,referral,created_by,created_at,updated_at,approval_status"

func scan(row interface{ Scan(...any) error }) (Opportunity, error) {
	var o Opportunity
	var types []string
	err := row.Scan(&o.ID, &o.recordID, &o.Title, &o.Description, &types, &o.Eligibility, &o.Steps, &o.Benefits, &o.Link, &o.Referral, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt, &o.ApprovalStatus)
	o.Types = make([]Type, len(types))
	for index, value := range types {
		o.Types[index] = Type(value)
	}
	return o, err
}
func (s *Store) List(ctx context.Context, filter Type, status ApprovalStatus, page, perPage int) (Page, error) {
	p := Page{Items: []Opportunity{}, Page: page, PerPage: perPage}
	where, args := " WHERE approval_status=$1", []any{status}
	if filter != "" {
		where += " AND $2::opportunity_type = ANY(types)"
		args = append(args, filter)
	}
	if err := s.db.QueryRow(ctx, "SELECT count(*) FROM opportunities"+where, args...).Scan(&p.Total); err != nil {
		return p, err
	}
	args = append(args, perPage, (page-1)*perPage)
	rows, err := s.db.Query(ctx, "SELECT "+columns+" FROM opportunities"+where+fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)-1, len(args)), args...)
	if err != nil {
		return p, err
	}
	defer rows.Close()
	for rows.Next() {
		o, err := scan(rows)
		if err != nil {
			return p, err
		}
		p.Items = append(p.Items, o)
	}
	return p, rows.Err()
}
func (s *Store) Get(ctx context.Context, displayID int64, status ApprovalStatus) (Opportunity, error) {
	return scan(s.db.QueryRow(ctx, "SELECT "+columns+" FROM opportunities WHERE display_id=$1 AND approval_status=$2", displayID, status))
}
func (s *Store) Create(ctx context.Context, in Input, userID uuid.UUID) (Opportunity, error) {
	return scan(s.db.QueryRow(ctx, "INSERT INTO opportunities (title,description,types,eligibility,steps,benefits,link,referral,created_by) VALUES ($1,$2,$3::opportunity_type[],$4,$5,$6,$7,$8,$9) RETURNING "+columns, in.Title, in.Description, typeStrings(in.Types), in.Eligibility, in.Steps, in.Benefits, in.Link, in.Referral, userID))
}
func (s *Store) Update(ctx context.Context, displayID int64, userID uuid.UUID, in Input) (Opportunity, error) {
	return scan(s.db.QueryRow(ctx, "UPDATE opportunities SET title=$1,description=$2,types=$3::opportunity_type[],eligibility=$4,steps=$5,benefits=$6,link=$7,referral=$8,approval_status='pending',updated_at=now() WHERE display_id=$9 AND created_by=$10 RETURNING "+columns, in.Title, in.Description, typeStrings(in.Types), in.Eligibility, in.Steps, in.Benefits, in.Link, in.Referral, displayID, userID))
}
func (s *Store) Approve(ctx context.Context, displayID int64) (Opportunity, error) {
	return scan(s.db.QueryRow(ctx, "UPDATE opportunities SET approval_status='approved',updated_at=now() WHERE display_id=$1 AND approval_status='pending' RETURNING "+columns, displayID))
}
func (s *Store) Delete(ctx context.Context, displayID int64) (bool, error) {
	result, err := s.db.Exec(ctx, "DELETE FROM opportunities WHERE display_id=$1", displayID)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() == 1, nil
}

func typeStrings(types []Type) []string {
	values := make([]string, len(types))
	for index, value := range types {
		values[index] = string(value)
	}
	return values
}
