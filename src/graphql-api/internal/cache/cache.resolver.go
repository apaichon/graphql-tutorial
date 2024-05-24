package cache

import (
	"encoding/json"
	"fmt"

	// "graphql-api/internal/auth"
	"log"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/spf13/viper"
)

var cacheDB *Cache

// Middleware to enforce authorization based on permission
func GetCacheResolver(next func(p graphql.ResolveParams) (interface{}, error)) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {

		args := concatMapToString(p.Args)
		hashKey := p.Info.ParentType.Name() + "." + p.Info.FieldName + "_" + args // auth.HashString(args)
		log.Printf("\nRead cached Key:[%s] to get data\n", hashKey)
		cacheDB = NewCache(IntToCacheBackend(viper.GetInt("CACHE_PROVIDER")))
		// if exist cached
		jsonData, err := cacheDB.Get(hashKey)
		// fmt.Printf("Cache Data:%s | error:%v \n", jsonData, err)

		if err == nil {
			// Convert the JSON strings to the appropriate data structures
			arrayData, err := convertToSliceOfMaps(jsonData)
			if err == nil {
				log.Println("[Success] Get Data from cached")
				return arrayData, nil
			}
			objectData, err := convertToMap(jsonData)
			if err == nil {
				log.Println("[Success] Get Data from cached")
				return objectData, nil
			}
		}

		// Execute the resolver if permission is granted
		return next(p)
	}
}

// Middleware to enforce authorization based on permission
func SetCacheResolver(p graphql.ResolveParams, data interface{}) {
	args := concatMapToString(p.Args)
	hashKey := p.Info.ParentType.Name() + "." + p.Info.FieldName + "_" + args // auth.HashString(args)
	log.Printf("\nRead cached Key:[%s] to set data", hashKey)
	jsonData, err := json.Marshal(data)
	cacheDB = NewCache(IntToCacheBackend(viper.GetInt("CACHE_PROVIDER")))
	if err == nil {
		cacheDB.Set(hashKey, string(jsonData))
		log.Printf("\nSet Data to cached [%s]", hashKey)
	}
}

func RemoveGetCacheResolver(key string) {
	cacheDB = NewCache(IntToCacheBackend(viper.GetInt("CACHE_PROVIDER")))
	cacheDB.Removes(key)
}

func concatMapToString(m map[string]interface{}) string {
	var result strings.Builder

	// Iterate over the map and concatenate values
	for _, value := range m {
		strValue := fmt.Sprintf("%v", value)
		result.WriteString(strValue)
	}

	return result.String()
}

// Function to convert a JSON string to []map[string]interface{} if it represents an array
func convertToSliceOfMaps(jsonStr string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	return data, err
}

// Function to convert a JSON string to map[string]interface{} if it represents an object
func convertToMap(jsonStr string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	return data, err
}
