MATCH (f:User)-[s:SUBSCRIBE]->(u:User)
WHERE ID(f)=$user_id AND ID(u)=$subscription_id
DELETE s;
