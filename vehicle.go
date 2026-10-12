package main

type Vehicle struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	RegNo      string `json:"reg_no"`
	Vehicle    string `json:"vehicle"`
}
