package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on http://localhost:8080/ping
}

// 1. API handler
// - 1.1. User
// 	- 1.1.1. User Registration
// 	- 1.1.2. User Login
// 	- 1.1.3. User Logout
// 	- 1.1.4. User Profile Update
// 	- 1.1.5. User Profile Delete
// 	- 1.1.6. User Profile View
// - 1.2. Doctor
// 	- 1.2.1. Doctor Registration
// 	- 1.2.2. Doctor Login
// 	- 1.2.3. Doctor Logout
// 	- 1.2.4. Doctor Profile Update
// 	- 1.2.5. Doctor Profile Delete
// 	- 1.2.6. Doctor Profile View
// - 1.3. Chat
// 	- 1.3.1. Chat Start
// 	- 1.3.2. Chat End
// 	- 1.3.3. Chat History
// - 1.4. Review
// 	- 1.4.1. Review Add
// 	- 1.4.2. Review Update
// 	- 1.4.3. Review Delete
// 	- 1.4.4. Review View
// - 1.5. AI Communication
// 	- 1.5.1. AI Communication Start
// 	- 1.5.2. AI Communication End
// 	- 1.5.3. AI Communication History
// - 1.6. Crowdsense Ranking
// 	- 1.6.1. Crowdsense Ranking Start
// 	- 1.6.2. Crowdsense Ranking End
// 	- 1.6.3. Crowdsense Ranking History

// 2. Database
// 3. Middleware
// 4. Model
// 5. Service
// 6. Repository
// 7. Controller
// 8. Route

// 2. AI Communication
// 3. Crowdsense Ranking
// 4. Chat history Store korte hobe
// 5. Doctore review store hobe
