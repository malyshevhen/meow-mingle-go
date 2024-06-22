MATCH
    (u:User),
    (p:Post)
WHERE
    ID(u)=$author_id AND ID(p)=$post_id
CREATE
    (c:Comment {content: $content})-[:ON]->(p),
    (u)-[:WRITE {role: 'Author'}]->(c)
RETURN
    ID(c) AS id,
    c.content AS content,
    ID(p) AS post_id,
    ID(u) AS author_id;
