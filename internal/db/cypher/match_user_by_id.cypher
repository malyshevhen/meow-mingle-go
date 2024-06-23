MATCH
    (u:User {id: $id})
RETURN
    u.id AS id,
    u.email AS email,
    u.first_name AS first_name,
    u.last_name AS last_name;
