# GoMeisa

## GoMeisa is designed like a Trello-Redmine prototype.
 
**Key features:**
- authentication based on sessions;
- server communicates with database (adding users, projects, tasks etc);
- generating invite link. Implemented with the help of _interim_ table "invitations". If user follows generated link
he is added to table "projects_users" and brand new generated link is deleted from "invitations".
- undesirable actions are prevented. User who is not in a certain project can't call its handlers. Every handler checks whether
user is in project or not.
- role functional. Only technical leader (admin) can change project description, pin tasks and remove employees. 
Regular users cannot do this;
- migration up support;
- views based on templates (layout, content, navbar). Minimalistic front-end.

### This project requires a lot of impovements such as writing middleware instead of using repeating code in controllers, building docker container and adding new functional.
