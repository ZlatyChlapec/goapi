## Development
Create `.env` file with TOKENS variable. For example:
```shell
$ cat .env
TOKENS=[{"token": "longRandomHash", "user": "Spiderman"},{"token": "longRandomHash", "user": "Batman"}]
```

When you have that you can just run with `docker-compose up --build`.

## Usage
API endpoints:
* /images/
  * GET
    * Returns `application/json`
      ```json
      {
        "1": { // image_id
          "mime": "image/jpeg", // type of the data in `data`
          "data": "_9j_4AAQSkZJRgABAQAAAQ..." // base64 encoded data
        }
      }
      ```
  * POST
    * Requires body of `Content-Type: application/json`
      ```json
      {"mime": "image/jpeg", "data": "_9j_4A..."}
      ```
    * Returns `application/json`
      ```json
      {"image_id": 1}
      ```
* /images/{image_id}
  * GET
    * Returns `application/json`
      ```json
      {"mime": "image/jpeg", "data": "_9j_4A..."}
      ```
* /images/direct/{image_id}
  * GET
    * Returns `Content-Type` of mime type stored with image. Thanks to that it shows image.

All endpoints are behind token authentication. With each request you need to provide Bearer token. You can provider one
by specifying similar header `Authorization: Bearer longRandomHash`