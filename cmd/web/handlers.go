package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rohitxdd/snippetbox/internal/models"
	"github.com/rohitxdd/snippetbox/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()

	if err != nil {
		app.ServerError(w, err)
		return
	}
	data := app.newTemplateData(r)

	data.Snippets = snippets

	app.Render(w, 200, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)

	data := app.newTemplateData(r)

	data.Snippet = snippet

	app.Render(w, 200, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.Render(w, http.StatusOK, "create.tmpl.html", data)
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	var form snippetCreateForm
	err := app.decodePostForm(r, &form)

	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type userCreateForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userCreateForm{}
	app.Render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var userForm userCreateForm
	err := app.decodePostForm(r, &userForm)

	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	//check for empty fields
	userForm.CheckField(validator.NotBlank(userForm.Name), "name", "Name cannot be null")
	userForm.CheckField(validator.NotBlank(userForm.Email), "email", "email cannot be null")
	userForm.CheckField(validator.NotBlank(userForm.Password), "password", "password cannot be null")
	//validate name, email, password

	userForm.CheckField(validator.ValidEmail(userForm.Email), "email", "Invalid Email")
	userForm.CheckField(validator.ValidUserName(userForm.Name), "name", "Name must be alpha numeric b/w 4 to 16 chars")
	userForm.CheckField(validator.ValidUserName(userForm.Password), "password", "password must be greater than 8 chars")

	if !userForm.Valid() {
		data := app.newTemplateData(r)
		data.Form = userForm
		app.Render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(userForm.Name, userForm.Email, userForm.Password)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			userForm.AddFieldError("email", "Email Address is already in use")
			data := app.newTemplateData(r)
			data.Form = userForm
			app.Render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
			return
		} else {
			app.ServerError(w, err)
			return
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.Render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
