MATCH (u:User {id: $user_id})-[l:LIKE]->(c:Comment {id: $comment_id})
DELETE l;
