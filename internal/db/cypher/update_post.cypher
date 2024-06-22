MATCH
    (p:Post),
    (u:User)
OPTIONAL MATCH
    (p)-[l:LIKE]-()
WITH
    u,p,count(l) AS likes
WHERE
    ID(p)=$id
    AND ID(u)=$author_id
    AND (u)-[:WRITE]->(p)
SET
    p.content=$content
RETURN
    ID(p) AS id,
    p.content AS content,
    ID(u) AS author_id,
    likes;
