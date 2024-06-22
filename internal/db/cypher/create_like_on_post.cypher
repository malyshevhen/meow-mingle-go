MATCH
    (u:User),
    (p:Post)
WHERE
    ID(u)=$user_id AND ID(p)=$post_id
CREATE
    (u)-[:LIKE]->(p);
