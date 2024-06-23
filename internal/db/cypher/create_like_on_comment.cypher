MATCH
    (u:User {id: $user_id}),
    (c:Comment {id: $comment_id})
CREATE
    (u)-[:LIKE]->(c);
