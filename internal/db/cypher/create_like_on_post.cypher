MATCH
    (u:User {id: $user_id}),
    (p:Post {id: $post_id})
CREATE
    (u)-[:LIKE]->(p);
