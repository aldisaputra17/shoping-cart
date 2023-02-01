package helpers

import (
	"strconv"
	"strings"

	"github.com/aldisaputra17/shoping-cart/entities"
	"github.com/gin-gonic/gin"
)

type PaginationResult struct {
	Result interface{}
	Error  error
}

func GeneratePaginationRequest(context *gin.Context) *entities.Pagination {
	// default limit, page & sort parameter
	limit := 10
	page := 0
	sortName := "name asc"
	sortQuantity := "quantity asc"

	var searchs []entities.Search

	query := context.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "sort_name":
			sortName = queryValue
			break
		case "sort_quantity":
			sortQuantity = queryValue
		}

		// check if query parameter key contains dot
		if strings.Contains(key, ".") {
			// split query parameter key by dot
			searchKeys := strings.Split(key, ".")

			// create search object
			search := entities.Search{Column: searchKeys[0], Action: searchKeys[1], Query: queryValue}

			// add search object to searchs array
			searchs = append(searchs, search)
		}
	}

	return &entities.Pagination{Limit: limit, Page: page, SortName: sortName, SortQuatity: sortQuantity, Searchs: searchs}
}
