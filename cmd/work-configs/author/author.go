/*
// Using the struct methods
authorInfo := author.NewAuthorInfo()
authorInfo.Print()

// Update status and date
authorInfo.UpdateStatus("Production Ready")
authorInfo.UpdateDate("2025-01-15")
authorInfo.Print()

// Using package-level functions (existing approach)
author.PrintAuthorInfo()
*/

package author

import (
	"fmt"
)

const (
	Author string = "LN Kornilov"
	Email  string = "ln.KornilovStar@gmail.com"
	Github string = "github.com/KornilovLN/go-na-practike"
	Data   string = "2025-20-06"
	Status string = "Testing - in progress"
	Info   string = "Go-na-practike"
)

type AuthorInfo struct {
	Author string
	Email  string
	Github string
	Data   string
	Status string
	Info   string
}

// Constructor method - creates new AuthorInfo instance
func NewAuthorInfo() *AuthorInfo {
	return &AuthorInfo{
		Author: Author,
		Email:  Email,
		Github: Github,
		Data:   Data,
		Status: Status,
		Info:   Info,
	}
}

// Method to get formatted string representation of AuthorInfo
func (a *AuthorInfo) String() string {
	return "=== " +
		"Author: " + Author + " ==================================================\n" +
		"Email : " + Email + "\n" +
		"GitHub: " + Github + "\n" +
		"Date  : " + Data + "\n" +
		"Status: " + Status + "\n" +
		"Info  : " + Info + "\n"
}

// Method to print AuthorInfo
func (a *AuthorInfo) Print() {
	fmt.Println(a.String())
}

// Method to update author information
func (a *AuthorInfo) UpdateStatus(newStatus string) {
	a.Status = newStatus
}

// Method to update date
func (a *AuthorInfo) UpdateDate(newDate string) {
	a.Data = newDate
}
