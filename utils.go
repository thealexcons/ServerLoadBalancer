package loadbalancer

import (
	"net/http"
)

const (
	attempts = iota
	retry
)

// getAttemptsFromContext returns the attempts for request
func getAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(attempts).(int); ok {
		return attempts
	}
	return 1
}

// getRetryFromContext returns the retry for request
func getRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(retry).(int); ok {
		return retry
	}
	return 0
}

// Some helper functions for the WRR algorithm

func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

func calculateGCD(weights []int) int {
	result := weights[0]
	for i := 1; i < len(weights); i++ {
		result = gcd(result, weights[i])
	}
	return result
}

func maxValue(weights []int) int {
	m := 0
	for i, e := range weights {
		if i == 0 || e > m {
			m = e
		}
	}
	return m
}

func allSame(weights []int) bool {
	f := weights[0]
	for i := 1; i < len(weights); i++ {
		if weights[i] != f {
			return false
		}
	}
	return true
}
