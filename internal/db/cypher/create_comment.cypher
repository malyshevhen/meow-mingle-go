MATCH
    (u:User {id: $author_id}),
    (p:Post {id: $post_id})
CREATE
    (u)-[:WRITE {role: 'Author'}]->(c:Comment {id: $id, content: $content})-[:ON]->(p)
RETURN
    c.id AS id,
    c.content AS content,
    p.id AS post_id,
    u.id AS author_id;
