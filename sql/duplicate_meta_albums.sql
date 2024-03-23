-- SQLite
SELECT
    a.id,
    b.id,
    a.simplified_name,
    b.simplified_name
FROM meta_albums as a
INNER JOIN meta_albums as b
ON (
    -- a.simplified_name = b.simplified_name
    a.id < b.id AND
    B.simplified_name LIKE (A.simplified_name || '%')
)
ORDER BY a.simplified_name, a.id