-- SQLite
SELECT
    A.Name,
    B.Name,
    A.simplified_name,
    B.simplified_name
FROM albums as A
INNER JOIN albums AS B
ON (
    a.id != b.id AND 
    B.name LIKE A.name || '%' AND
    substr(a.simplified_name, 0, 4) = substr(b.simplified_name, 0, 4)
)
WHERE NOT (
    B.Name LIKE '%Remix%' OR
    B.Name LIKE '%. (Deluxe Edition)' OR
    B.Name LIKE '% (10th Anniversary Edition)' OR
    B.Name LIKE '% Reconfigured'
)
ORDER BY A.simplified_name, A.Name