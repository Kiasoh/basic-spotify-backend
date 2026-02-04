package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiasoh/basic-spotify-backend/models"
)

type UserRepository interface {
	CreateUserInTx(ctx context.Context, tx pgx.Tx, user *models.User) (int, error)
	UpdateRecommPlaylistIDInTx(ctx context.Context, tx pgx.Tx, userID int, playlistID int) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context) ([]models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUserInTx(ctx context.Context, tx pgx.Tx, user *models.User) (int, error) {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	var id int
	err := tx.QueryRow(ctx, query, user.Username, user.Password.Bytes()).Scan(&id)
	return id, err
}

func (r *userRepository) UpdateRecommPlaylistIDInTx(ctx context.Context, tx pgx.Tx, userID int, playlistID int) error {
	query := `UPDATE users SET recomm_plylist_id = $1 WHERE id = $2`
	_, err := tx.Exec(ctx, query, playlistID, userID)
	return err
}


func (r *userRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, username, password, recomm_plylist_id, created_at FROM users WHERE id = $1`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Password.Hash, &user.RecommPlaylistID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, password, recomm_plylist_id, created_at FROM users WHERE username = $1`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password.Hash, &user.RecommPlaylistID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET username = $1, password = $2, recomm_plylist_id = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, user.Username, user.Password.Bytes(), user.RecommPlaylistID, user.ID)
	return err
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *userRepository) ListUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, username, created_at FROM users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
