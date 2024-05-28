package handlers

import (
	"fmt"
	"log"
	"net/http"
	"service/service/internal/lib/model"
	"service/service/pkg/cache"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// GetOrderHandler - handler для получения order по id
func GetOrderHandler(csh *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "id")

		order, exist := csh.Get(orderID)
		if !exist {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("order [%v] not found in cache", orderID)
			w.Write([]byte(fmt.Sprintf("order [%v] not found", orderID)))
			return
		}

		if str, ok := order.(*model.Order); ok {
			log.Printf("Read order [%v] from cache", orderID)
			render.JSON(w, r, str)
		} else {
			log.Printf("Something was wrong when trying type assertion with order [%v] from cache", str)
			w.Write([]byte(fmt.Sprintf("can`t found order [%v]", orderID)))
			return
		}

	}
}
