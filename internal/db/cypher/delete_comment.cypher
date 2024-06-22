MATCH (u:User)-[:WRITE]->(c:Comment)
WHERE ID(c)=$id AND ID(u)=$author_id
DETACH DELETE c;
