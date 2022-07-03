package main

import (
  "fmt"
  "net/http"
  // "log"
  "html/template"
  "github.com/gorilla/mux"
  "database/sql"
  _"github.com/go-sql-driver/mysql"
  "golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

type Article struct {
  Id uint16 // >0
  Title, Anons, Full_text string //can be less then 0
}

type Pass struct {
  Login, Password, StorageAccess string
  Success bool
}

var posts = []Article{}
var showPost Article
// var (
//   tmpl = template.Must(template.ParseFiles("forms.html"))
// )

type User struct {
  Name string //char?
  Age uint16 // >0
  Trust_lvl int16 //can be less then 0
  Avg_g, Hpnss float64
  Skills [] string
}

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}

  	username := req.FormValue("username")
  	password := req.FormValue("password")
    var user string

    err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

    switch {
    case err == sql.ErrNoRows:
      hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
      if err != nil {
        http.Error(res, "Server error, unable to create your account.", 500)
        return
      }
      _, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
      if err != nil {
        http.Error(res, "Server error, unable to create your account.", 500)
        return
      }
      res.Write([]byte("User created!"))
  		return
  	case err != nil:
  		http.Error(res, "Server error, unable to create your account.", 500)
  		return
  	default:
  		http.Redirect(res, req, "/", 301)
  	}
  }


  func loginPage(res http.ResponseWriter, req *http.Request) {
  	if req.Method != "POST" {
  		http.ServeFile(res, req, "login.html")
  		return
  	}

  	username := req.FormValue("username")
  	password := req.FormValue("password")

  	var databaseUsername string
  	var databasePassword string

  	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

  	if err != nil {
  		http.Redirect(res, req, "/login", 301)
  		return
  	}

  	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
  	if err != nil {
  		http.Redirect(res, req, "/login", 301)
  		return
  	}

  	res.Write([]byte("Hello" + databaseUsername))

  }
// func login(w http.ResponseWriter, r *http.Request)  {
//   data := Pass {
//     Login: r.FormValue("login")
//     Password: r.FormValue("password")
//   }
//   tmpl.ExecuteTemplate(w, "forms.html", nil)
// }

func (u User)  getAllInfo() string {
  return fmt.Sprintf("Username is %s. She is %d, her trust " +
    "level is %d", u.Name, u.Age, u.Trust_lvl)

}

func (u *User) setNewName(newName string)  {
  u.Name = newName
}

