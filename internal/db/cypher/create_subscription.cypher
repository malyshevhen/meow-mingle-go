MATCH
    (u:User),
    (s:User)
WHERE
    ID(u)=$user_id AND ID(s)=$subscription_id
CREATE
    (u)-[:SUBSCRIBE]->(s);
