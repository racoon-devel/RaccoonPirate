package representation

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

const (
	moviesDirectory = "Фильмы и Сериалы"
	musicDirectory  = "Музыка"
	booksDirectory  = "Книги"
	otherDirectory  = "Другое"

	filmsDirectory    = "Фильмы"
	tvSeriesDirectory = "Сериалы"

	byAlphabetDirectory = "Алфавит"
	byYearDirectory     = "Год"
	byGenreDirectory    = "Жанры"
)

func mapMediaTypeToDir(mediaType media.ContentType) string {
	switch mediaType {
	case media.Movies:
		return moviesDirectory
	case media.Music:
		return musicDirectory
	case media.Books:
		return booksDirectory
	default:
		return otherDirectory
	}
}

func mapMovieTypeToDir(movieType model.MovieType) string {
	switch movieType {
	case model.MovieType_Movie:
		return filmsDirectory
	case model.MovieType_TvSeries:
		return tvSeriesDirectory
	default:
		return filmsDirectory
	}
}
