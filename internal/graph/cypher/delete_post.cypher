MATCH (p:Post {id: $id})<-[:ON]-(c:Comment)
DETACH DELETE p,c;
