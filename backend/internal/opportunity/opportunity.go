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

const (
	Scholarship   Type = "scholarship"
	Hackathon     Type = "hackathon"
	Internship    Type = "internship"
	ResearchExtra Type = "research_extra"
)

func (t Type) Valid() bool {
	return t == Scholarship || t == Hackathon || t == Internship || t == ResearchExtra
}

type Opportunity struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        Type      `json:"type"`
	Eligibility string    `json:"eligibility"`
	Steps       string    `json:"steps"`
	Benefits    string    `json:"benefits"`
	Link        string    `json:"link"`
	Referral    string    `json:"referral"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type Input struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        Type   `json:"type"`
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
	if !i.Type.Valid() {
		return errors.New("type must be scholarship, hackathon, internship, or research_extra")
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

const columns = "id,title,description,type,eligibility,steps,benefits,link,referral,created_by,created_at,updated_at"

func scan(row interface{ Scan(...any) error }) (Opportunity, error) {
	var o Opportunity
	err := row.Scan(&o.ID, &o.Title, &o.Description, &o.Type, &o.Eligibility, &o.Steps, &o.Benefits, &o.Link, &o.Referral, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}
func (s *Store) List(ctx context.Context, filter Type, page, perPage int) (Page, error) {
	p := Page{Items: []Opportunity{}, Page: page, PerPage: perPage}
	where, args := "", []any{}
	if filter != "" {
		where = " WHERE type=$1"
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
func (s *Store) Get(ctx context.Context, id uuid.UUID) (Opportunity, error) {
	return scan(s.db.QueryRow(ctx, "SELECT "+columns+" FROM opportunities WHERE id=$1", id))
}
func (s *Store) Create(ctx context.Context, in Input, userID uuid.UUID) (Opportunity, error) {
	return scan(s.db.QueryRow(ctx, "INSERT INTO opportunities (title,description,type,eligibility,steps,benefits,link,referral,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "+columns, in.Title, in.Description, in.Type, in.Eligibility, in.Steps, in.Benefits, in.Link, in.Referral, userID))
}
func (s *Store) Update(ctx context.Context, id, userID uuid.UUID, in Input) (Opportunity, error) {
	return scan(s.db.QueryRow(ctx, "UPDATE opportunities SET title=$1,description=$2,type=$3,eligibility=$4,steps=$5,benefits=$6,link=$7,referral=$8,updated_at=now() WHERE id=$9 AND created_by=$10 RETURNING "+columns, in.Title, in.Description, in.Type, in.Eligibility, in.Steps, in.Benefits, in.Link, in.Referral, id, userID))
}
