## Removing bad data

This is by no means recommended.. But during testing I somehow scrobbled movies.

    SET FOREIGN_KEY_CHECKS=0;
    DELETE FROM artists WHERE `name` = "%!s(<nil>)";
    DELETE FROM albums WHERE `name` = "%!s(<nil>)";
    DELETE album_artist FROM album_artist LEFT JOIN artists ON artists.uuid = album_artist.artist WHERE artists.uuid is null;
    DELETE album_artist FROM album_artist LEFT JOIN albums ON albums.uuid = album_artist.album WHERE albums.uuid is null;
    DELETE track_artist FROM track_artist LEFT JOIN artists ON artists.uuid = track_artist.artist WHERE artists.uuid is null;
    DELETE tracks FROM tracks LEFT JOIN track_artist ON track_artist.track = tracks.uuid WHERE track_artist.track IS NULL;
    DELETE scrobbles FROM scrobbles LEFT JOIN tracks ON tracks.uuid = scrobbles.track WHERE tracks.uuid is null;
    SET FOREIGN_KEY_CHECKS=1;
