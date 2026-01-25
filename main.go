package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Description string    `json:"description"`
}

var category = [] Category {
	{ID: 1, Name: "Sabun", Description: "Sabun"},
	{ID: 2, Name: "Pasta Gigi", Description: "Pasta Gigi"},
	{ID: 3, Name: "Shampoo", Description: "Shampoo"},
}


func getCreateCategoryHandler (w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(category)
		
	} else if r.Method == "POST" {
		// baca data dari request
		var newCategory Category
		err := json.NewDecoder(r.Body).Decode(&newCategory)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// masukkin data ke dalam variable category
		newCategory.ID = len(category) + 1
		category = append(category, newCategory)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201
		json.NewEncoder(w).Encode(newCategory)
	}
}


func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	// URL: /api/category/123 -> ID = 123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Cari kategori dengan ID tersebut
	for _, p := range category {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	// Kalau tidak found
	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}


func updateCategory(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop kategori, cari id, ganti sesuai data dari request
	for i := range category {
		if category[i].ID == id {
			updateCategory.ID = id
			category[i] = updateCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}
	
	http.Error(w, "Category not found", http.StatusNotFound)
}


func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	
	// ganti id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}
	
	// loop produk cari ID, dapet index yang mau dihapus
	for i, p := range category {
		if p.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			category = append(category[:i], category[i+1:]...)
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Data deleted successfully",
			})
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}



func main(){

	http.HandleFunc("/api/category", getCreateCategoryHandler)

	http.HandleFunc("/api/category/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}

	})

	fmt.Println("Server running at localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed to run")
	}
}