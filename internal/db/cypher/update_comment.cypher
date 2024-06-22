MATCH
    (u:User)-[:WRITE]->(c:Comment)-[:ON]->(p:Post)
OPTIONAL MATCH
    (c)-[l:LIKE]-()
WITH
    u,p,c,count(l) AS likes
WHERE
    ID(c)=$id AND ID(u)=$author_id
SET
    c.content=$content
RETURN
    ID(c) AS id,
    c.content AS content,
    ID(u) AS author_id,
    ID(p) AS post_id,
    likes;
