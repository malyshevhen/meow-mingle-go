MATCH (u:User)-[l:LIKE]->(p:Post)
WHERE ID(u)=$user_id AND ID(p)=$post_id
DELETE l;
