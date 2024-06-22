// MATCH (p:Post) WHERE ID(p)=$id RETURN p;
MATCH
    (p:Post),
    (u:User)
WHERE
    ID(p)=11 AND (u)-[:WRITE]->(p)
RETURN
    ID(p) AS id,
    p.content AS content,
    ID(u) AS author_id;
