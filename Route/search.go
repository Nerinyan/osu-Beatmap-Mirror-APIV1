package Route

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nerina1241/osu-beatmap-mirror-api/osu"
	"github.com/nerina1241/osu-beatmap-mirror-api/src"
)

func parseSort(s string) (ss string) { //sort

	s = strings.ToLower(s)
	switch s {
	case "ranked_asc":
		ss += "ranked_date asc"
	case "favourites_asc":
		ss += "favourite_count asc"
	case "favourites_desc":
		ss += "favourite_count desc"
	case "plays_asc":
		ss += "play_count asc"
	case "plays_desc":
		ss += "play_count desc"
	case "updated_asc":
		ss += "last_updated asc"
	case "updated_desc":
		ss += "last_updated desc"
	default:
		ss += "ranked_date desc"
	}

	return
}

func parsePage(s string) (ss string) {
	atoi, err := strconv.Atoi(s)
	if err != nil || atoi <= 0 {
		return " limit 50 "
	}
	return fmt.Sprintf("limit %d,50", atoi*50)
}

func parseMode(s string) (ss string) {
	s = strings.ToLower(s)
	switch s {
	case "0":
		ss = "0"
	case "1":
		ss = "1"
	case "2":
		ss = "2"
	case "3":
		ss = "3"
	default:
		ss = "0,1,2,3"
	}
	return
}

func parseStatus(s string) (ss string) {
	switch s {
	case "ranked":
		ss = "1,2"
	case "qualified":
		ss = "3"
	case "loved":
		ss = "4"
	case "pending":
		ss = "0"
	case "wip":
		ss = "-1"
	case "graveyard":
		ss = "-2"
	case "any":
		ss = "4,3,2,1,0,-1,-2"
	default:
		ss = "4,2,1"

	}
	return
}

func Search(c echo.Context) (err error) {
	var q string
	var rows *sql.Rows
	if c.QueryParam("q") == "" {
		q = fmt.Sprintf(src.QuerySearchBeatmapSet,
			parseStatus(c.QueryParam("s")), //ranked
			parseStatus(c.QueryParam("s")), //ranked
			parseMode(c.QueryParam("m")),   //osu,mania
			parseSort(c.QueryParam("sort")),
			parsePage(c.QueryParam("p")), //page

		)
		rows, err = src.Maria.Query(q)
	} else {
		q = fmt.Sprintf(src.QuerySearchBeatmapSetWhitQueryText,
			parseStatus(c.QueryParam("s")), //ranked
			parseStatus(c.QueryParam("s")), //ranked
			parseMode(c.QueryParam("m")),   //osu,mania
			parseSort(c.QueryParam("sort")),
			parsePage(c.QueryParam("p")), //page

		)
		rows, err = src.Maria.Query(q, c.QueryParam("q"))
	}

	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var sets []osu.BeatmapSetsOUT
	var index = map[int]int{}
	var mapids []int
	for rows.Next() {
		var set osu.BeatmapSetsOUT

		err = rows.Scan(
			// beatmapset_id, artist, artist_unicode, creator, favourite_count, hype_current,
			//hype_required, nsfw, play_count, source, status, title, title_unicode, user_id,
			//video, availability_download_disabled, availability_more_information, bpm, can_be_hyped,
			//discussion_enabled, discussion_locked, is_scoreable, last_updated, legacy_thread_url,
			//nominations_summary_current, nominations_summary_required, ranked, ranked_date, storyboard,
			//submitted_date, tags, has_favourited, description, genre_id, genre_name, language_id, language_name, ratings
			&set.Id, &set.Artist, &set.ArtistUnicode, &set.Creator, &set.FavouriteCount, &set.Hype.Current,
			&set.Hype.Required, &set.Nsfw, &set.PlayCount, &set.Source, &set.Status, &set.Title, &set.TitleUnicode, &set.UserId,
			&set.Video, &set.Availability.DownloadDisabled, &set.Availability.MoreInformation, &set.Bpm, &set.CanBeHyped,
			&set.DiscussionEnabled, &set.DiscussionLocked, &set.IsScoreable, &set.LastUpdated, &set.LegacyThreadUrl,
			&set.NominationsSummary.Current, &set.NominationsSummary.Required, &set.Ranked, &set.RankedDate, &set.Storyboard,
			&set.SubmittedDate, &set.Tags, &set.HasFavourited, &set.Description.Description, &set.Genre.Id, &set.Genre.Name,
			&set.Language.Id, &set.Language.Name, &set.RatingsString)
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return
		}
		index[*set.Id] = len(sets)
		mapids = append(mapids, *set.Id)
		sets = append(sets, set)
	}

	if len(sets) < 1 {
		c.NoContent(http.StatusNotFound)
		return
	}
	st := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(mapids)), ", "), "[]")
	rows, err = src.Maria.Query(fmt.Sprintf(src.QueryBeatmap, st))
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var Map osu.BeatmapOUT
		err = rows.Scan(
			//beatmap_id, beatmapset_id, mode, mode_int, status, ranked, total_length, max_combo, difficulty_rating,
			//version, accuracy, ar, cs, drain, bpm, convert, count_circles, count_sliders, count_spinners, deleted_at,
			//hit_length, is_scoreable, last_updated, passcount, playcount, checksum, user_id
			&Map.Id, &Map.BeatmapsetId, &Map.Mode, &Map.ModeInt, &Map.Status, &Map.Ranked, &Map.TotalLength, &Map.MaxCombo, &Map.DifficultyRating,
			&Map.Version, &Map.Accuracy, &Map.Ar, &Map.Cs, &Map.Drain, &Map.Bpm, &Map.Convert, &Map.CountCircles, &Map.CountSliders, &Map.CountSpinners, &Map.DeletedAt,
			&Map.HitLength, &Map.IsScoreable, &Map.LastUpdated, &Map.Passcount, &Map.Playcount, &Map.Checksum, &Map.UserId,
		)
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return
		}
		sets[index[*Map.BeatmapsetId]].Beatmaps = append(sets[index[*Map.BeatmapsetId]].Beatmaps, Map)

	}

	return c.JSON(http.StatusOK, sets)
}
