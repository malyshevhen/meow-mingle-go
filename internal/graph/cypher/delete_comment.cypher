MATCH (u:User {id: $author_id})-[:WRITE]->(c:Comment {id: $id})
DETACH DELETE c;
