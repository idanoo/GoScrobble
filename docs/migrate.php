<?php

// Temp script to migrate old data from MySQL to PostgreSQL

// PLEASE RUN ALL MIGRATIONS ON POSTGRES BEFORE RUNNING THIS.

echo PHP_EOL . "Loading connections...";

global $mysqli;
$mysqli = new mysqli('172.27.138.37', 'goscrobble', 'X9u7jdfy', 'goscrobble');
if ($mysqli->connect_errno) {
    die("Failed to connect to MySQL");
}

global $postgres;
$postgres = new PDO("pgsql:host=127.0.0.1;port=5432;dbname=goscrobble;user=goscrobble;password=supersecretdatabasepassword1");

function getArray($query): array
{
    global $mysqli;
    if (!$result = $mysqli->query($query)) {
        die($mysqli->error);
    }

    while ($row = $result->fetch_assoc()) {
        $data[] = $row;
    }
    return $data;
}

echo PHP_EOL . "Skipping schema_migrations (Already exists)";

echo PHP_EOL . "Migrating config";
// $config = getArray("SELECT * FROM config");
// $update = $postgres->prepare("UPDATE config SET value = ? WHERE key = ?");
// foreach ($config as $row) {
//     $update->execute([
//         $row['value'],
//         $row['key'],
//     ]);
// }

echo PHP_EOL . "Migrating users";
$users = getArray("SELECT
        BIN_TO_UUID(uuid, true) as uuid,
        created_at,
        inet_ntoa(conv(created_ip, 16, 10)) as created_ip,
        modified_at, 
        inet_ntoa(modified_ip) as modified_ip,
        username,
        password,
        email,
        verified,
        active,
        admin,
        `mod`,
        token,
        private,
        timezone
    FROM users;");

$update = $postgres->prepare("INSERT INTO users (uuid, created_at, created_ip, modified_at, modified_ip, username, password, email, verified, active, admin, mod, token, private, timezone) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)");
foreach ($users as $row) {
    echo PHP_EOL . $row['username'];
    $update->execute([
        $row['uuid'],
        $row['created_at'],
        $row['created_ip'],
        $row['modified_at'],
        $row['modified_ip'],
        $row['username'],
        $row['password'],
        $row['email'],
        $row['verified'],
        $row['active'],
        $row['admin'],
        $row['mod'],
        $row['token'],
        $row['private'],
        $row['timezone'],
    ]);
}

// echo PHP_EOL . "Migrating albums";
// echo PHP_EOL . "Migrating artists";
// echo PHP_EOL . "Migrating tracks";
// echo PHP_EOL . "Migrating genres";
// echo PHP_EOL . "Migrating links";
// echo PHP_EOL . "Migrating oauth_tokens";
// echo PHP_EOL . "Migrating refresh_tokens";
// echo PHP_EOL . "Migrating resettoken";
// echo PHP_EOL . "Migrating album_artist";
// echo PHP_EOL . "Migrating track_album";
// echo PHP_EOL . "Migrating track_artist";
// echo PHP_EOL . "Migrating scrobbles";
