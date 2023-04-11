Endpoints:

Method: POST Path: /register

Example body:
{
"publicAddress": "0x3110752149AF23Ee65968C2019b7c86D12B32229"
}

Returns status code: 
   - 201 Created when the user successfully created
   - 409 When the user already exists
   - 400 when the request body is invalid
   - 500 if there is a server error

Returns Body:
    - empty

_____________________________________________
Method: GET Path: /users/:publicAddress/nonce

Returns Example Body:
{
"nonce": "63210018627757926317526024290391413217358890723641523832149966690207267728843150165831744512767436400627528585164026452344678510"
}

Returns StatusCode:
- 200 when a nonce is fetched
- 404 when the publicAddress is not registered
- 500 if there is a server error

_____________________________________________
Method: POST Path:/signin

Example body:
{
"publicAddress": "0x3110752149AF23Ee65968C2019b7c86D12B32229",
"signature": "LHB/Efh/BB4JyCUGDIFYp46nutMLyHvwENwd2sss",
"nonce": "63210018627757926317526024290391413217358890723641523832149966690207267728843150165831744512767436400627528585164026452344678510"
}

Returns StatusCode:
- 200 when ok
- 401 when the user is not authenticate
- 500 if there is a server error Returns Example Body: 
      Define me

_____________________________________________
Method: GET Path: /welcome

Returns Example Body:
```
{
 "publicAddress": "0x3110752149AF23Ee65968C2019b7c86D12B32229"
}
```

Returns StatusCode:
- 200 when ok
- 401 when the user is not authenticate
- 500 if there is a server error
