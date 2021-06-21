package src

const UpsertBeatmap = `
INSERT INTO BeatmapMirror.beatmap
(
	beatmap_id,beatmapset_id,mode,mode_int,status,	ranked,total_length,max_combo,difficulty_rating,version,
	accuracy,ar,cs,drain,bpm,` + "`convert`" + `,count_circles,count_sliders,count_spinners,deleted_at,
	hit_length,is_scoreable,last_updated,passcount,playcount,	checksum,user_id
)VALUES(
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?
)ON DUPLICATE KEY UPDATE 
	beatmapset_id = ?, mode = ?, mode_int = ?, status = ?, 
	ranked = ?, total_length = ?, max_combo = ?, difficulty_rating = ?, version = ?, 
	accuracy = ?, ar = ?, cs = ?, drain = ?, bpm = ?, 
	` + "`convert`" + ` = ?, count_circles = ?, count_sliders = ?, count_spinners = ?, deleted_at = ?, 
	hit_length = ?, is_scoreable = ?, last_updated = ?, passcount = ?, playcount = ?, 
	checksum = ?, user_id = ?
;
`

const UpsertBeatmapSet = `
INSERT INTO BeatmapMirror.beatmapset(
	beatmapset_id,artist,artist_unicode,creator,favourite_count,
	hype_current,hype_required,nsfw,play_count,source,
	status,title,title_unicode,user_id,video,
	availability_download_disabled,availability_more_information,bpm,can_be_hyped,discussion_enabled,
	discussion_locked,is_scoreable,last_updated,legacy_thread_url,nominations_summary_current,
	nominations_summary_required,ranked,ranked_date,storyboard,submitted_date,
	tags,has_favourited,description,genre_id,genre_name,
	language_id,language_name,ratings
)VALUES(
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?
)ON DUPLICATE KEY UPDATE 
artist= ?, artist_unicode= ?, creator= ?, favourite_count= ?, 
	hype_current= ?, hype_required= ?, nsfw= ?, play_count= ?, source= ?, 
	status= ?, title= ?, title_unicode= ?, user_id= ?, video= ?, 
	availability_download_disabled= ?, availability_more_information= ?, bpm= ?, can_be_hyped= ?, discussion_enabled= ?, 
	discussion_locked= ?, is_scoreable= ?, last_updated= ?, legacy_thread_url= ?, nominations_summary_current= ?, 
	nominations_summary_required= ?, ranked= ?, ranked_date= ?, storyboard= ?, submitted_date= ?, 
	tags= ?, has_favourited= ?, description= ?, genre_id= ?, genre_name= ?, 
	language_id= ?, language_name= ?, ratings= ?
;
`

const UpsertBeatmapSet2 = `
INSERT INTO BeatmapMirror.beatmapset(
	beatmapset_id,artist,artist_unicode,creator,favourite_count,
	nsfw,play_count,source,
	status,title,title_unicode,user_id,video,
	availability_download_disabled,availability_more_information,bpm,can_be_hyped,discussion_enabled,
	discussion_locked,is_scoreable,last_updated,legacy_thread_url,nominations_summary_current,
	nominations_summary_required,ranked,ranked_date,storyboard,submitted_date,
	tags,has_favourited
)VALUES(
	?,?,?,?,?,	
	?,?,?,
	?,?,?,?,?,	
	?,?,?,?,?,
	?,?,?,?,?,	
	?,?,?,?,?,
	?,?	
)ON DUPLICATE KEY UPDATE 
	artist= ?, artist_unicode= ?, creator= ?, favourite_count= ?, 
	nsfw= ?, play_count= ?, source= ?, 
	status= ?, title= ?, title_unicode= ?, user_id= ?, video= ?, 
	availability_download_disabled= ?, availability_more_information= ?, bpm= ?, can_be_hyped= ?, discussion_enabled= ?, 
	discussion_locked= ?, is_scoreable= ?, last_updated= ?, legacy_thread_url= ?, nominations_summary_current= ?, 
	nominations_summary_required= ?, ranked= ?, ranked_date= ?, storyboard= ?, submitted_date= ?, 
	tags= ?, has_favourited= ?
;
`

const QueryBeatmap = `select * from BeatmapMirror.beatmap where beatmapset_id in( %s );`

const QuerySearchBeatmapSet = `
select * from (select * from BeatmapMirror.beatmapset where ranked in( %s ) ) A 
inner join (select distinct beatmapset_id from BeatmapMirror.beatmap where ranked in( %s ) AND mode_int in ( %s ) ) B using (beatmapset_id)
order by A.%s %s ;
`
const QuerySearchBeatmapSetWhitQueryText = `
select A.* from (select * from BeatmapMirror.search_index where MATCH(text) AGAINST(?)) S
inner join (select * from BeatmapMirror.beatmapset where ranked in( %s ) ) A using (beatmapset_id)
inner join (select distinct beatmapset_id from BeatmapMirror.beatmap where ranked in( %s ) AND mode_int in ( %s ) ) B using (beatmapset_id) 
order by A.%s %s ;
`

const GetDownloadBeatmapData = `SELECT beatmapset_id,artist,title,last_updated FROM BeatmapMirror.beatmapset WHERE beatmapset_id = ?`

const QueryAPILog = `INSERT INTO BeatmapMirror.api_log (time, request_id, remote_ip, host, method, uri, user_agent, status, error, latency, latency_human, bytes_in, bytes_out) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?);`

