package src

var UpsertMaps = `
INSERT INTO BeatmapMirror.beatmaps (
		id, set_id, set_ranked, set_ranked_txt, ranked, ranked_txt, mode, 
		mode_txt, title, title_unicode, artist, artist_unicode, version, creator, creator_id, set_submitted_date, set_last_updated, 
		set_ranked_date, last_updated, favourite_count, set_playcount, difficulty_rating, set_bpm, bpm, ar, cs, hp, 
		od, max_combo, playcount, passcount, total_length, hit_length, count_circles, count_spinners, count_sliders, has_storyboard, 
		has_video, language_id, language_name, genre_id, genre_name, tags, beatmaps_count
	) VALUES (
		?,?,?,?,?,?,?,
		?,?,?,?,?,?,?,?,?,?,
		?,?,?,?,?,?,?,?,?,?,
		?,?,?,?,?,?,?,?,?,?,
		?,?,?,?,?,?,?
    ) ON DUPLICATE KEY UPDATE
		set_id=?,set_ranked=?, set_ranked_txt=?, ranked=?, ranked_txt=?, mode=?, 
		mode_txt=?, title=?, title_unicode=?, artist=?, artist_unicode=?, version=?, creator=?, creator_id=?, set_submitted_date=?, set_last_updated=?, 
		set_ranked_date=?, last_updated=?, favourite_count=?, set_playcount=?, difficulty_rating=?, set_bpm=?, bpm=?, ar=?, cs=?, hp=?, 
		od=?, max_combo=?, playcount=?, passcount=?, total_length=?, hit_length=?, count_circles=?, count_spinners=?, count_sliders=?, has_storyboard=?, 
		has_video=?, language_id=?, language_name=?, genre_id=?, genre_name=?, tags=?, beatmaps_count =?
`
var UpsertMapsSet = `
REPLACE INTO BeatmapMirror.sets (
        beatmapset_id, title, title_unicode, artist, artist_unicode, creator, submitted_date, 
        ranked, ranked_date, last_updated, play_count, bpm, tags, genre_id,
        genre_name, language_id, language_name, favourite_count
    ) VALUES (
        ?,?,?,?,?,?,?,
        ?,?,?,?,?,?,?,
        ?,?,?,?
    )ON DUPLICATE KEY UPDATE
		title=?, title_unicode=?, artist=?, artist_unicode=?, creator=?, submitted_date=?, 
        ranked=?, ranked_date=?, last_updated=?, play_count=?, bpm=?, tags=?, genre_id=?,
        genre_name=?, language_id=?, language_name=?, favourite_count=?
`

var UpsertMaps2 = `
INSERT INTO BeatmapMirror.beatmaps (
		id, set_id, set_ranked, set_ranked_txt, ranked, ranked_txt, mode, 
		mode_txt, title, title_unicode, artist, artist_unicode, version, creator, creator_id, set_submitted_date, set_last_updated, 
		set_ranked_date, last_updated, favourite_count, set_playcount, difficulty_rating, set_bpm, bpm, ar, cs, hp, 
		od, max_combo, playcount, passcount, total_length, hit_length, count_circles, count_spinners, count_sliders, has_storyboard, 
		has_video, tags, beatmaps_count
	) VALUES (
		?,?,?,?,?,?,?,
		?,?,?,?,?,?,?,?,?,?,
		?,?,?,?,?,?,?,?,?,?,
		?,?,?,?,?,?,?,?,?,?,
		?,?,?
    ) ON DUPLICATE KEY UPDATE
		set_id=?,set_ranked=?, set_ranked_txt=?, ranked=?, ranked_txt=?, mode=?, 
		mode_txt=?, title=?, title_unicode=?, artist=?, artist_unicode=?, version=?, creator=?, creator_id=?, set_submitted_date=?, set_last_updated=?, 
		set_ranked_date=?, last_updated=?, favourite_count=?, set_playcount=?, difficulty_rating=?, set_bpm=?, bpm=?, ar=?, cs=?, hp=?, 
		od=?, max_combo=?, playcount=?, passcount=?, total_length=?, hit_length=?, count_circles=?, count_spinners=?, count_sliders=?, has_storyboard=?, 
		has_video=?, tags=?, beatmaps_count =?
`
var UpsertMapsSet2 = `
INSERT INTO BeatmapMirror.sets (
        beatmapset_id, title, title_unicode, artist, artist_unicode, creator, submitted_date, 
        ranked, ranked_date, last_updated, play_count, bpm, tags, favourite_count
    ) VALUES (
        ?,?,?,?,?,?,?,
        ?,?,?,?,?,?,?
    )ON DUPLICATE KEY UPDATE
		title=?, title_unicode=?, artist=?, artist_unicode=?, creator=?, submitted_date=?, 
        ranked=?, ranked_date=?, last_updated=?, play_count=?, bpm=?, tags=?, favourite_count=?
`

var CheckDownloadable = `SELECT download_disabled FROM osu.beatmap_set WHERE id = ?`
var GetDownloadBeatmapData = `SELECT beatmapset_id as id,artist,title,last_updated from BeatmapMirror.sets where beatmapset_id = ?`

var SearchBeatmaps = `
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
`
