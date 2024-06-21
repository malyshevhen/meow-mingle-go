MATCH (u:User {email: $id}) RETURN u.email AS email, u.first_name AS first_name, u.last_name AS last_name;
