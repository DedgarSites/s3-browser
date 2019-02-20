package controllers

import (
	"strings"

	"github.com/dedgarsites/s3-browser/bucket"
	"github.com/dedgarsites/s3-browser/datastores"
	"github.com/dedgarsites/s3-browser/tree"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"

	"github.com/gorilla/sessions"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	asession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"fmt"
	"net/http"
)

type filePath struct {
	Path string
}

// GET tree, /all/
func GetTree(c echo.Context) error {
	return c.Render(http.StatusOK, "folder.html", tree.RootFolder)
}

// GET /all/*
func GetTreeAll(c echo.Context) error {
	// TODO check PostMap before continuing with recursive logic
	tempFolder := tree.RootFolder

	if len(c.ParamValues()) == 0 {
		return c.Render(http.StatusOK, "folder.html", tempFolder)
	}

	fullPath := c.ParamValues()[0] // errors if base path
	fmt.Println(tempFolder)

	path := strings.Split(fullPath, "/")

	for _, dir := range path {
		foundFolder := tree.FindNode(tempFolder, dir)
		if foundFolder.Name != "" {
			tempFolder = foundFolder
		}
	}
	if tempFolder.Name == "" || strings.HasSuffix(fullPath, "/") {
		fmt.Println(tempFolder)
		return c.Render(http.StatusOK, "folder.html", tempFolder)
	}

	extension := strings.Split(path[len(path)-1], ".")[1]

	gotItem := bucket.PreSigned(fullPath)

	fileMap := map[string]interface{}{
		"gotItem": gotItem,
	}

	if _, ok := datastores.ImageTypes[extension]; ok {
		return c.Render(http.StatusOK, "image.html", fileMap)
	}

	if _, ok := datastores.VideoTypes[extension]; ok {
		return c.Render(http.StatusOK, "video.html", fileMap)
	}

	return c.Render(http.StatusNotFound, "e04.html", "404 not found")
}

// GET /
func GetMain(c echo.Context) error {
	return c.Render(http.StatusOK, "main.html", datastores.PostMap)
}

// GET /login
func GetLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

// GET /about
func GetAbout(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", nil)
}

// GET /contact
func GetContact(c echo.Context) error {
	return c.Render(http.StatusOK, "contact.html", nil)
}

// GET /login
func GetRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

// GET /privacy
func GetPrivacy(c echo.Context) error {
	return c.Render(http.StatusOK, "privacy.html", nil)
}

// GET /graph
func GetGraph(c echo.Context) error {
	sess, _ := session.Get("session", c)
	if _, ok := sess.Values["current_user"].(string); ok {
		graphGet := map[string]int{"January": 100, "February": 200, "March": 300, "April": 400, "May": 500, "June": 600, "July": 700, "August": 800, "September": 900, "October": 1000, "November": 1100, "December": 1200}

		graphMap := map[string]interface{}{"graphMap": graphGet}
		return c.Render(http.StatusOK, "graph_a.html", graphMap)
	}
	return c.Redirect(http.StatusPermanentRedirect, "/login")
}

// GET /api/graph
func GetApiGraph(c echo.Context) error {
	callback := c.QueryParam("callback")
	month := []string{"January", "February", "March", "April", "May"} //, "June", "July", "August", "September", "October", "November", "December"}
	content := make(map[string]int)
	for i, item := range month {
		content[item] = (i + 1) * 300
	}
	return c.JSONP(http.StatusOK, callback, &content)
}

// GET /api/file
func GetApiFile(c echo.Context) error {
	callback := c.QueryParam("callback")

	content := "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg=="
	return c.JSONP(http.StatusOK, callback, &content)
}

// POST /post-contact
func PostContact(c echo.Context) error {

	if strings.Contains(c.FormValue("message"), "http") && strings.Contains(c.FormValue("message"), "dedgar.com/") == false {
		return c.String(http.StatusOK, "Form submitted")
	}

	TextBody := c.FormValue("name") + "\n" + c.FormValue("email") + "\n" + c.FormValue("message")

	sess, err := asession.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(datastores.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(datastores.CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(datastores.CharSet),
				Data:    aws.String(datastores.Subject),
			},
		},
		Source: aws.String(datastores.Sender),
	}

	result, err := svc.SendEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}

	}
	fmt.Println(c.FormValue("name"))
	fmt.Println(c.FormValue("email"))
	fmt.Println(c.FormValue("message"))
	fmt.Println("Email Sent to address: " + datastores.Recipient)
	fmt.Println(result)
	return c.String(http.StatusOK, "Form submitted")
}

// GET /post/:postname
func GetPost(c echo.Context) error {
	post := c.Param("postname")
	if _, ok := datastores.PostMap[post]; ok {
		return c.Render(http.StatusOK, post+".html", post)
	}
	return c.Render(http.StatusNotFound, "e04.html", "404 Post not found")
}

// GET /post
func GetPostView(c echo.Context) error {
	return c.Render(http.StatusOK, "post_view.html", datastores.PostMap)
}

// GET /trial
func GetTrial(c echo.Context) error {
	sess, _ := session.Get("session", c)
	logged_in_user := sess.Values["current_user"].(string)
	return c.String(http.StatusOK, logged_in_user)
}

// handle any error by attempting to render a custom page for it
func Custom404Handler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.Render(code, errorPage, code); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
			MainSession(c)
			sess, _ := session.Get("session", c)

			if sess.Values["authenticated"] == "true" {
				fmt.Println("in if block")
				fmt.Println(sess.Values)
				return next(c)
			}
			//return next(c)
			return c.Redirect(http.StatusTemporaryRedirect, "/login/google")
		}
	}
}

func MainSession(c echo.Context) { //error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 86400, // * 7,
		//HttpOnly: true,
	}
	//sess.Values["authenticated"] = "true"
	sess.Save(c.Request(), c.Response())
}
