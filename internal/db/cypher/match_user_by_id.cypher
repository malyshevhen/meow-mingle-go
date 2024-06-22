MATCH
    (u:User)
WHERE
    ID(u)=$id
RETURN
    ID(u) AS id,
    u.email AS email,
    u.first_name AS first_name,
    u.last_name AS last_name;
