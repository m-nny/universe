SELECT ma.*, a.*
		FROM meta_albums as ma
		LEFT JOIN meta_album_artists as maa ON (ma.simplified_name = maa.meta_album_id)
		LEFT JOIN artists as a ON (maa.artist_id = a.id)
		WHERE ma.simplified_name IN (?)
