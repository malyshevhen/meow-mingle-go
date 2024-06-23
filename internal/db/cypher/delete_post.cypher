MATCH (u:User {id: $author_id})-[:WRITE]->(p:Post {id: $id})<-[:ON]-(c:Comment)
DETACH DELETE p,c;