func create(w http.ResponseWriter, r *http.Request)  {
  // t, err := template.ParseFiles("templates/check_pass.html", "templates/header.html",
  // "templates/footer.html")
  // if err != nil {
  //   fmt.Fprintf(w, err.Error())
  // }
  //
  // t.ExecuteTemplate(w, "check_pass", nil)
  //
  //   inputPassword2 := r.FormValue("inputPassword2")
  //
  //   if inputPassword2 == "memmove" || inputPassword2 == "guestpass" {
  //     http.Redirect(w, r, "/create/", http.StatusSeeOther)
  //     t.ExecuteTemplate(w, "create", nil)
  //   } else {
  //       // fmt.Fprint(w, "Раньше говорили я бы с ним в разведку не пошел, я б с тобой в тур не поехал, ты проверку не прошел")
  // }

  t, err := template.ParseFiles("templates/create.html", "templates/header.html",
  "templates/footer.html")
  if err != nil {
    fmt.Fprintf(w, err.Error())
  }

  t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request)  {
  title := r.FormValue("title")
  anons := r.FormValue("anons")
  full_text := r.FormValue("full_text")


  if title == "" || anons == "" || full_text == "" {
    fmt.Fprint(w, "Enter all data")
  } else {
      db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golfing")//access to database fuck autocorrect WTF is golfing lol???
      if err != nil {
        panic(err)
      }

      defer db.Close()

      insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`,`full_text`) VALUES('%s', '%s', '%s')", title, anons, full_text))
      if err != nil {
        panic(err)
      }
      defer insert.Close()

      http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func show_post(w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)

  t, err := template.ParseFiles("templates/header.html", "templates/footer.html", "templates/show.html")

  if err != nil {
    fmt.Fprintf(w, err.Error())
    return
  }

  db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golfing")//access to database fuck autocorrect WTF is golfing lol???
  if err != nil {
    panic(err)
  }

  defer db.Close()
  res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
  	if err != nil{
  		panic(err)
  	}
  	showPost = Article{}
  	for res.Next(){
  		var post Article
  		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_text)
  		if err != nil{
  			panic(err)
  		}

  		showPost = post

  	}

  	t.ExecuteTemplate(w, "show", showPost)

}

func index(w http.ResponseWriter, r *http.Request)  {

  t, err := template.ParseFiles("templates/header.html", "templates/footer.html", "templates/index.html")

  if err != nil {
    fmt.Fprintf(w, err.Error())
    return
  }

  // db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golfing")
  // if err != nil {
  //   panic(err)
  // }
  //
  // defer db.Close()
  //
  // res, err := db.Query("SELECT * FROM `articles`")
  // if err != nil{
  //   panic(err)
  // }
  //
  //
  // posts = []Article{}
  // for res.Next() {
  //   var post Article
  //   err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_text) //scan if any
  //   if err != nil{
  //         panic(err)
  //         }
  //
  // // fmt.Println(fmt.Sprintf("Post: %s with id %d", post.Title, post.Id))
  t.ExecuteTemplate(w, "index", nil)
  //   posts = append(posts, post)
  // }
    // t.ExecuteTemplate(w, "index", posts)
}

func display_posts(w http.ResponseWriter, r *http.Request)  {
  t, err := template.ParseFiles("templates/header.html", "templates/footer.html", "templates/dachi_archives.html")

  if err != nil {
    fmt.Fprintf(w, err.Error())
    return
  }

  db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golfing")
  if err != nil {
    panic(err)
  }

  defer db.Close()

  res, err := db.Query("SELECT * FROM `articles`")
  if err != nil{
    panic(err)
  }


  posts = []Article{}
  for res.Next() {
    var post Article
    err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_text) //scan if any
    if err != nil{
          panic(err)
          }

  // // fmt.Println(fmt.Sprintf("Post: %s with id %d", post.Title, post.Id))
  // t.ExecuteTemplate(w, "dachi_archives", nil)
    posts = append(posts, post)
  }
    t.ExecuteTemplate(w, "dachi_archives", posts)
}

func check_pass(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("templates/checkpass.html", "templates/header.html",
  "templates/footer.html")
  if err != nil {
    fmt.Fprintf(w, err.Error())
  }

  t.ExecuteTemplate(w, "checkpass", nil)
  // http.Redirect(res, req, "/login", 301)
    // Passcode := r.FormValue("Passcode")
  //
  //   if Passcode == "memmove" || Passcode == "guestpass" {
  //     http.Redirect(w, r, "/create/", http.StatusSeeOther)
  //   } else {
  //       fmt.Fprint(w, "Раньше говорили я бы с ним в разведку не пошел, я б с тобой в тур не поехал, ты проверку не прошел")
  // // }
  // if Passcode != "memmove" {
  //   fmt.Fprint(w, "Enter correct key")
  // } else {
  //       http.Redirect(w, r, "/create", http.StatusSeeOther)
  //     }
}

func passfail(w http.ResponseWriter, r *http.Request)  {

  t, err := template.ParseFiles("templates/passfail.html", "templates/header.html",
  "templates/footer.html")
  if err != nil {
    fmt.Fprintf(w, err.Error())
  }

  t.ExecuteTemplate(w, "passfail", nil)
}

func pass_correct(w http.ResponseWriter, r *http.Request)  {
    Passcode := r.FormValue("Passcode")

    if Passcode == "memmove" {
      http.Redirect(w, r, "/create", http.StatusSeeOther)
    } else {
      http.Redirect(w, r, "/passfail", http.StatusSeeOther)
  }
}


func handleFunc()  {

  rtr := mux.NewRouter()
  rtr.HandleFunc("/", index).Methods("GET")
  rtr.HandleFunc("/checkpass/", check_pass).Methods("GET")
  rtr.HandleFunc("/pass_correct", pass_correct).Methods("POST")
  rtr.HandleFunc("/create", create).Methods("GET")
  rtr.HandleFunc("/passfail", passfail).Methods("GET")
  rtr.HandleFunc("/save_article", save_article).Methods("POST")
  rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
  rtr.HandleFunc("/display/", display_posts).Methods("GET")
  http.Handle("/", rtr)
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
  http.ListenAndServe(":8080", nil)
}


func main(){
//2:14
  handleFunc()
}
