# URLShortener

URL Sortener implemented in Go using gin-gonic framework for routing and Gorm for Modeling & Database operations.

Through the web application the user can generate a short link for his/her URL. 
The user also has an option to upload a json file with an array of URLs to create short link for (which are uploaded concurrently using Goroutines).


After starting the app through go run main.go & setting up the Env. variables for Database 
the web app will be running on http://localhost:8080

Rest of the web app flow is self-explanatory

