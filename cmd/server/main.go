package main

import (
	"log"
	"net/http"
	"time"
)

type PostData struct {
	Vectors [][]float64 `json:"vectors"`
}

type PostQuery struct {
	K     int       `json:"k"`
	Query []float64 `json:"query"`
}

func main() {
	http.HandleFunc("/bulk-add", loggingMiddleware(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// body, err := ioutil.ReadAll(r.Body)
		// if err != nil {
		// 	http.Error(rw, "Error reading request body", http.StatusInternalServerError)
		// 	return
		// }

		// var postData PostData
		// if err := json.Unmarshal(body, &postData); err != nil {
		// 	http.Error(rw, "Error decoding JSON", http.StatusBadRequest)
		// 	return
		// }

		// rw.WriteHeader(http.StatusOK)
	}))

	http.HandleFunc("/query", loggingMiddleware(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// body, err := ioutil.ReadAll(r.Body)
		// if err != nil {
		// 	http.Error(rw, "Error reading request body", http.StatusInternalServerError)
		// 	return
		// }

		// var postData PostQuery
		// if err := json.Unmarshal(body, &postData); err != nil {
		// 	http.Error(rw, "Error decoding JSON", http.StatusBadRequest)
		// 	return
		// }

		// log.Printf("Query: Starting\n")
		// result := idx.Query(postData.K, postData.Query)

		// jsonData, err := json.Marshal(result)
		// if err != nil {
		// 	http.Error(rw, "Error converting to JSON", http.StatusInternalServerError)
		// 	return
		// }

		// rw.Header().Set("Content-Type", "application/json")
		// rw.Write(jsonData)
	}))

	log.Println("Starting server")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next(w, r)

		duration := time.Since(start)
		log.Printf("Request took: %v", duration)
	}
}
