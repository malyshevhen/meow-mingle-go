CREATE (u:User {email: $email, first_name: $first_name, last_name: $last_name, password: $password}) RETURN u;
