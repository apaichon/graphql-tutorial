# Module 8: Performance and Monitoring
## Lab8.1 - Cache
**Objective:** เข้าใจจุดประสงค์การใช้ Cache และประยุกต์​ใช้
**ไฟล์ทีี่เกี่ยวข้องใน Lab นี้**
```plantuml
@startmindmap
* data
** event.db
* src
** graphql-api
*** cmd
**** server
***** main.go
*** config
**** config.go
*** internal
**** cache
***** cache-resolver.go
***** cache.go
***** redis.go
***** sqlite.go
*** pkg
**** graphql
***** resolvers
****** contact.resolver.go
***** types.go
@endmindmap

```

**Prerequisit**
ติดตั้ง Redis ใน Docker
```sh
docker run -p 6379:6379 --name redis -d redis

```

1. สร้าง Redis object ที่โฟลเดอร์ internale/cache/cache.go ป้อนโค้ด ดังนี้

```go
package cache

import (
	"context"
	"fmt"
	"time"
	"sync"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	redisOnce     sync.Once
	redisInstance *RedisClient
)
// RedisClient represents a simple Redis client.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client.
func NewRedisClient() (*RedisClient, error) {

	addr:= viper.GetString("CACHE_CON_STR")
	password:= viper.GetString("CACHE_PASSWORD")
	db:= viper.GetInt("CACHE_INDEX")
	
	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Ping the Redis server to check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisClient{client: client}, nil
}

// GetInstance returns the singleton instance of the Redis client.
func  GetRedisInstance() (*RedisClient, error) {
	redisOnce.Do(func() {
		var err error
		redisInstance, err = NewRedisClient()
		if err != nil {
			log.Fatalf("Error creating Redis client: %v", err)
		}
	})
	return redisInstance, nil
}

// Close closes the Redis connection.
func (rc *RedisClient) Close() error {
	return rc.client.Close()
}

// Get retrieves the value associated with the given key from Redis.
func (rc *RedisClient) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key '%s' not found", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get value for key '%s': %v", key, err)
	}
	return val, nil
}

// Set sets the value associated with the given key in Redis.
func (rc *RedisClient) Set(key, value string) error {
	ctx := context.Background()
	cacheAge:= viper.GetInt("CACHE_AGE")
	err := rc.client.Set(ctx, key, value, time.Duration(cacheAge) *time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to set value for key '%s': %v", key, err)
	}
	return nil
}

// Remove removes the specified key from Redis.
func (rc *RedisClient) Remove(key string) error {
	ctx := context.Background()
	deleted, err := rc.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to remove key '%s' from Redis: %v", key, err)
	}
	if deleted == 0 {
		return fmt.Errorf("key '%s' does not exist in Redis", key)
	}
	return nil
}

// Remove removes the specified key from Redis.
func (rc *RedisClient) Removes(key string) {
	ctx := context.Background()
	 // Use Lua script to delete keys by pattern
	 script := `
	 local keys = redis.call('KEYS', ARGV[1])
	 for i=1,#keys do
		 redis.call('DEL', keys[i])
	 end
	 return keys
 `
	// Execute Lua script
	result, err := rc.client.Eval(ctx, script, []string{}, key +"*").Result()
	if err != nil {
		panic(err)
	}

	// Print deleted keys
	deletedKeys, _ := result.([]interface{})
	log.Printf("deletedKeys%v", deletedKeys...)
	for _, key := range deletedKeys {
		// deleted, err := 
		rc.client.Del(ctx, key.(string)) //.Result() 
	}
}
```

2. สร้าง Sqlite Object สำหรับทำ Cache ด้วย SQLite ที่ internal/cached/sqlite.go ป้อนโค้ด ดังนี้

