# golang-api-test
REST API test by Golang

### Including
 * auto build and restart server [pilu/fresh](https://github.com/pilu/fresh)
 * Web Framework gin [gin-gonic/gin](https://github.com/gin-gonic/gin)
 * MySQL driver [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
 * ORM [jinzhu/gorm](https://github.com/jinzhu/gorm)
 * .env load [joho/godotenv](https://github.com/joho/godotenv)

### Install & Run

```bash
go get github.com/pilu/fresh
go get github.com/gin-gonic/gin
go get -u github.com/go-sql-driver/mysql
go get -u github.com/jinzhu/gorm
go get github.com/joho/godotenv
fresh
```
POST to `http://localhost:4000/api/v1/accounts`
