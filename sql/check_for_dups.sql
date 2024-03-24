WITH
meta_album_rank as (
    SELECT
        ma.id,
        COUNT(*) as rank
    FROM meta_albums as ma
    INNER JOIN spotify_albums as sa ON (sa.meta_album_id = ma.id)
    GROUP BY 1
),
dup_meta_albums as (
    SELECT *
    FROM meta_album_rank as ma
    WHERE rank > 1
),
dup_meta_tracks as (
    SELECT
        mt.id,
        mt.simplified_name
    FROM meta_tracks as mt
    INNER JOIN dup_meta_albums as ma ON (ma.id = mt.meta_album_id)
    ORDER BY 2, 1
),
dup_spotify_tracks as (
    select
        a.name,
        a.meta_track_id,
        group_concat(a.spotify_id) as spotify_ids,
        COUNT(a.spotify_id) as rank
    FROM spotify_tracks as a
    -- INNER JOIN spotify_tracks as b ON (a.id < b.id AND a.meta_track_id = b.meta_track_id)
    GROUP BY 1, 2
    ORDER BY rank DESC, 1, 2
)
SELECT * FROM dup_spotify_tracks