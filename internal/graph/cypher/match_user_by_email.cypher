MATCH
    (u:User {email: $email})
RETURN
    u.id AS id,
    u.email AS email,
    u.first_name AS first_name,
    u.last_name AS last_name,
    u.password AS password;
