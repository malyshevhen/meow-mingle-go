MATCH (u:User)-[:WRITE]->(p:Post)<-[:ON]-(c:Comment)
WHERE ID(p)=$id AND ID(u)=$author_id
DETACH DELETE p,c;
