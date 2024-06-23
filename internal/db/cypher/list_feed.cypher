MATCH
    (u:User {id: $id})-[:SUBSCRIBE]->(s:User)-[:WRITE]->(p:Post)
OPTIONAL MATCH
    (p)-[l:LIKE]-()
WITH
    u,p,s,count(l) AS likes
RETURN
    p.id AS id,
    p.content AS content,
    s.id AS author_id,
    likes;