/*
CREATE TABLE `BeatmapMirror.beatmap` (
  `beatmap_id` int(10) NOT NULL,
  `beatmapset_id` int(10) NOT NULL,
  `mode` varchar(6) COLLATE utf8_unicode_ci DEFAULT NULL,
  `mode_int` tinyint(1) DEFAULT NULL,
  `status` varchar(9) COLLATE utf8_unicode_ci DEFAULT NULL,
  `ranked` tinyint(1) DEFAULT NULL,
  `total_length` int(10) DEFAULT NULL,
  `max_combo` int(10) DEFAULT NULL,
  `difficulty_rating` decimal(63,2) DEFAULT NULL,
  `version` varchar(254) COLLATE utf8_unicode_ci DEFAULT NULL,
  `accuracy` decimal(63,2) DEFAULT NULL,
  `ar` decimal(63,2) DEFAULT NULL,
  `cs` decimal(63,2) DEFAULT NULL,
  `drain` decimal(63,2) DEFAULT NULL,
  `bpm` decimal(63,2) DEFAULT NULL,
  `convert` tinyint(1) DEFAULT NULL,
  `count_circles` int(10) DEFAULT NULL,
  `count_sliders` int(10) DEFAULT NULL,
  `count_spinners` int(10) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `hit_length` int(10) DEFAULT NULL,
  `is_scoreable` tinyint(1) DEFAULT NULL,
  `last_updated` datetime DEFAULT NULL,
  `passcount` int(10) DEFAULT NULL,
  `playcount` int(10) DEFAULT NULL,
  `checksum` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_id` int(10) DEFAULT NULL,
  PRIMARY KEY (`beatmap_id`),
  UNIQUE KEY `beatmap_id_UNIQUE` (`beatmap_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE `BeatmapMirror.beatmapset` (
  `beatmapset_id` int(1) NOT NULL,
  `artist` varchar(254) COLLATE utf8_unicode_ci DEFAULT NULL,
  `artist_unicode` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `creator` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `favourite_count` int(1) DEFAULT NULL,
  `hype_current` int(1) DEFAULT NULL,
  `hype_required` int(1) DEFAULT NULL,
  `nsfw` tinyint(1) DEFAULT NULL,
  `play_count` int(1) DEFAULT NULL,
  `source` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `status` varchar(9) CHARACTER SET utf8mb4 DEFAULT NULL,
  `title` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `title_unicode` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `user_id` int(1) DEFAULT NULL,
  `video` tinyint(1) DEFAULT NULL,
  `availability_download_disabled` tinyint(1) DEFAULT NULL,
  `availability_more_information` text CHARACTER SET utf8mb4 DEFAULT NULL,
  `bpm` decimal(63,2) DEFAULT NULL,
  `can_be_hyped` tinyint(1) DEFAULT NULL,
  `discussion_enabled` tinyint(1) DEFAULT NULL,
  `discussion_locked` tinyint(1) DEFAULT NULL,
  `is_scoreable` tinyint(1) DEFAULT NULL,
  `last_updated` datetime DEFAULT NULL,
  `legacy_thread_url` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `nominations_summary_current` int(1) DEFAULT NULL,
  `nominations_summary_required` int(1) DEFAULT NULL,
  `ranked` tinyint(1) DEFAULT NULL,
  `ranked_date` datetime DEFAULT NULL,
  `storyboard` tinyint(1) DEFAULT NULL,
  `submitted_date` datetime DEFAULT NULL,
  `tags` text CHARACTER SET utf8mb4 DEFAULT NULL,
  `has_favourited` tinyint(1) DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 DEFAULT NULL,
  `genre_id` int(1) DEFAULT NULL,
  `genre_name` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `language_id` int(1) DEFAULT NULL,
  `language_name` varchar(254) CHARACTER SET utf8mb4 DEFAULT NULL,
  `ratings` varchar(254) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`beatmapset_id`),
  UNIQUE KEY `id_UNIQUE` (`beatmapset_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE `api_log` (
  `time` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `request_id` varchar(32) NOT NULL,
  `remote_ip` varchar(21) DEFAULT NULL,
  `host` varchar(128) DEFAULT NULL,
  `method` varchar(10) DEFAULT NULL,
  `uri` longtext DEFAULT NULL,
  `user_agent` longtext DEFAULT NULL,
  `status` int(10) unsigned DEFAULT NULL,
  `error` longtext DEFAULT NULL,
  `latency` bigint(20) unsigned DEFAULT NULL,
  `latency_human` varchar(16) DEFAULT NULL,
  `bytes_in` bigint(20) unsigned DEFAULT NULL,
  `bytes_out` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 트리거================
DELIMITER $$
CREATE DEFINER=`root`@`%` TRIGGER `osu`.`beatmapset_BEFORE_INSERT` BEFORE INSERT ON `beatmapset` FOR EACH ROW
BEGIN
	INSERT INTO `osu`.`search_index`(`beatmapset_id`,`text`)
    VALUES(
		NEW.beatmapset_id,
		concat_ws(' ',NEW.beatmapset_id,NEW.artist,NEW.creator,NEW.title)
    )ON DUPLICATE KEY UPDATE text = concat_ws(' ',NEW.beatmapset_id,NEW.artist,NEW.creator,NEW.title);
END$$
DELIMITER ;
-- 트리거================
DELIMITER $$
CREATE DEFINER=`root`@`%` TRIGGER `osu`.`beatmapset_BEFORE_UPDATE` BEFORE UPDATE ON `beatmapset` FOR EACH ROW
BEGIN
	INSERT INTO `osu`.`search_index`(`beatmapset_id`,`text`)
    VALUES(
		NEW.beatmapset_id,
		concat_ws(' ',NEW.beatmapset_id,NEW.artist,NEW.creator,NEW.title)
    )ON DUPLICATE KEY UPDATE text = concat_ws(' ',NEW.beatmapset_id,NEW.artist,NEW.creator,NEW.title);
END$$
DELIMITER ;
-- 트리거================
*/
