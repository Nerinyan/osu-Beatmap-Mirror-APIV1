package src

import (
	"fmt"
	"strings"

	"github.com/nerina1241/osu-beatmap-mirror-api/ConsoleLogger"
	"github.com/nerina1241/osu-beatmap-mirror-api/osu"
)

func upsertMap(m osu.BeatmapIN, ch chan struct{}) {

	Upsert(UpsertBeatmap, []interface{}{
		m.Id, m.BeatmapsetId, m.Mode, m.ModeInt, m.Status, m.Ranked, m.TotalLength, m.MaxCombo, m.DifficultyRating, m.Version,
		m.Accuracy, m.Ar, m.Cs, m.Drain, m.Bpm, m.Convert, m.CountCircles, m.CountSliders, m.CountSpinners, m.DeletedAt,
		m.HitLength, m.IsScoreable, m.LastUpdated, m.Passcount, m.Playcount, m.Checksum, m.UserId,
	})
	ch <- struct{}{}
}

const (
	setUpsert = `
		INSERT INTO BeatmapMirror.beatmapset (
			beatmapset_id,artist,artist_unicode,creator,favourite_count,
			nsfw,play_count,source,
			status,title,title_unicode,user_id,video,
			availability_download_disabled,availability_more_information,bpm,can_be_hyped,discussion_enabled,
			discussion_locked,is_scoreable,last_updated,legacy_thread_url,nominations_summary_current,
			nominations_summary_required,ranked,ranked_date,storyboard,submitted_date,
			tags,has_favourited )
		VALUES %s ON DUPLICATE KEY UPDATE 
			artist = VALUES(artist), artist_unicode = VALUES(artist_unicode), creator = VALUES(creator), favourite_count = VALUES(favourite_count), 
			nsfw = VALUES(nsfw), play_count = VALUES(play_count), source = VALUES(source), 
			status = VALUES(status), title = VALUES(title), title_unicode = VALUES(title_unicode), user_id = VALUES(user_id), video = VALUES(video), 
			availability_download_disabled = VALUES(availability_download_disabled), availability_more_information = VALUES(availability_more_information), 
			bpm = VALUES(bpm), can_be_hyped = VALUES(can_be_hyped), discussion_enabled = VALUES(discussion_enabled), 
			discussion_locked = VALUES(discussion_locked), is_scoreable = VALUES(is_scoreable), last_updated = VALUES(last_updated), 
			legacy_thread_url = VALUES(legacy_thread_url), nominations_summary_current = VALUES(nominations_summary_current), 
			nominations_summary_required = VALUES(nominations_summary_required), ranked = VALUES(ranked), ranked_date = VALUES(ranked_date), 
			storyboard = VALUES(storyboard), submitted_date = VALUES(submitted_date), 
			tags = VALUES(tags), has_favourited = VALUES(has_favourited);`
	setValues = `(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)` //30

	mapUpsert = `
		INSERT INTO BeatmapMirror.beatmap (	
			beatmap_id,beatmapset_id,mode,mode_int,status,	ranked,total_length,max_combo,difficulty_rating,version,
			accuracy,ar,cs,drain,bpm,` + "`convert`" + `,count_circles,count_sliders,count_spinners,deleted_at,
			hit_length,is_scoreable,last_updated,passcount,playcount,	checksum,user_id
		)VALUES %s ON DUPLICATE KEY UPDATE 
			beatmapset_id = VALUES(beatmapset_id), mode = VALUES(mode), mode_int = VALUES(mode_int), status = VALUES(status), 
			ranked = VALUES(ranked), total_length = VALUES(total_length), max_combo = VALUES(max_combo), 
			difficulty_rating = VALUES(difficulty_rating), version = VALUES(version), 
			accuracy = VALUES(accuracy), ar = VALUES(ar), cs = VALUES(cs), drain = VALUES(drain), bpm = VALUES(bpm), 
			` + "`convert`" + ` = VALUES(` + "`convert`" + `), count_circles = VALUES(count_circles), count_sliders = VALUES(count_sliders),
			count_spinners = VALUES(count_spinners), deleted_at = VALUES(deleted_at), 
			hit_length = VALUES(hit_length), is_scoreable = VALUES(is_scoreable), last_updated = VALUES(last_updated), 
			passcount = VALUES(passcount), playcount = VALUES(playcount), 
			checksum = VALUES(checksum), user_id = VALUES(user_id);`
	mapValues         = `(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)` //27
	selectDeletedMaps = `select beatmap_id from BeatmapMirror.beatmap where beatmapset_id in (%s) AND beatmap_id not in (%s)`
	deleteMap         = `delete from BeatmapMirror.beatmap where beatmap_id in (%s);`
)

func buildSqlValues(s string, count int) (r string) {
	var sbuf []string
	for i := 0; i < count; i++ {
		sbuf = append(sbuf, s)
	}
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(sbuf)), ","), "[]")
}

