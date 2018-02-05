# TODO

## This sprint

Goal:

1) front end encryption
  - add onclick function to paste page
    - ~~generate key~~
    - ~~get the form parameter~~
    - ~~perform the encryption~~
    - ~~replace the form parameters and send~~
  - add decrypt to view page
    - ~~get the key~~
    - ~~get the encrypted message~~
    - ~~decrypt~~
    - ~~replae the message~~

## Other
Broad points:

- write server.go
- write html/\*
    - use boostrap

For later:

- use godoc
- write about page copy 
- improve INSTALL.txt
- figure out setting global debug and in\_production
- refactor viewHandler
- write a log++ package
- use template caching

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
