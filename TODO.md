# TODO

## This sprint

Goal:

- write view handler and page

Up next:

- ~~fix redirect in view~~
- rewrite view handler to show read from db
  - ~~process URI correctly~~
  - ~~query database~~
  - ~~set of Message struct~~
  - check for sql injection for view handler
- rewrite view template to display paste
- *refactor*
- add time left for paste on view
- add view count to table
- implement burn after reading
- implement 401 for deleted views
- implement periodic msg deletion
- *refactor*

## Other
Broad points:

- write server.go
- write html/\*
    - use boostrap

For later:

- use godoc
- write about page copy 
- write install, authors, and readme
- figure out setting global debug and in\_production

v1 complete when:

- implements all features on 0bin.net (except file upload)
  - upload text with timeout
  - upon upload redirects to viewpage
  - viewpage works with # in uri

v2 complete when:

- user can set encryption type and key
  - encryption types should be AES and RSA


v3 complete when:

- user can upload different filetypes
