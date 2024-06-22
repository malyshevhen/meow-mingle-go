MATCH
    (u:User),
    (c:Comment)
WHERE
    ID(u)=$user_id AND ID(c)=$comment_id
CREATE
    (u)-[:LIKE]->(c);
