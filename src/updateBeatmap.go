package src

import (
	"fmt"

	"github.com/nerina1241/osu-beatmap-mirror-api/osu"
)

func upsertMap(m osu.Beatmap, ch chan struct{}) {

	Upsert(UpsertBeatmap, []interface{}{
		m.Id, m.BeatmapsetId, m.Mode, m.ModeInt, m.Status, m.Ranked, m.TotalLength, m.MaxCombo, m.DifficultyRating, m.Version,
		m.Accuracy, m.Ar, m.Cs, m.Drain, m.Bpm, m.Convert, m.CountCircles, m.CountSliders, m.CountSpinners, m.DeletedAt,
		m.HitLength, m.IsScoreable, m.LastUpdated, m.Passcount, m.Playcount, m.Checksum, m.UserId,
	})
	ch <- struct{}{}
}

func updateMapset(s *osu.BeatmapSets) {
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

func updateSearchBeatmaps(data *[]osu.BeatmapSets) (err error) {
	if *data == nil {
		return
	}
	for _, s := range *data {
		Upsert(UpsertBeatmapSet2, []interface{}{
			s.Id, s.Artist, s.ArtistUnicode, s.Creator, s.FavouriteCount,
			s.Nsfw, s.PlayCount, s.Source,
			s.Status, s.Title, s.TitleUnicode, s.UserId, s.Video,
			s.Availability.DownloadDisabled, s.Availability.MoreInformation, s.Bpm, s.CanBeHyped, s.DiscussionEnabled,
			s.DiscussionLocked, s.IsScoreable, s.LastUpdated, s.LegacyThreadUrl, s.NominationsSummary.Current,
			s.NominationsSummary.Required, s.Ranked, s.RankedDate, s.Storyboard, s.SubmittedDate,
			s.Tags, s.HasFavourited,
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
	return
}
