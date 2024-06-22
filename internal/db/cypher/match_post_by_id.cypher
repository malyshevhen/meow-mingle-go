MATCH
    (p:Post),
    (u:User)
OPTIONAL MATCH
    (p)-[l:LIKE]-()
WITH
    u,p,count(l) AS likes
WHERE
    ID(p)=$id AND (u)-[:WRITE]->(p)
RETURN
    ID(p) AS id,
    p.content AS content,
    ID(u) AS author_id,
    likes;