```go
package cache

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	sqliteOnce     sync.Once
	sqliteInstance *SQLiteInMemClient
)

// SQLiteClient represents a simple SQLite client.
type SQLiteInMemClient struct {
	db *sql.DB
}

// NewSQLiteClient creates a new SQLite client with an in-memory database.
func NewSQLiteInMemClient() (*SQLiteInMemClient, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite database: %v", err)
	}
	return &SQLiteInMemClient{db: db}, nil
}

// GetInstance returns the singleton instance of the SQLite client.
func GetSqliteInMemInstance() (*SQLiteInMemClient, error) {
	sqliteOnce.Do(func() {
		var err error
		sqliteInstance, err = NewSQLiteInMemClient()
		if err != nil {
			log.Fatalf("Error creating SQLite client: %v", err)
		}
	})
	return sqliteInstance, nil
}

// Close closes the SQLite database connection.
func (sc *SQLiteInMemClient) Close() error {
	return sc.db.Close()
}

// CreateTable creates a table in the SQLite database.
func (sc *SQLiteInMemClient) CreateTable() error {
	_, err := sc.db.Exec(`CREATE TABLE IF NOT EXISTS cache (
		key TEXT PRIMARY KEY,
		value TEXT
	)`)
	if err != nil {
		return fmt.Errorf("failed to create table in SQLite database: %v", err)
	}
	return nil
}

// Get retrieves the value associated with the given key from the SQLite database.
func (sc *SQLiteInMemClient) Get(key string) (string, error) {
	var value string
	err := sc.db.QueryRow("SELECT value FROM cache WHERE key = ?", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("key '%s' not found", key)
		}
		return "", fmt.Errorf("failed to get value for key '%s': %v", key, err)
	}
	return value, nil
}

// Set sets the value associated with the given key in the SQLite database.
func (sc *SQLiteInMemClient) Set(key, value string) error {
	_, err := sc.db.Exec("INSERT INTO cache(key, value) VALUES(?, ?)", key, value)
	if err != nil {
		return fmt.Errorf("failed to set value for key '%s': %v", key, err)
	}
	return nil
}

// Remove removes the specified key from the SQLite database.
func (sc *SQLiteInMemClient) Remove(key string) error {
	res, err := sc.db.Exec("DELETE FROM cache WHERE key = ?", key)
	if err != nil {
		return fmt.Errorf("failed to remove key '%s' from SQLite database: %v", key, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows after removing key '%s': %v", key, err)
	}
	if count == 0 {
		return fmt.Errorf("key '%s' does not exist in SQLite database", key)
	}
	return nil
}

// Remove removes the specified key from the SQLite database.
func (sc *SQLiteInMemClient) Removes(key string)  {
	sc.db.Exec("DELETE FROM cache WHERE key Like ?", "%"+key+"*%")
}

```

3. สร้าง Cache Interface สำหรับรองรับการเลือกใช้งาน Cache ระหว่าง Redis หรือ Sqlite ที่ไฟล์ *internal/cache/cache.go*  ป้อนโค้ด ดังนี้ 

```go
package cache

import (
	"log"
)

// CacheBackend represents the backend for the cache.
type CacheBackend int

const (
	RedisBackend CacheBackend = iota
	SQLiteBackend
)

func IntToCacheBackend(i int) CacheBackend {
	switch i {
	case 0:
		return RedisBackend
	case 1:
		return SQLiteBackend
	default:
		return RedisBackend // or return an error, depending on your use case
	}
}

// Cache represents a cache with support for different backends.
type Cache struct {
	backend CacheBackend
	db      CacheDB
}

// CacheDB represents the interface for interacting with the cache database.
type CacheDB interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Remove(key string) error
	Removes(key string)
	Close() error
}

// NewCache creates a new cache with the specified backend.
func NewCache(backend CacheBackend) *Cache {
	var db CacheDB
	switch backend {
	case RedisBackend:
		db,err := GetRedisInstance()
		if (err !=nil) {
			panic(err)
		}
		return &Cache{backend: backend, db: db}
	case SQLiteBackend:
		// db = &SQLiteInMemClient{}
		db,err := GetSqliteInMemInstance()
		if (err !=nil) {
			panic(err)
		}
		return &Cache{backend: backend, db: db}
	default:
		log.Fatalf("Unsupported cache backend: %v", backend)
	}
	return &Cache{backend: backend, db: db}
}

// Get retrieves the value associated with the given key from the cache.
func (c *Cache) Get(key string) (string, error) {
	return c.db.Get(key)
}

// Set sets the value associated with the given key in the cache.
func (c *Cache) Set(key, value string) error {
	return c.db.Set(key, value)
}

// Remove removes the specified key from the cache.
func (c *Cache) Remove(key string) error {
	return c.db.Remove(key)
}

// Remove removes the specified key from the cache.
func (c *Cache) Removes(key string) {
	c.db.Removes(key)
}
```

