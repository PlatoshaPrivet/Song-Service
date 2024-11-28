package controllers

import (
	"database/sql"
	"net/http"
	db "song-service/database"
	"song-service/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Get songs via filtration
// @Tags Search
// @Description Get a list of songs
// @Produce json
// @Param id query int false "Song ID"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param group query string false "Group filter"
// @Param song query string false "Song filter"
// @Param realiseDate query string false "Realise date filter"
// @Param text query string false "Text filter"
// @Param link query string false "Link filter"
// @Success 200 {array} models.Song
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /songs/search/filters [get]
func GetSongs(c *gin.Context) {
	id := c.Query("id")
	group := c.Query("group")
	title := c.Query("song")
	realiseDate := c.Query("realiseDate")
	text := c.Query("text")
	link := c.Query("link")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	//Need bigint for postgres
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 1 {
		offsetInt = 0
	}

	rows, err := db.DB.Query(`
        SELECT id, group_name, song_name, release_date, text, link 
        FROM songs 
        WHERE (group_name ILIKE $1 OR $1 = '') AND (song_name ILIKE $2 OR $2 = '') 
		AND (release_date ILIKE $3 OR $3 = '') AND (text ILIKE $4 OR $4 = '') AND (link ILIKE $5 OR $5 = '') AND (id ILIKE $6 OR $6 = '')
        LIMIT $7 OFFSET $8
    `, "%"+group+"%", "%"+title+"%", "%"+realiseDate+"%", "%"+text+"%", "%"+link+"%", "%"+id+"%", limitInt, offsetInt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	songs := []models.Song{}
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		songs = append(songs, song)
	}

	c.JSON(http.StatusOK, songs)
}

// @Summary Get song via cuplet
// @Tags Search
// @Description Get a song by a verse
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param id query int true "Song ID"
// @Success 200 {array} models.Song
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /songs/search/verse [get]
func GetVerse(c *gin.Context) {
	id := c.Query("id")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("page", "1")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 1 {
		offsetInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10
	}

	row := db.DB.QueryRow(`
        SELECT text 
        FROM songs 
        WHERE id = $1
    `, idInt)

	var song models.Song
	err = row.Scan(&song.Text)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	verses := strings.Split(song.Text, `\n`)
	totalVerses := len(verses)
	start := (offsetInt - 1) * limitInt
	end := start + limitInt

	if start >= totalVerses {
		c.JSON(http.StatusOK, gin.H{"verses": []string{}, "page": offsetInt, "pageSize": limitInt})
		return
	}

	if end > totalVerses {
		end = totalVerses
	}

	c.JSON(http.StatusOK, gin.H{
		"songID":     song.ID,
		"verses":     verses[start:end],
		"page":       offsetInt,
		"pageSize":   limitInt,
		"totalCount": totalVerses,
	})
}

// @Summary Add New Song
// @Tags Post
// @Description Add a new song
// @Accept json
// @Produce json
// @Param group query string true "Group filter"
// @Param song query string true "Song filter"
// @Success 200 {object} models.Song
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /songs [post]
func AddSong(c *gin.Context) {
	var newSong models.Song
	if err := c.ShouldBindJSON(&newSong); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid input")
		return
	}
	query := `INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := db.DB.QueryRow(query, newSong.Group, newSong.Song, "empty", "empty", "empty").Scan(&newSong.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to add song:"+err.Error())
		return
	}

	c.JSON(http.StatusOK, newSong)
}

// @Summary Update song
// @Tags Update
// @Description Update the song details by ID
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Param song body models.Required true "Updated Song Info"
// @Success 200 {object} models.Song
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	id := c.Param("id")
	var updateSong models.Required
	if err := c.ShouldBindJSON(&updateSong); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.DB.Exec(`
        UPDATE songs SET group_name = $1, song_name = $2, 
        release_date = $3, text = $4, link = $5 WHERE id = $6`,
		updateSong.Group, updateSong.Song, updateSong.ReleaseDate,
		updateSong.Text, updateSong.Link, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

// @Summary Delete song
// @Description Delete song by ID
// @Tags Delete
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Success 200 {object} models.Song
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	id := c.Param("id")

	_, err := db.DB.Exec(`DELETE FROM songs WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
