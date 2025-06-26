MATCH
    (u:User {id: $author_id})
CREATE
    (u)-[:WRITE {role: 'Author'}]->(p:Post {id: $id, content: $content})
RETURN
    p.id AS id,
    p.content AS content,
    u.id AS author_id;