4. เพิ่มค่า Configuration เกี่ยวกับ Cache ที่ไฟล์ .env
```bash
CACHE_PROVIDER=0 #Redis
CACHE_CON_STR=127.0.0.1:6379
CACHE_PASSWORD=
CACHE_INDEX=0
CACHE_AGE=300
```
5. สร้าง Cache Resovlver เพื่อนำไปใส่ใน Resolver ที่ไฟล์ *internale/cache/cache-resolver.go* ป้อนโค้ด ดังนี้

```go
package cache

import (
	"encoding/json"
	"fmt"
	// "graphql-api/internal/auth"
	"strings"
	"log"
	"github.com/graphql-go/graphql"
	"github.com/spf13/viper"
)

var cacheDB *Cache

func init() {
	cacheDB = NewCache(IntToCacheBackend(viper.GetInt("CACHE_PROVIDER")))
}

// Middleware to enforce authorization based on permission
func GetCacheResolver(next func(p graphql.ResolveParams) (interface{}, error)) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {

		args := concatMapToString(p.Args)
		hashKey := p.Info.ParentType.Name() + "." + p.Info.FieldName + "_" + args // auth.HashString(args)
		log.Printf("\nRead cached Key:[%s] to get data\n", hashKey)
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
	if err == nil {
		cacheDB.Set(hashKey, string(jsonData))
		log.Printf("\nSet Data to cached [%s]", hashKey)
	}
}

func RemoveGetCacheResolver(key string) {
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

```

6. นำฟังก์ชั่น GetCacheResolver ไปใส่ในฟังก์ชั่น ContactQueriesType ในส่วนของ Gets ต่างๆ

```go
// pkg/graphql/types.go
import (
    ...
	"graphql-api/internal/cache"
)
var ContactQueriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactQueries",
	Fields: graphql.Fields{
		"gets": &graphql.Field{
			Type:    graphql.NewList(ContactGraphQLType),
			Args:    SearhTextQueryArgument,
            // Add cache.GetCacheResolver
			Resolve: auth.AuthorizeResolverClean("contacts.gets", cache.GetCacheResolver(resolvers.GetContactResolve)),
		},
		"getPagination": &graphql.Field{
			Type:    ContactPaginationGraphQLType,
			Args:    SearhTextPaginationQueryArgument,
			 // Add cache.GetCacheResolver
            Resolve: auth.AuthorizeResolverClean("contacts.getPagination", cache.GetCacheResolver(resolvers.GetContactsPaginationResolve)),
		},
		"getById": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    IdArgument,
			 // Add cache.GetCacheResolver
            Resolve: auth.AuthorizeResolverClean("contacts.getById", cache.GetCacheResolver(resolvers.GetContactByIdResolve)),
		},
	},
})
```

7. นำฟังก์ชั่น SetCacheResolver ไปใส่ใน pkg/graphql/resolvers/contact.resolver.go หรือ resolver get ที่ต้องการ ในส่วนของ Gets ต่างๆ หลังจากประมวลผลสำเร็จแล้ว

```go
func GetContactResolve(params graphql.ResolveParams) (interface{}, error) {
	// Update limit and offset if provided
	limit, ok := params.Args["limit"].(int)
	if !ok {
		limit = 10
	}

	offset, ok := params.Args["offset"].(int)
	if !ok {
		offset = 0
	}

	searchText, ok := params.Args["searchText"].(string)
	if !ok {
		searchText = ""
	}
	contactRepo := contact.NewContactRepo()

	// Fetch contacts from the database
	contacts, err := contactRepo.GetContactsBySearchText(searchText, limit, offset)
	if err != nil {
		return nil, err
	}
	// time.Sleep(1 * time.Minute)
    // Set Cached
	go cache.SetCacheResolver(params, contacts)
	
	return contacts, nil
}

func GetContactsPaginationResolve(params graphql.ResolveParams) (interface{}, error) {
	// Update limit and offset if provided
	page, ok := params.Args["page"].(int)
	if !ok {
		page = 1
	}

	pageSize, ok := params.Args["pageSize"].(int)
	if !ok {
		pageSize = 10
	}

	searchText, ok := params.Args["searchText"].(string)
	if !ok {
		searchText = ""
	}
	contactRepo := contact.NewContactRepo()

	// Fetch contacts from the database
	contacts, pager, err := contactRepo.GetContactsBySearchTextPagination(searchText, page, pageSize)
	var contactPagination = models.ContactPaginationModel{
		Contacts:   contacts,
		Pagination: pager,
	}

	if err != nil {
		return nil, err
	}

    // Set Cache
	go cache.SetCacheResolver(params, contacts)

	return contactPagination, nil
}

func GetContactByIdResolve(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(int)
	contactRepo := contact.NewContactRepo()

	// Fetch contacts from the database
	contact, err := contactRepo.GetContactByID(id)
	if err != nil {
		return nil, err
	}
    // Set Cache
	go cache.SetCacheResolver(params, contact)
	return contact, nil
}
```

