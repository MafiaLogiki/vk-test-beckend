package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type contextKey string

type ValidatorFunc func(value string) (any, error)

const (
	PageKey contextKey = "page"
	LimitKey contextKey = "limit"
	SortKey contextKey = "sort_by"
	MinPrice contextKey = "min_price"
	MaxPrice contextKey = "max_price"
)

type ValidationRule struct {
	ParamName    string
	DefaultValue string
	Validator    ValidatorFunc
	ContextKey   contextKey
}

func ValidatePositiveInt(value string) (any, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("%s must be an integer", value)
	}
	if i <= 0 {
		return nil, fmt.Errorf("%s must be positive integer", value)
	}

	return i, nil
}

func ValidateByMap(m map[string]string) ValidatorFunc {
	return func (value string) (any, error ) {
		if _, ok := m[value]; !ok {
			return nil, fmt.Errorf("invalid value")
		}
		return m[value], nil
	}
}

func ValidateQueryParams(rules ...ValidationRule) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			query := r.URL.Query()

			for _, rule := range rules {
				valueStr := query.Get(rule.ParamName)
				if valueStr == "" {
					valueStr = rule.DefaultValue
				}

				validatedValue, err := rule.Validator(valueStr)
				if err != nil {
					msg := fmt.Sprintf("Invalid query parameter %s: %v", rule.ParamName, err)
					http.Error(w, msg, http.StatusBadRequest) 
					return
				}
				ctx = context.WithValue(ctx, rule.ContextKey, validatedValue)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidateRequestBody[T any](validatorFunc func(T) error) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc( func (w http.ResponseWriter, r *http.Request){
			var data T
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}

			if err := validatorFunc(data); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), "validatedData", data)
			next.ServeHTTP(w, r.WithContext(ctx))			
		})
	}
}

func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
