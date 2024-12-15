package psql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/leyl1ne/rest-api-parser/pkg/models"
	"github.com/leyl1ne/rest-api-parser/pkg/storage"
	"github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "leyline"
	password = "pass"
	dbname   = "parserDB"
)

type Storage struct {
	db *sql.DB
}

func Connect() (*sql.DB, error) {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	conn, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func CloseConnection(db *sql.DB) {
	defer db.Close()
}

func New() (*Storage, error) {
	const op = "storage.postgresql.NewStorage"

	db, err := Connect()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var exists bool
	if err := db.QueryRow(`
	SELECT EXISTS 
	(SELECT FROM pg_tables 
	WHERE schemaname = 'public' 
	AND tablename = 'songs' );`).Scan(&exists); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		_, err := db.Query(`
		CREATE TABLE songs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		artist VARCHAR(255) NOT NULL,
		album VARCHAR(255),
		release_year YEAR,
		genre VARCHAR(100),
		duration INT,
		lyrics TEXT NOT NULL);`)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveSong(song models.Song) (int64, error) {
	const op = "storage.psql.SaveSong"

	stmt, err := s.db.Prepare(`
	INSERT INTO songs (title,artist,album,release_year,genre,duration,lyrics)
	VALUES ($2,$3,$4,$5,$6,$7,$8);`)
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(song.Title,
		song.Artist, song.Album, song.ReleaseYear,
		song.Genre, song.Duration, song.Lyrics)
	if err != nil {
		if postgeErr, ok := err.(*pq.Error); ok && postgeErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrSongExists)
		}

		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil

}

func (s *Storage) DeleteSong(id int) error {
	const op = "storage.psql.DeleteSong"

	stmt, err := s.db.Prepare("DELTE FROM songs WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rowsAffected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrSongNotFound)
	}

	return nil
}

func (s *Storage) GetSong(id int) (models.Song, error) {
	const op = "storage.psql.GetSong"

	stmt, err := s.db.Prepare("SELECT * FROM songs WHERE id = $1")
	if err != nil {
		return models.Song{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	song := models.Song{}

	err = stmt.QueryRow(id).Scan(&song.ID, &song.Title,
		&song.Artist, &song.Album, &song.ReleaseYear,
		&song.Genre, &song.Duration, &song.Lyrics)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Song{}, fmt.Errorf("%s: %w", op, storage.ErrSongNotFound)
	}
	if err != nil {
		return models.Song{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return song, nil
}

func (s *Storage) GetAllSong() ([]models.Song, error) {
	const op = "storage.psql.GetAllSong"
	var songs = make([]models.Song, 0)

	stmt, err := s.db.Prepare("SELECT * FROM songs;")
	if err != nil {
		return songs, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Query()
	if err != nil {
		return songs, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	for res.Next() {
		var song models.Song
		err = res.Scan(&song.ID, &song.Title,
			&song.Artist, &song.Album, &song.ReleaseYear,
			&song.Genre, &song.Duration, &song.Lyrics)
		if err != nil {
			return songs, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (s *Storage) UpdateSong(id int, updatedSong models.Song) error {
	const op = "storage.psql.UpdateSong"

	stmt, err := s.db.Prepare(`
	UPDATE songs SET 
	title = $2, artist = $3, album = $4, 
	release_year = $5, genre = $6, duration = $7, lyrics = $8 
	WHERE id = $1;`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(id, updatedSong.Title,
		updatedSong.Artist, updatedSong.Album, updatedSong.ReleaseYear,
		updatedSong.Genre, updatedSong.Duration, updatedSong.Lyrics)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: getting rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrSongNotFound)
	}

	return nil

}
