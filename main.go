package main

import (
  "fmt"
  "net/http"
  "html/template"
  "database/sql"
  _"github.com/go-sql-driver/mysql"
)

type Article struct {
  Id uint16 // >0
  Title, Anons, Full_text string //can be less then 0
}


type User struct {
  Name string //char?
  Age uint16 // >0
  Trust_lvl int16 //can be less then 0
  Avg_g, Hpnss float64
  Skills [] string
}

func (u User)  getAllInfo() string {
  return fmt.Sprintf("Username is %s. She is %d, her trust " +
    "level is %d", u.Name, u.Age, u.Trust_lvl)

}

func (u *User) setNewName(newName string)  {
  u.Name = newName
}

func create(w http.ResponseWriter, r *http.Request)  {
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

func index(w http.ResponseWriter, r *http.Request)  {

  t, err := template.ParseFiles("templates/header.html", "templates/footer.html", "templates/index.html")

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

  for res.Next() {
    var post Article
    err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_text) //scan if any
    if err != nil{
      panic(err)
    }

  fmt.Println(fmt.Sprintf("Post: %s with id %d", post.Title, post.Id))

  t.ExecuteTemplate(w, "index", nil)
  }
}

func handleFunc()  {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
  http.HandleFunc("/", index)
  http.HandleFunc("/create/", create)
  http.HandleFunc("/save_article", save_article)
  http.ListenAndServe(":8080", nil)
}


func main(){
//2:14
  handleFunc()
}