8. นำฟังก์ชั่น RemoveGetCacheResolver ไปใส่ใน CretateContactResolve เมื่อมีการสร้างข้อมูลใหม่เกิดขึ้น เพื่อให้แน่ใจว่าผู้ใช้งานได้ข้อมูลที่อัพเดตเมื่อมีการดึงข้อมูล

```go
func CretateContactResolve(params graphql.ResolveParams) (interface{}, error) {
	// Map input fields to Contact struct
	input := params.Args["input"].(map[string]interface{})
	invalids := validateContact(input)

	if len(invalids) > 0 {
		return nil, fmt.Errorf("%v", invalids)
	}

	contactInput := models.ContactModel{

		Name:      input["name"].(string),
		FirstName: input["first_name"].(string),
		LastName:  input["last_name"].(string),
		GenderId:  input["gender_id"].(int),
		Dob:       input["dob"].(time.Time),
		Email:     input["email"].(string),
		Phone:     input["phone"].(string),
		Address:   input["address"].(string),
		PhotoPath: input["photo_path"].(string),
		CreatedBy: "test-api",
		CreatedAt: time.Now(),
	}

	contactRepo := contact.NewContactRepo()

	// Insert Contact to the database
	id, err := contactRepo.InsertContact(&contactInput)
	if err != nil {
		return nil, err
	}
	contactInput.ContactId = int64(id)
    // Clear Cache to make sure that user will got updated data.
	go cache.RemoveGetCacheResolver("ContactQueries")
	return contactInput, nil
}
```
ึ9. ทดสอบ Query และ Mutation 

```graphql
{
  contacts {
    gets(searchText: "") {
      name
      first_name
      last_name
      gender_id
      dob
      email
      phone
      address
      photo_path
      created_at
      created_by
    }
    getById(id: 1) {
      name
      first_name
      last_name
      gender_id
      dob
      email
      phone
      address
      photo_path
      created_at
      created_by
    }
  }
}

```
**Mutation**
*Query*
```graphql
mutation CreateContact($input: CreateContactInput, $input2: CreateContactInput) {
  contactMutations {
    contact1: createContact(input: $input) {
      contact_id
    }
    contact2: createContact(input: $input2) {
      contact_id
    }
  }
}

```

*Variables*
```json
{
  "input": {
    "name": "ABC",
    "first_name": "PUP3",
    "last_name": "Apaichon4",
    "gender_id": 1,
    "dob": "1979-11-13T00:00:00Z",
    "email": "john@example.com",
    "phone": "123-456-7890",
    "address": "123 Main St",
    "photo_path": "path/to/photo.jpg"
  },
  "input2": {
    "name": "Dr.PUP3",
    "first_name": "PUP3",
    "last_name": "Apaichon3",
    "gender_id": 1,
    "dob": "1979-11-13T00:00:00Z",
    "email": "john@example.com",
    "phone": "123-456-7890",
    "address": "123 Main St",
    "photo_path": "path/to/photo.jpg"
  }
}
```

## Lab8.2 - Open Telemetry
**Objective:** เข้าใจจุดประสงค์การใช้ Open Telemetry สำหร้บทำ ​Trace Monitoring
**ไฟล์ทีี่เกี่ยวข้องใน Lab นี้**
```plantuml
@startmindmap
* data
** event.db
* src
** graphql-api
*** cmd
**** server
***** main.go
*** config
**** config.go
*** internal
**** monitoring
***** otel.go
***** otel.resolver.go
*** pkg
**** graphql
***** resolvers
****** contact.resolver.go
***** types.go
@endmindmap

```

