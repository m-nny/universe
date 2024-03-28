WITH saved_tracks as (
    SELECT
        user_username,
        meta_track_id
    FROM user_saved_tracks
    WHERE user_username = "m-nny"
),
saved_albums as (
    SELECT
        DISTINCT mt.meta_album_id as id
    FROM saved_tracks as st
    INNER JOIN meta_tracks as mt ON (st.meta_track_id = mt.id)
),
releases as (
    SELECT *
    FROM discogs_releases as dr
    INNER JOIN saved_albums as sa ON (sa.id = dr.meta_album_id)
)
-- SELECT
--     ma.*
-- FROM saved_albums as sa
-- INNER JOIN meta_albums as ma ON (sa.id = ma.id)
SELECT * FROM releases
ORDER BY 4, 3
;
