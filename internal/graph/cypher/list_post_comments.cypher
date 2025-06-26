MATCH
    (u:User)-[:WRITE]->(c:Comment)-[:ON]->(p:Post {id: $id})
OPTIONAL MATCH
    (c)-[l:LIKE]-()
WITH
    u,p,c,count(l) AS likes
RETURN
    c.id AS id,
    c.content AS content,
    u.id AS author_id,
    p.id AS post_id,
    likes;