1. ติดตั้งไลบรารีของ Open Telemetry และ exporters zipkin

```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/zipkin
go get go.opentelemetry.io/otel/sdk/trace
go get go.opentelemetry.io/otel/semconv/v1.24.0
```
2. ติดตั้ง zipkin เป็น UI สำหรับดูผลลัพธ์การ Monitor Trace

```bash
docker run -d -p 9411:9411 openzipkin/zipkin
```
3. เขียนโค้ด go ที่ internal/monitoring/otel.go ป้อนโค้ด ดังนี้

```go
package monitoring

import (
	"log"
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)


var logger = log.New(os.Stderr, "graphql-api", log.Ldate|log.Ltime|log.Llongfile)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func InitTracer(url string) (func(context.Context) error, error) {
	// Create Zipkin Exporter and install it as a global tracer.
	//
	// For demoing purposes, always sample. In a production application, you should
	// configure the sampler to a trace.ParentBased(trace.TraceIDRatioBased) set at the desired
	// ratio.
	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("graphql-api"),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

```

4. สร้าง otel resolver ที่ internal/monitoring/otel.resolver.go ป้อนโค้ด ดังนี้

```go
package monitoring

import (
	"github.com/graphql-go/graphql"
	"go.opentelemetry.io/otel"
)

// Middleware to trace resolver functions
func TraceResolver(resolverFunc func(p graphql.ResolveParams) (interface{}, error)) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		// Create a new span for the resolver function
		spanName := p.Info.ParentType.Name() + "." + p.Info.FieldName
		ctx, span := otel.GetTracerProvider().Tracer("").Start(p.Context, spanName)
		defer span.End()

		// Perform the resolver logic
		return resolverFunc(graphql.ResolveParams{Source: p.Source, Args: p.Args, Info: p.Info, Context: ctx})
	}
}

```
5. เพิ่ม Configuration สำหรับ อ้างถึง Trace Server ใน .env ดังนี้
```bash
TRACE_EXPORTER_URL=http://localhost:9411/api/v2/spans
```
6. นำ otel ไปใช้ในไฟล์ cmd/server/main.go import package monitoring ใน main.go และเรียกใช้ฟังก์ชัน InitTracer ดังนี้

```go
import
    (
        ...
        "graphql-api/internal/monitoring")

    
    func main() {
	
	shutdown, err :=  monitoring.InitTracer(viper.GetString("TRACE_EXPORTER_URL"))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()
    ...
```

5. นำ otel resolver ไปใส่ใน resolver ที่ต้องการ monitor
```go
// pkg/graphql/types.go
// Define the ContactQueries type
var ContactQueriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactQueries",
	Fields: graphql.Fields{
		"gets": &graphql.Field{
			Type:    graphql.NewList(ContactGraphQLType),
			Args:    SearhTextQueryArgument,
			Resolve: auth.AuthorizeResolverClean("contacts.gets",
            // add monitor Trace Resolver
            monitoring.TraceResolver( cache.GetCacheResolver(resolvers.GetContactResolve))),
		},
		"getPagination": &graphql.Field{
			Type:    ContactPaginationGraphQLType,
			Args:    SearhTextPaginationQueryArgument,
			Resolve: auth.AuthorizeResolverClean("contacts.getPagination", cache.GetCacheResolver(resolvers.GetContactsPaginationResolve)),
		},
		"getById": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    IdArgument,
			Resolve: auth.AuthorizeResolverClean("contacts.getById", cache.GetCacheResolver(resolvers.GetContactByIdResolve)),
		},
	},
})
```
6. test graphql query
```graphql
{
  contacts {
    gets(searchText: "") {
      name
      first_name
      last_name
      gender_id
      dob
      email
      phone
      address
      photo_path
      created_at
      created_by
    }
  }
}
```
7. เปิด zipkin เพื่อดู monitor ที่ http://localhost:9411

![SQL Tools](./images/lab8/figure8.1.png)
*Figure 8.1 Zipkin Trace Monitor.*


## Lab8.3 - Metric Monitoring

