MATCH (u:User)-[l:LIKE]->(c:Comment)
WHERE ID(u)=$user_id AND ID(c)=$comment_id
DELETE l;
