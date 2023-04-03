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

