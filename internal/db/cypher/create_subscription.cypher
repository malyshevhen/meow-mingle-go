MATCH
    (u:User {id: $user_id}),
    (s:User {id: $subscription_id})
CREATE
    (u)-[:SUBSCRIBE]->(s);
