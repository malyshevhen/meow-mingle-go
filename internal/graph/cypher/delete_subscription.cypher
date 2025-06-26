MATCH (f:User {id: $user_id})-[s:SUBSCRIBE]->(u:User {id: $subscription_id})
DELETE s;
