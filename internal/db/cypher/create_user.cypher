CREATE
    (u:User {
        email: $email,
        first_name: $first_name,
        last_name: $last_name,
        password: $password
    })
RETURN
    ID(u) AS id,
    u.email AS email,
    u.first_name AS first_name,
    u.last_name AS last_name;
