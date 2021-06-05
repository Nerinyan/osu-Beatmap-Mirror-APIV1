package Route

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/labstack/echo"
	"github.com/nerina1241/osu-beatmap-mirror-api/src"
)

type minMax struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}
type query struct {
	SortBy struct {
		Column string
		// 마지막 업데이트/랭크 날짜/ 등록 날짜/셋id/곡 길이/타이틀/아티스트/난이도/Rating/플카/좋아요 수
		Desc bool //1,0
	}
	CS               *minMax //and cs between 0 and 10
	AR               *minMax //and ar between 0 and 10
	OD               *minMax //and od between 0 and 10
	HP               *minMax //and hp between 0 and 10
	BPM              *minMax //and bpm between 0 and 10
	Length           *minMax //and total_length between 0 and 10
	DifficultyRating *minMax //and difficulty_rating between 0 and 10
	Mode             *int    // mode = 0
	Categories       int     // 랭크
	Genre            int     //장르
	Language         int     // 언어
	Extra            string  //스토리보드/비디오
	ExplicitContent  bool    //nsfw
	Search           string
	Limit            int
	Skip             int
}

/*
select * from
	(
		select
			S.id ,S.bpm as set_bpm, S.is_scoreable, S.ranked_date, S.last_updated, S.submitted_date, S.title, S.artist, S.creator, S.source, S.tags,
			S.language, S.favourite_count, S.play_count, S.user_id, S.genre, S.download_disabled, S.storyboard, S.video, S.ranked,

			M.id as map_id , M.ar, M.cs, M.od, M.hp, M.bpm, M.difficulty_rating, M.version, M.mode, M.total_length,
			M.playcount, M.count_spinners, M.count_sliders, M.hit_length, M.count_circles, M.checksum
		from ( select * from osu.beatmap_set ) as S inner join ( select * from osu.beatmaps ) as M  on M.set_id = S.id
	) A
where
 	concat_ws(' ',id, artist, source, title, creator,tags) like '%koi no uta%'
	AND genre = 2
	AND download_disabled = 0
	AND storyboard = 1
	AND video = 1
	AND language = 1
	AND ranked = 1
    AND	mode = 0
	AND ar between 0 AND 10
	AND od between 0 AND 10
	AND cs between 0 AND 10
	AND hp between 0 AND 10
	AND bpm between 0 AND 360
	AND total_length between 0 AND 500
	AND difficulty_rating between 0 AND 10
    order by ranked_date desc ;
*/
func beatmapQuery(q *query) (code int, data string) {
	code = 500
	qs := src.SearchBeatmaps
	if q.Search != "" {
		qs += `concat_ws(' ',id, artist, source, title, creator,tags) like concat('%','?','%')`
	}
	fmt.Println(qs)
	fmt.Println(q.Search)
	rows, err := src.Maria.Query(qs)
	if err != nil {
		fmt.Println(err)
		return
	}
	//cols, _ := rows.Columns()
	for rows.Next() {
		var y interface{}
		err := rows.Scan(&y)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(y)
	}
	return
}

func Search(c echo.Context) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(500)
	}
	var Query query
	err = json.Unmarshal(b, &Query)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(500)
	}

	return c.JSON(beatmapQuery(&Query))
}
