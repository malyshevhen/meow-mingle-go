MATCH (u:User) WHERE ID(u)=$author_id
CREATE (p:Post {content: $content})<-[:WRITE {role: 'Author'}]-(u)
RETURN p;
