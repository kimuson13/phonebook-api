# phonebook-api
curl --request POST --header "Content-Type: application/json" --data "{\"name\":\"test\",\"phone\":\"09012121212\"}" "http://localhost:8080/api/phonebooks"

curl --request PUT --header "Content-Type: application/json" --data "{\"name\":\"put\",\"phone\":\"09087876565\"}" "http://localhost:8080/api/phonebooks/2"

curl --request DELETE "http://localhost:8080/api/phonebooks/5"

curl --request POST --data '{"name":"test", "phone":"09087654321"}' "http://localhost:8080/api/phonebooks"

curl -X PUT -d "{\"name\":\"put\",\"phone\":\"09087876565\"}" "http://localhost:8080/api/phonebooks/4"

curl -X POST -d "{\"name\":\"add\",\"phone\":\"09087876565\"}" "http://localhost:8080/api/phonebooks"