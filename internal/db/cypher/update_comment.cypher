MATCH
    (u:User {id: $author_id})-[:WRITE]->(c:Comment {id: $id})-[:ON]->(p:Post)
OPTIONAL MATCH
    (c)-[l:LIKE]-()
WITH
    u,p,c,count(l) AS likes
SET
    c.content=$content
RETURN
    c.id AS id,
    c.content AS content,
    u.id AS author_id,
    p.id AS post_id,
    likes;
