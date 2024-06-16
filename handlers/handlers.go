package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Course struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Heading     string `json:"heading"`
	IsPublished bool   `json:"is_published"`
	ImageURL    string `json:"image_url"`
}

type CoursesResponse struct {
	Courses []Course `json:"courses"`
	Meta    Meta     `json:"meta"`
}

type Meta struct {
	Total         int `json:"total"`
	Page          int `json:"page"`
	From          int `json:"from"`
	To            int `json:"to"`
	PerPage       int `json:"per_page"`
	NumberOfPages int `json:"number_of_pages"`
}

type Enrollment struct {
	UserID          int    `json:"user_id"`
	EnrolledAt      string `json:"enrolled_at"`
	CompletedAt     string `json:"completed_at"`
	PercentComplete int    `json:"percent_complete"`
	ExpiresAt       string `json:"expires_at"`
}

type EnrollmentsResponse struct {
	Enrollments []Enrollment `json:"enrollments"`
	Meta        Meta         `json:"meta"`
}

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	LastSignInIP string `json:"last_sign_in_ip"`
	Role         string `json:"role"`
	Courses      []struct {
		CourseID           int    `json:"course_id"`
		CourseName         string `json:"course_name"`
		EnrolledAt         string `json:"enrolled_at"`
		IsActiveEnrollment bool   `json:"is_active_enrollment"`
		CompletedAt        string `json:"completed_at"`
		PercentComplete    int    `json:"percent_complete"`
	} `json:"courses"`
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func getCourses() []Course {
	apiKey := os.Getenv("API_KEY")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://developers.teachable.com/v1/courses", nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("apiKey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching courses:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil
	}

	var coursesResponse CoursesResponse
	err = json.Unmarshal(body, &coursesResponse)
	if err != nil {
		log.Println("Error unmarshaling courses:", err)
		return nil
	}

	return coursesResponse.Courses
}

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	courseName := r.URL.Query().Get("course_name")

	var courses []Course
	if courseName == "" {
		courses = getCourses()
	} else {
		courses = []Course{}
		foundCourse := findCourseByName(courseName)
		if foundCourse != nil {
			courses = append(courses, *foundCourse)
		}
	}

	enrollmentDetails := make([]map[string]interface{}, 0)
	for _, course := range courses {
		enrollments := getEnrollments(course.ID)
		courseEnrollments := make([]map[string]interface{}, 0)
		for _, enrollment := range enrollments {
			user := getUserDetails(enrollment.UserID)
			courseEnrollments = append(courseEnrollments, map[string]interface{}{
				"user_id":          enrollment.UserID,
				"enrolled_at":      enrollment.EnrolledAt,
				"completed_at":     enrollment.CompletedAt,
				"percent_complete": enrollment.PercentComplete,
				"expires_at":       enrollment.ExpiresAt,
				"user_name":        user.Name,
				"user_email":       user.Email,
			})
		}
		enrollmentDetails = append(enrollmentDetails, map[string]interface{}{
			"course":      course,
			"enrollments": courseEnrollments,
		})
	}

	json.NewEncoder(w).Encode(enrollmentDetails)
}

func findCourseByName(courseName string) *Course {
	courses := getCourses()
	for _, course := range courses {
		if course.Name == courseName {
			return &course
		}
	}
	return nil
}

func getEnrollments(courseID int) []Enrollment {
	apiKey := os.Getenv("API_KEY")
	client := &http.Client{}
	url := fmt.Sprintf("https://developers.teachable.com/v1/courses/%d/enrollments", courseID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("apiKey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching enrollments:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil
	}

	var enrollmentsResponse EnrollmentsResponse
	err = json.Unmarshal(body, &enrollmentsResponse)
	if err != nil {
		log.Println("Error unmarshaling enrollments:", err)
		return nil
	}

	return enrollmentsResponse.Enrollments
}

func getUserDetails(userID int) User {
	apiKey := os.Getenv("API_KEY")
	client := &http.Client{}
	url := fmt.Sprintf("https://developers.teachable.com/v1/users/%d", userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return User{}
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("apiKey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching user details:", err)
		return User{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return User{}
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("Error unmarshaling user details:", err)
		return User{}
	}

	return user
}
