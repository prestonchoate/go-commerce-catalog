# Go Commerce Catalog Microservice

This project is a simple CRUD microservice that stores a catalog of products. I built this project as a tool to help me learn and understand GO. The project utilizes the [go-chi](https://github.com/go-chi/chi) router for its extension of the native net/http implementation and its complex route paramaterization, grouping, and middleware handling. Currently I am utilizing [planetscale](https://planetscale.com/) for hosting MySql, but any MySql database could be used for this project. In the future I will provide a docker-compose orchestration to handle all of the software dependencies to run this project, but for now nothing is containerized.

## Running Locally

- Make sure Go is installed locally and available when you run `go version` from your terminal
- Create a `.env` file with correct data based on the `.env.sample` file
- `go run .` will run the project by executing the `main.go` file
- Browse to the site based on the `PORT` specified in the .env and verify the server is serving content
