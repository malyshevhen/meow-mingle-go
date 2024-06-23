MATCH
    (u:User)-[:SUBSCRIBE]->(s:User)-[:WRITE]->(p:Post)
OPTIONAL MATCH
    (p)-[l:LIKE]-()
WITH
    u,p,s,count(l) AS likes
WHERE
    ID(u)=$id
RETURN
    ID(p) AS id,
    p.content AS content,
    ID(s) AS author_id,
    likes;
