WITH discogs_releases_rank as (
    SELECT
        dr.meta_album_id,
        COUNT(*) as rank
    FROM discogs_releases as dr
    WHERE dr.meta_album_id IS NOT NULL
    GROUP BY 1
)
SELECT
    drr.rank as rank,
    dr.artist_name,
    dr.name,
    ma.simplified_name,
    dr.discogs_id,
    dr.meta_album_id
FROM discogs_releases as dr
INNER JOIN discogs_releases_rank as drr ON (drr.meta_album_id = dr.meta_album_id)
LEFT JOIN meta_albums as ma ON (ma.id = dr.meta_album_id)
ORDER BY 1 DESC, 2, 3