func updateMapset(s *osu.BeatmapSetsIN) {
	//	beatmapset_id,artist,artist_unicode,creator,favourite_count,
	//	hype_current,hype_required,nsfw,play_count,source,
	//	status,title,title_unicode,user_id,video,
	//	availability_download_disabled,availability_more_information,bpm,can_be_hyped,discussion_enabled,
	//	discussion_locked,is_scoreable,last_updated,legacy_thread_url,nominations_summary_current,
	//	nominations_summary_required,ranked,ranked_date,storyboard,submitted_date,
	//	tags,has_favourited,description,genre_id,genre_name,
	//	language_id,language_name,ratings

	r := *s.Ratings
	Upsert(UpsertBeatmapSet, []interface{}{
		s.Id, s.Artist, s.ArtistUnicode, s.Creator, s.FavouriteCount,
		s.Hype.Current, s.Hype.Required, s.Nsfw, s.PlayCount, s.Source,
		s.Status, s.Title, s.TitleUnicode, s.UserId, s.Video,
		s.Availability.DownloadDisabled, s.Availability.MoreInformation, s.Bpm, s.CanBeHyped, s.DiscussionEnabled,
		s.DiscussionLocked, s.IsScoreable, s.LastUpdated, s.LegacyThreadUrl, s.NominationsSummary.Current,
		s.NominationsSummary.Required, s.Ranked, s.RankedDate, s.Storyboard, s.SubmittedDate,
		s.Tags, s.HasFavourited, s.Description.Description, s.Genre.Id, s.Genre.Name,
		s.Language.Id, s.Language.Name,
		fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d", r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7], r[8], r[9], r[10]),
	})

	if *s.Beatmaps == nil {
		return
	}
	ch := make(chan struct{}, len(*s.Beatmaps))
	for _, m := range *s.Beatmaps {
		go upsertMap(m, ch)
	}
	for i := 0; i < len(ch); i++ {
		<-ch
	}
}

func updateSearchBeatmaps(data *[]osu.BeatmapSetsIN) (err error) {
	if data == nil {
		return
	}
	if len(*data) < 1 {
		return
	}

	var (
		setInsertBuf []interface{}
		mapInsertBuf []interface{}
		beatmapSets  []int
		beatmaps     []int
		deletedMaps  []int
	)

	for _, s := range *data {
		beatmapSets = append(beatmapSets, s.Id)
		setInsertBuf = append(setInsertBuf,
			s.Id, s.Artist, s.ArtistUnicode, s.Creator, s.FavouriteCount, s.Nsfw, s.PlayCount, s.Source,
			s.Status, s.Title, s.TitleUnicode, s.UserId, s.Video,
			s.Availability.DownloadDisabled, s.Availability.MoreInformation, s.Bpm, s.CanBeHyped, s.DiscussionEnabled,
			s.DiscussionLocked, s.IsScoreable, s.LastUpdated, s.LegacyThreadUrl, s.NominationsSummary.Current,
			s.NominationsSummary.Required, s.Ranked, s.RankedDate, s.Storyboard, s.SubmittedDate,
			s.Tags, s.HasFavourited,
		)
		for _, m := range *s.Beatmaps {
			beatmaps = append(beatmaps, m.Id)
			mapInsertBuf = append(mapInsertBuf,
				m.Id, m.BeatmapsetId, m.Mode, m.ModeInt, m.Status, m.Ranked, m.TotalLength, m.MaxCombo, m.DifficultyRating, m.Version,
				m.Accuracy, m.Ar, m.Cs, m.Drain, m.Bpm, m.Convert, m.CountCircles, m.CountSliders, m.CountSpinners, m.DeletedAt,
				m.HitLength, m.IsScoreable, m.LastUpdated, m.Passcount, m.Playcount, m.Checksum, m.UserId,
			)
		}
	}
	//맵셋
	if _, err = Maria.Exec(fmt.Sprintf(setUpsert, buildSqlValues(setValues, len(beatmapSets))), setInsertBuf...); err != nil {
		fmt.Println(err)
		return err
	}

	//맵
	if _, err = Maria.Exec(fmt.Sprintf(mapUpsert, buildSqlValues(mapValues, len(beatmaps))), mapInsertBuf...); err != nil {
		fmt.Println(err)
		return err
	}

	//삭제된 맵 제거
	sets := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(beatmapSets)), ","), "[]")
	maps := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(beatmaps)), ","), "[]")
	rows, err := Maria.Query(fmt.Sprintf(selectDeletedMaps, sets, maps))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var i int
		if err = rows.Scan(&i); err != nil {
			return err
		}
		deletedMaps = append(deletedMaps, i)
	}
	if len(deletedMaps) > 1 {
		ConsoleLogger.DangersConsolelog("DELETED MAPS", strings.Join(strings.Fields(fmt.Sprint(deletedMaps)), ","))
		dmaps := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(deletedMaps)), ","), "[]")
		// fmt.Println(fmt.Sprintf(deleteMap, dmaps))
		if _, err = Maria.Exec(fmt.Sprintf(deleteMap, dmaps)); err != nil {
			ConsoleLogger.WarningConsolelog("ERROR", err.Error())
			return err
		}
	}
	return
}
