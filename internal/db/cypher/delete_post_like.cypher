MATCH (u:User {id: $user_id})-[l:LIKE]->(p:Post {id: $post_id})
DELETE l;
