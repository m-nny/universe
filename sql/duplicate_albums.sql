-- SQLite
SELECT
    A.Name as a_name,
    B.Name as b_name,
    A.simplified_name as a_simplified_name,
    B.simplified_name as b_simplified_name,
    A.spotify_id as a_spotify_id,
    B.spotify_id as b_spotify_id
FROM spotify_albums as A
INNER JOIN spotify_albums AS B
ON (
    a.id < b.id AND
    B.name LIKE (A.name || '%') AND
    substr(a.simplified_name, 0, 4) = substr(b.simplified_name, 0, 4)
    -- TRUE
)
WHERE NOT (
    B.Name LIKE '%Remix%' OR
    -- B.Name LIKE '%. (Deluxe Edition)' OR
    -- B.Name LIKE '% (10th Anniversary Edition)' OR
    -- B.Name LIKE '% Reconfigured' OR
    FALSE
)
ORDER BY A.simplified_name, A.Name, A.spotify_id, B.spotify_id
