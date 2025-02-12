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
	WHERE schemename = 'public'
	AND tablename = 'artists');`).Scan(&exists); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		_, err := db.Query(`
		CREATE TABLE artists (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		genre VARCHAR(100),
		country VARCHAR(100));`)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

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
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		artist_id INT NOT NULL,
		album VARCHAR(255),
		release_year INT,
		genre VARCHAR(100),
		lyrics TEXT NOT NULL,
		CONSTRAINT fk_aritst FOREIGN KEY (artist_id) REFERENCES artists (id) ON DELETE CASCADE);`)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	

	return &Storage{db: db}, nil
}

func (s *Storage) SaveSong(song models.Song) (int64, error) {
	const op = "storage.psql.SaveSong"

	stmt, err := s.db.Prepare(`
	INSERT INTO songs (title,artist,album,release_year,genre,lyrics)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id;`)
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var id int64
	err = stmt.QueryRow(song.Title,
		song.Artist, song.Album, song.ReleaseYear,
		song.Genre, song.Lyrics).Scan(&id)
	if err != nil {
		if postgeErr, ok := err.(*pq.Error); ok && postgeErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrSongExists)
		}

		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
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
		&song.Genre, &song.Lyrics)
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
			&song.Genre, &song.Lyrics)
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
	release_year = $5, genre = $6, lyrics = $7 
	WHERE id = $1;`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(id, updatedSong.Title,
		updatedSong.Artist, updatedSong.Album, updatedSong.ReleaseYear,
		updatedSong.Genre, updatedSong.Lyrics)
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
